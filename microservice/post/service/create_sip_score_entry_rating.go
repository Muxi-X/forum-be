package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"gorm.io/gorm"
)

func (s *PostService) CreateSipScoreEntryCommentRating(ctx context.Context, req *pb.CreateSipScoreEntryCommentRatingRequest, resp *pb.Response) error {
	logger.Info("PostService CreateSipScoreEntryCommentRating")

	rating := req.GetRating()
	if rating < 1 || rating > 5 {
		return errno.ServerErr(errno.ErrBadRequest, "rating must be between 1 and 5")
	}

	return s.Dao.Transaction(func(tx *gorm.DB) error {
		// 1. 检查 entry 是否存在
		_, err := s.Dao.GetSipScoreEntry(req.GetSipScoreId(), req.GetSipScoreEntryId(), tx)
		if err != nil {
			return err
		}

		// 2. 创建评论
		comment := &dao.CommentModel{
			TargetID:   req.GetSipScoreEntryId(),
			TargetType: constvar.SipScoreEntry,
			TypeName:   constvar.FirstLevelComment,
			Content:    req.GetComment(),
			FatherId:   0,
			CreatorId:  req.GetUserId(),
			LikeNum:    0,
			ImgUrl:     req.GetImgUrl(),
			IsReport:   false,
		}

		commentID, err := s.Dao.CreateComment(comment, tx)
		if err != nil {
			return err
		}

		// 3. 创建评分记录
		data := &dao.SipScoreEntryCommentRating{
			CreatorID:      req.GetUserId(),
			LastModifiedBy: req.GetUserId(),
			SipScoreID:     req.GetSipScoreId(),
			EntryID:        req.GetSipScoreEntryId(),
			CommentID:      commentID,
			LikeNum:        0,
		}

		_, err = s.Dao.CreateSipScoreEntryCommentRating(data, tx)
		if err != nil {
			return err
		}

		// 4. 更新 SipScore participant_count
		err = s.Dao.IncrSipScoreParticipantCount(req.GetSipScoreId(), 1, tx)
		if err != nil {
			return err
		}

		// 5. 更新 entry 统计字段
		err = s.Dao.IncrSipScoreEntryScore(req.GetSipScoreId(), req.GetSipScoreId(), rating, 1, tx)

		return err
	})
}
