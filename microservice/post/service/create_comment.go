package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *PostService) CreateComment(_ context.Context, req *pb.CreateCommentRequest, resp *pb.CreateCommentResponse) error {
	logger.Info("PostService CreateComment")
	if req.TargetType == "" {
		return errno.ServerErr(errno.ErrBadRequest, "target_type is required")
	}
	if req.TargetId == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "target_id is required")
	}

	switch req.TargetType {
	case constvar.Post:
		return s.createPostComment(req, resp)

	case constvar.SipScoreEntryCommentRating:
		return s.createSipScoreEntryCommentRatingComment(req, resp)

	default:
		return errno.ServerErr(errno.ErrBadRequest, "unsupported target_type: "+req.TargetType)
	}
}

func (s *PostService) createPostComment(req *pb.CreateCommentRequest, resp *pb.CreateCommentResponse) error {
	post, err := s.Dao.GetPost(req.TargetId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	switch req.TypeName {
	case constvar.SubPost:
		req.FatherId = req.TargetId
		resp.UserId = post.CreatorId
		resp.FatherContent = post.Title

	case constvar.FirstLevelComment:
		if req.FatherId == 0 {
			req.FatherId = req.TargetId
		}
		if req.FatherId == req.TargetId {
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

		resp.BeRepliedUserId = comment.CreatorId
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
		if err := s.Dao.ChangePostScore(req.TargetId, constvar.CommentScore); err != nil {
			logger.Error(errno.ErrChangeScore.Error(), logger.String(err.Error()))
		}
	}()

	commentInfo, err := s.Dao.GetCommentInfo(commentId)
	if err != nil {
		return databaseErr(err)
	}

	resp.Id = commentId
	resp.TypeName = post.Domain
	resp.CreatedAt = timestamppb.New(commentInfo.CreatedAt)
	resp.CreatorName = commentInfo.CreatorName
	resp.CreatorAvatar = commentInfo.CreatorAvatar
	resp.TargetType = constvar.Post
	resp.TargetId = post.Id

	return nil
}

// todo 由网关或者这里检测权限 -> 该用户存在对该榜单的阅读权限
func (s *PostService) createSipScoreEntryCommentRatingComment(req *pb.CreateCommentRequest, resp *pb.CreateCommentResponse) error {
	rating, err := s.Dao.GetSipScoreEntryCommentRatingByID(req.TargetId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	if rating == nil {
		return errno.ServerErr(errno.ErrBadRequest, "rating not found")
	}

	data := &dao.CommentModel{
		Content:    req.Content,
		FatherId:   req.FatherId,
		CreatorId:  req.CreatorId,
		TargetID:   rating.ID,
		TargetType: constvar.SipScoreEntryCommentRating,
		ImgUrl:     req.ImgUrl,
		TypeName:   constvar.SecondLevelComment, // 一般而言都是二级评论
	}

	var commentId uint32
	err = s.Dao.Transaction(func(tx *gorm.DB) error {
		id, err := s.Dao.CreateComment(data, tx)
		if err != nil {
			return err
		}
		commentId = id

		if err := s.Dao.IncrSipScoreEntryCommentRatingCommentNum(rating.SipScoreID, rating.EntryID, rating.ID, tx); err != nil {
			return err
		}

		if err := s.Dao.IncrSipScoreEntryCommentCount(rating.SipScoreID, rating.EntryID, tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return databaseErr(err)
	}

	if err := model.AddPolicy(req.CreatorId, constvar.Comment, commentId, constvar.Write); err != nil {
		return errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	commentInfo, err := s.Dao.GetCommentInfo(commentId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.Id = commentId
	resp.UserId = rating.CreatorID
	resp.CreatedAt = timestamppb.New(commentInfo.CreatedAt)
	resp.CreatorName = commentInfo.CreatorName
	resp.CreatorAvatar = commentInfo.CreatorAvatar
	resp.TargetType = constvar.SipScoreEntryCommentRating
	resp.TargetId = rating.ID
	resp.TypeName = commentInfo.TypeName
	resp.BeRepliedUserId = rating.CreatorID
	resp.FatherContent = rating.Content

	return nil
}
