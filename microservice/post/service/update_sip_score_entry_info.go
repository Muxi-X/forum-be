package service

import (
	"context"
	"errors"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"

	"gorm.io/gorm"
)

func (s *PostService) UpdateSipScoreEntryInfo(_ context.Context, req *pb.UpdateSipScoreEntryInfoRequest, _ *pb.Response) error {
	logger.Info("PostService UpdateSipScoreEntryInfo")

	lastModifiedBy := req.GetLastModifiedBy()
	if lastModifiedBy == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "last_modified_by required")
	}

	updateMask := req.GetUpdateMask()
	if updateMask == nil || len(updateMask.Paths) == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "update_mask required")
	}

	fieldMap := map[string]interface{}{
		"name":        req.GetName(),
		"description": req.GetDescription(),
		"cover_img":   req.GetCoverImg(),
	}

	update := map[string]interface{}{
		"last_modified_by": lastModifiedBy,
	}
	for _, path := range req.UpdateMask.Paths {
		if val, ok := fieldMap[path]; ok {
			update[path] = val
		} else {
			return errno.ServerErr(errno.ErrBadRequest, "invalid update_mask path: "+path)
		}
	}

	err := s.Dao.UpdateSipScoreEntry(req.GetSipScoreId(), req.GetSipScoreEntryId(), update)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ServerErr(errno.ErrItemNotFound, "sip_score_entry not found")
		}
		return errno.ServerErr(errno.ErrDatabase, "failed to update sip_score_entry: "+err.Error())
	}

	return nil
}
