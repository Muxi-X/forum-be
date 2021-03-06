package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"forum/util"
	"strconv"
)

func (s *PostService) UpdatePostInfo(_ context.Context, req *pb.UpdatePostInfoRequest, _ *pb.Response) error {
	logger.Info("PostService UpdatePostInfo")

	if req.Title == "" || req.Content == "" {
		return errno.ServerErr(errno.ErrBadRequest, "title and content can't be null")
	}

	post, err := s.Dao.GetPost(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if post == nil {
		return errno.NotFoundErr(errno.ErrItemNotFound, "post-"+strconv.Itoa(int(req.Id)))
	}

	post.Title = req.Title
	post.Content = req.Content
	post.LastEditTime = util.GetCurrentTime()
	post.Category = req.Category

	if err := post.Save(); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
