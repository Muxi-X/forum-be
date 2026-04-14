package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

func (s *PostService) ListCollection(ctx context.Context, req *pb.ListPostPartInfoRequest, resp *pb.ListPostPartInfoResponse) error {
	logger.Info("PostService ListCollections")

	postIds, err := s.Dao.ListCollectionByUserId(req.TargetUserId, constvar.CollectionPost)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	var filter = &dao.PostModel{}

	domain, err := s.GetUserDomain(ctx, req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrRPC, err.Error())
	}
	if domain == constvar.NormalDomain {
		filter.Domain = domain
	}

	posts, err := s.Dao.ListPostInfoByPostIds(postIds, filter, req.Offset, req.Limit, req.LastId, req.Pagination)
	if err != nil {
		return errno.ServerErr(errno.ErrListPostInfoByPostIds, err.Error())
	}

	resp.Posts = posts

	return nil
}
