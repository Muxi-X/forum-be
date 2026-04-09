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
)

// todo 可以进一步优化 - 消息队列

func (s *PostService) CreateSipScore(_ context.Context, req *pb.CreateSipScoreRequest, resp *pb.CreateSipScoreResponse) error {
	logger.Info("PostService CreateSipScore")

	// 参数检验
	domain := req.GetDomain()
	if domain != constvar.NormalDomain && domain != constvar.MuxiDomain {
		return errno.ServerErr(errno.ErrBadRequest, "domain not legal")
	}

	tags := req.GetTags()
	uniqueTags := unique.UniqueStrings(tags)
	for _, content := range uniqueTags {
		if content == "" {
			return errno.ServerErr(errno.ErrBadRequest, "tag content cannot be empty")
		}
	}

	data := &dao.SipScoreModel{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		CoverImg:    req.GetCoverImg(),
		CreatorID:   req.GetCreatorId(),
		Domain:      domain,
		Category:    req.GetCategory(),
	}

	sipScoreID, err := s.Dao.CreateSipScore(data)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// 创建者具有写权限
	if err = model.AddPolicy(req.CreatorId, constvar.SipScore, sipScoreID, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	if err = model.AddResourceRole(constvar.SipScore, sipScoreID, domain); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	// 获取 tagID
	tagsModel, err := s.Dao.BatchGetOrCreateTags(uniqueTags)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// 顺序一样，直接构建
	sipScoreTags := make([]*dao.SipScoreTagModel, 0, len(uniqueTags))
	for _, tag := range tagsModel {
		sipScoreTags = append(sipScoreTags, &dao.SipScoreTagModel{
			TagID:      tag.Id,
			SipScoreID: sipScoreID,
		})
	}

	err = s.Dao.BatchCreateSipScoreTags(sipScoreTags)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	tagIDs := make([]uint32, 0, len(tagsModel))
	for _, tag := range tagsModel {
		tagIDs = append(tagIDs, tag.Id)
	}

	category := req.GetCategory()

	go func(tagIDs []uint32, category string) {
		if err := s.Dao.BatchAddTagsToSortedSet(tagIDs, category); err != nil {
			logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
		}
	}(tagIDs, category)

	resp.Id = sipScoreID
	return nil
}
