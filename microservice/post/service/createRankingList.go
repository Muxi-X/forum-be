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

func (s *PostService) CreateRankingList(_ context.Context, req *pb.CreateRankingListRequest, resp *pb.CreateRankingListResponse) error {
	logger.Info("PostService CreateRankingList")

	domain := req.GetDomain()
	if domain != constvar.NormalDomain && domain != constvar.MuxiDomain {
		return errno.ServerErr(errno.ErrBadRequest, "domain not legal")
	}

	data := &dao.RankingListModel{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		CoverImg:    req.GetCoverImg(),
		CreatorID:   req.GetCreatorId(),
		Domain:      domain,
		Category:    req.GetCategory(),
	}

	rankingListID, err := s.Dao.CreateRankingList(data)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// 创建者具有写权限
	if err = model.AddPolicy(req.CreatorId, constvar.RankingList, rankingListID, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	if err = model.AddResourceRole(constvar.RankingList, rankingListID, domain); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	rankingTags := req.GetTags()
	uniqueTags := unique.UniqueStrings(rankingTags)
	for _, content := range uniqueTags {
		if content == "" {
			return errno.ServerErr(errno.ErrBadRequest, "tag content cannot be empty")
		}
	}

	// 获取 tagID
	tagsModel, err := s.Dao.BatchGetOrCreateTags(uniqueTags)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// 顺序一样，直接构建
	rankingListTags := make([]*dao.RankingListTagModel, 0, len(uniqueTags))
	for _, tag := range tagsModel {
		rankingListTags = append(rankingListTags, &dao.RankingListTagModel{
			TagID:         tag.Id,
			RankingListID: rankingListID,
		})
	}

	err = s.Dao.BatchCreateRankingListTags(rankingListTags)
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

	resp.Id = rankingListID
	return nil
}
