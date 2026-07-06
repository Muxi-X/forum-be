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

func (s *PostService) CreatePostComment(_ context.Context, req *pb.CreatePostCommentRequest, resp *pb.CreatePostCommentResponse) error {
	logger.Info("PostService CreatPostComment")

	post, err := s.Dao.GetPost(req.PostId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// check if the FatherId is valid
	switch req.TypeName {
	case constvar.SubPost:
		req.FatherId = req.PostId
		resp.UserId = post.CreatorId
		resp.FatherContent = post.Title

	case constvar.FirstLevelComment:
		if req.FatherId == 0 {
			req.FatherId = req.PostId
		}
		if req.FatherId == req.PostId {
			resp.UserId = post.CreatorId
			resp.FatherContent = post.Title
		} else {
			comment, err := s.Dao.GetComment(req.FatherId)
			if err != nil {
				return errno.ServerErr(errno.ErrDatabase, err.Error())
			}
			if comment == nil {
				return errno.ServerErr(errno.ErrBadRequest, "the comment not found")
			}

			resp.FatherUserId = comment.CreatorId
			resp.UserId = comment.CreatorId
			resp.FatherContent = comment.Content
		}

	case constvar.SecondLevelComment:
		comment, err := s.Dao.GetComment(req.FatherId)
		if err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
		if comment == nil {
			return errno.ServerErr(errno.ErrBadRequest, "the comment not found")
		}

		resp.FatherUserId = comment.CreatorId
		resp.UserId = comment.CreatorId
		resp.FatherContent = comment.Content

	default:
		return errno.ServerErr(errno.ErrBadRequest, "type_name not legal")
	}

	data := &dao.CommentModel{
		TypeName:   req.TypeName,
		Content:    req.Content,
		FatherId:   req.FatherId,
		CreatorId:  req.CreatorId,
		TargetID:   post.Id,
		TargetType: constvar.Post,
		ImgUrl:     req.ImgUrl,
	}

	commentId, err := s.Dao.CreateComment(data)
	if err != nil {
		return databaseErr(err)
	}

	// 二级评论需要递增父评论的子评论数
	if req.TypeName == constvar.SecondLevelComment {
		if err := s.Dao.IncrCommentSubNum(req.FatherId); err != nil {
			logger.Error("incr comment sub_num error", logger.String(err.Error()))
		}
	}

	if err := model.AddPolicy(req.CreatorId, constvar.Comment, commentId, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	if err := model.AddResourceRole(constvar.Comment, commentId, post.Domain); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	go func() {
		if err := s.Dao.ChangePostScore(req.PostId, constvar.CommentScore); err != nil {
			logger.Error(errno.ErrChangeScore.Error(), logger.String(err.Error()))
		}
	}()

	commentInfo, err := s.Dao.GetCommentInfo(commentId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.Id = commentId
	resp.TypeName = post.Domain
	resp.CreateTime = util.GetCurrentTime()
	resp.CreatorName = commentInfo.CreatorName
	resp.CreatorAvatar = commentInfo.CreatorAvatar

	return nil
}
