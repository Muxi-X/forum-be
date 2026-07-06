package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/util"
)

func (s *PostService) CreatePost(_ context.Context, req *pb.CreatePostRequest, resp *pb.CreatePostResponse) error {
	logger.Info("PostService CreatePost")

	if req.Domain != constvar.NormalDomain && req.Domain != constvar.MuxiDomain {
		return errno.ServerErr(errno.ErrBadRequest, "domain not legal")
	}

	data := &dao.PostModel{
		Domain:          req.Domain,
		Content:         req.Content,
		Title:           req.Title,
		CreateTime:      util.GetCurrentTime(),
		Category:        req.Category,
		Re:              false,
		CreatorId:       req.UserId,
		ContentType:     req.ContentType,
		CompiledContent: req.CompiledContent,
		LastEditTime:    util.GetCurrentTime(),
		Summary:         req.Summary,
	}

	postId, err := s.Dao.CreatePost(data)
	if err != nil {
		return databaseErr(err)
	}

	if err := model.AddPolicy(req.UserId, constvar.Post, postId, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	if err := model.AddResourceRole(constvar.Post, postId, req.Domain); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	for _, content := range req.Tags {
		if content == "" {
			return errno.ServerErr(errno.ErrBadRequest, "tag content can't be null")
		}

		tag, err := s.Dao.GetTagByContent(content)
		if err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}

		item := dao.Post2TagModel{
			PostId: postId,
			TagId:  tag.Id,
		}
		if err := s.Dao.CreatePost2Tag(item); err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}

		go func() {
			if err := s.Dao.AddTagToSortedSet(tag.Id, req.Category); err != nil {
				logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
			}
		}()
	}

	resp.Id = postId

	return nil
}
