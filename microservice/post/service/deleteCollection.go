package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	"forum-post/worker"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

// DeleteCollection
// NOTE: 这个好像没有用到，废弃废弃
func (s *PostService) DeleteCollection(_ context.Context, req *pb.Request, _ *pb.Response) error {
	logger.Info("PostService DeleteCollection")

	collection := &dao.CollectionModel{
		ID: req.Id,
	}

	if err := s.Dao.DeleteCollection(collection); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if err := model.DeletePermission(req.UserId, constvar.Collection, req.Id, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	err := worker.PublishPostInteraction(worker.PostInteractionWriter, &worker.InteractionMessage{
		Type:   worker.TypeCollection,
		PostId: req.Id,
		Score:  -constvar.CollectionScore,
	})
	if err != nil {
		logger.Error("publish post interaction error", logger.String(err.Error()))
	}

	return nil
}
