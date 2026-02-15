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

	if req.TypeName == constvar.Post {
		err = s.Dao.DeletePost(req.Id)
	} else if req.TypeName == constvar.Comment {
		err = s.Dao.DeleteComment(req.Id)
	} else if req.TypeName == constvar.QualityPost {
		err = s.Dao.ChangeQualityPost(req.Id, false)
	} else {
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
