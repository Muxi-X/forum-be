package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

// todo 需要再传一个参数，表示什么类型的收藏，帖子、评论等

func (s *PostService) CreateOrRemoveCollection(_ context.Context, req *pb.Request, resp *pb.CreateOrRemoveCollectionResponse) error {
	logger.Info("PostService CreateOrRemoveCollection")

	var score int

	collection := &dao.CollectionModel{
		UserID:      req.UserId,
		ContentID:   req.Id,
		ContentType: constvar.CollectionPost,
	}

	isCollection, err := s.Dao.IsUserCollected(req.UserId, constvar.CollectionPost, req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if isCollection {
		err = s.Dao.DeleteCollection(collection)

		score = -constvar.CollectionScore
	} else {
		resp.Id, err = s.Dao.CreateCollection(collection)

		score = constvar.CollectionScore
	}
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	go func() {
		if err := s.Dao.ChangePostScore(req.Id, score); err != nil {
			logger.Error(errno.ErrChangeScore.Error(), logger.String(err.Error()))
		}
	}()

	post, err := s.Dao.GetPost(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.UserId = post.CreatorId
	resp.Content = post.Title
	resp.TypeName = post.Domain

	return nil
}
