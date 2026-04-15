package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/pkg/unique"
	"time"

	"go.uber.org/zap"
)

// todo 可以进一步优化 - 将 tag 放到消息队列
// todo 这里的时间性能分析暂时不删除，暂时没测试到 7s 的情况，前端之后对接的时候看下日志

func (s *PostService) CreateSipScore(_ context.Context, req *pb.CreateSipScoreRequest, resp *pb.CreateSipScoreResponse) error {
	start := time.Now()
	logger.Info("PostService CreateSipScore")

	// 参数检验
	t1 := time.Now()
	domain := req.GetDomain()
	if domain != constvar.NormalDomain && domain != constvar.MuxiDomain {
		return errno.ServerErr(errno.ErrBadRequest, "domain not legal")
	}
	logger.Info("check domain cost:", zap.String("time", time.Since(t1).String()))

	t2 := time.Now()
	tags := req.GetTags()
	uniqueTags := unique.UniqueStrings(tags)
	for _, content := range uniqueTags {
		if content == "" {
			return errno.ServerErr(errno.ErrBadRequest, "tag content cannot be empty")
		}
	}
	logger.Info("check tags cost:", zap.String("time", time.Since(t2).String()))

	t3 := time.Now()
	creatorID := req.GetCreatorId()
	data := &dao.SipScoreModel{
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		CoverImg:       req.GetCoverImg(),
		CreatorID:      creatorID,
		Domain:         domain,
		Category:       req.GetCategory(),
		LastModifiedBy: creatorID,
	}

	sipScoreID, err := s.Dao.CreateSipScore(data)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	logger.Info("CreateSipScore cost:", zap.String("time", time.Since(t3).String()))

	// 创建者具有写权限
	t4 := time.Now()
	if err = model.AddPolicy(req.CreatorId, constvar.SipScore, sipScoreID, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	if err = model.AddResourceRole(constvar.SipScore, sipScoreID, domain); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}
	logger.Info("add casbin policy cost:", zap.String("time", time.Since(t4).String()))

	// 获取 tagID
	t5 := time.Now()
	tagsModel, err := s.Dao.BatchGetOrCreateTags(uniqueTags)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// 顺序一样，直接构建
	sipScoreTags := make([]*dao.SipScoreTagModel, 0, len(uniqueTags))
	tagIDs := make([]uint32, 0, len(tagsModel))
	for _, tag := range tagsModel {
		tagIDs = append(tagIDs, tag.Id)

		sipScoreTags = append(sipScoreTags, &dao.SipScoreTagModel{
			TagID:      tag.Id,
			SipScoreID: sipScoreID,
		})
	}
	logger.Info("BatchGetOrCreateTags cost:", zap.String("time", time.Since(t5).String()))

	t6 := time.Now()
	err = s.Dao.BatchCreateSipScoreTags(sipScoreTags)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	category := req.GetCategory()
	logger.Info("BatchCreateSipScore cost:", zap.String("time", time.Since(t6).String()))

	go func(tagIDs []uint32, category string) {
		if err := s.Dao.BatchAddTagsToSortedSet(tagIDs, category); err != nil {
			logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
		}
	}(tagIDs, category)

	resp.Id = sipScoreID
	logger.Info("PostService CreateSipScore cost:", zap.String("time", time.Since(start).String()))
	return nil
}
