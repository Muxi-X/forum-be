package service

import (
	"context"
	"errors"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"
	"sync"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *PostService) GetSipScore(_ context.Context, req *pb.Request, resp *pb.SipScore) error {
	logger.Info("PostService GetSipScore")

	sipScoreID := req.GetId()
	if sipScoreID == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "id required")
	}
	sipScore, err := s.Dao.GetSipScore(req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ServerErr(errno.ErrItemNotFound, "sip_score-"+strconv.Itoa(int(sipScoreID)))
		}
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	var (
		collected bool
		tags      []string
	)

	var wg sync.WaitGroup
	wg.Add(2)

	// 收藏状态
	go func() {
		defer wg.Done()
		c, err := s.Dao.IsUserCollected(req.UserId, constvar.CollectionSipScore, sipScoreID)
		if err != nil {
			logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
			return
		}

		collected = c
	}()

	// tags
	go func() {
		defer wg.Done()
		t, _, err := s.Dao.ListTagsBySipScoreId(sipScoreID)
		if err != nil {
			logger.Error(errno.ErrDatabase.Error(), logger.String(err.Error()))
			return
		}

		tags = t
	}()

	wg.Wait()

	resp.Id = sipScoreID
	resp.CreatedAt = timestamppb.New(sipScore.CreatedAt)
	resp.UpdatedAt = timestamppb.New(sipScore.UpdatedAt)
	resp.CreatorId = sipScore.CreatorID
	resp.LastModifiedBy = sipScore.LastModifiedBy
	resp.EntryCount = sipScore.EntryCount
	resp.CollectCount = sipScore.CollectCount
	resp.ParticipantCount = sipScore.ParticipantCount
	resp.Name = sipScore.Name
	resp.Description = sipScore.Description
	resp.CoverImg = sipScore.CoverImg
	resp.Domain = sipScore.Domain
	resp.Category = sipScore.Category
	resp.Tags = tags
	resp.IsCollected = collected

	return nil
}
