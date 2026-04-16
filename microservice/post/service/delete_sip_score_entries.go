package service

import (
	"context"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"

	"gorm.io/gorm"
)

func (s *PostService) DeleteSipScoreEntries(_ context.Context, req *pb.DeleteSipScoreEntriesRequest, _ *pb.Response) error {
	logger.Info("PostService DeleteSipScoreEntries")

	sipScoreID := req.GetSipScoreId()
	if sipScoreID == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "sipScoreId cannot be 0")
	}

	entryIDs := req.GetSipScoreEntryId()
	if len(entryIDs) == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "sipScoreEntryId cannot be empty")
	}

	var deletedIDs []uint32
	fc := func(tx *gorm.DB) error {
		var err error

		deletedIDs, err = s.Dao.GetSipScoreEntryIDs(sipScoreID, entryIDs, tx)
		if err != nil {
			return err
		}

		if len(deletedIDs) == 0 {
			return nil
		}

		entryCount, participantCount, err := s.Dao.GetSipScoreEntryStats(sipScoreID, deletedIDs, tx)
		if err != nil {
			return err
		}

		if err = s.Dao.DeleteSipScoreEntries(sipScoreID, deletedIDs, tx); err != nil {
			return err
		}

		return s.Dao.DecrSipScoreStats(sipScoreID, entryCount, participantCount, tx)
	}

	err := s.Dao.Transaction(fc)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// TODO 投递 MQ
	// 删除 sipScoreEntryComment
	// 删除 sipScoreEntryReview
	// ……

	return nil
}
