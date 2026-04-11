package service

import (
	"context"
	"errors"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"strings"

	"gorm.io/gorm"
)

func (s *PostService) CreateSipScoreEntry(_ context.Context, req *pb.CreateSipScoreEntryRequest, resp *pb.CreateSipScoreEntryResponse) error {
	logger.Info("PostService CreateSipScoreEntry")

	// 参数检验
	entries := req.GetEntries()
	if len(entries) == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "entries required")
	}

	// todo 这里随便写的 100，之后压测
	if len(entries) > 100 {
		return errno.ServerErr(errno.ErrBadRequest, "too many entries, maximum is 100")
	}

	// 单次请求内去重
	seen := make(map[string]struct{})
	for _, entry := range entries {
		name := strings.TrimSpace(entry.GetName())
		if name == "" {
			return errno.ServerErr(errno.ErrBadRequest, "entry name cannot be empty")
		}
		if entry.GetCoverImg() == "" {
			return errno.ServerErr(errno.ErrBadRequest, "entry cover image cannot be empty")
		}
		if _, ok := seen[name]; ok {
			return errno.ServerErr(errno.ErrBadRequest, "duplicate entry name: "+name)
		}
		seen[name] = struct{}{}
	}

	// 构建
	sipScoreID := req.GetSipScoreId()
	creatorID := req.GetCreatorId()

	sipScoreEntries := make([]*dao.SipScoreEntryModel, 0, len(entries))
	for _, entry := range entries {
		sipScoreEntries = append(sipScoreEntries, &dao.SipScoreEntryModel{
			SipScoreID:       sipScoreID,
			LastModifiedBy:   creatorID,
			CreatorID:        creatorID,
			IsReport:         false,
			ParticipantCount: 0,
			CommentCount:     0,
			ScoreTotal:       0,
			Name:             strings.TrimSpace(entry.GetName()),
			Description:      strings.TrimSpace(entry.GetDescription()),
			CoverImg:         strings.TrimSpace(entry.GetCoverImg()),
		})
	}

	fc := func(tx *gorm.DB) error {
		sipScoreID = req.GetSipScoreId()
		if _, err := s.Dao.GetSipScore(sipScoreID, tx); err != nil {
			return err
		}

		if err := s.Dao.BatchCreateSipScoreEntries(sipScoreEntries, tx); err != nil {
			return err
		}

		incr := int64(len(sipScoreEntries))
		if err := s.Dao.IncrSipScoreEntryCount(sipScoreID, incr, tx); err != nil {
			return err
		}

		return nil
	}

	err := s.Dao.Transaction(fc)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ServerErr(errno.ErrItemNotFound, "sip score not found")
		} else if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errno.ServerErr(errno.ErrDatabase, "duplicate entry name")
		}
		return errno.ServerErr(errno.ErrDatabase, "database transaction error")
	}

	// 构建相应
	resp.EntryIds = make([]uint32, 0, len(sipScoreEntries))
	for _, entry := range sipScoreEntries {
		resp.EntryIds = append(resp.EntryIds, entry.ID)
	}

	return nil
}
