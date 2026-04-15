package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
)

func (s *PostService) DeleteItem(_ context.Context, req *pb.DeleteItemRequest, _ *pb.Response) error {
	logger.Info("PostService DeleteItem")

	var err error

	switch req.TypeName {
	case constvar.Post:
		err = s.Dao.DeletePost(req.Id)
	case constvar.Comment:
		err = s.Dao.DeleteComment(req.Id)
	case constvar.QualityPost:
		err = s.Dao.ChangeQualityPost(req.Id, false)
	case constvar.SipScore:
		err = s.deleteSipScore(req.Id)
	default:
		return errno.ServerErr(errno.ErrBadRequest, "wrong TypeName")
	}

	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if req.TypeName != constvar.QualityPost {
		if err := model.DeletePermission(req.UserId, req.TypeName, req.Id, constvar.Write); err != nil {
			return errno.ServerErr(errno.ErrCasbin, err.Error())
		}
	}

	return nil
}

func (s *PostService) deleteSipScore(id uint32) error {
	// 删除 主体
	if err := s.Dao.DeleteSipScore(id); err != nil {
		return err
	}

	// todo 消息队列中比较好
	// 删除 sipScore tags
	// 删除 sipScoreEntry, sipScoreEntryComment, sipScoreEntryReview
	// 删除 sipScore collection
	return nil
}
