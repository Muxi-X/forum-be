package service

import (
	"context"
	"errors"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"forum/pkg/unique"
	"strings"

	"gorm.io/gorm"
)

// UpdateSipScoreInfo
// NOTE:
// 只更新 SipScore 基本信息和标签，不更新关联的统计信息。
// 更新标签属于重操作，涉及多次 DB 访问和 Redis 统计修正，尽量避免频繁更新。
// 更新流程如下：
// 1. 更新 SipScore 基本信息 -> 1 次 DB
//
// =====以下步骤仅在标签更新时执行：=====
// 2. 获取 SipScore 当前信息（用于获取旧 category） -> 1 次 DB
// 3. 获取旧 tagIDs -> 1 次 DB
// 4. 删除旧标签关联 -> 1 次 DB
// 5. BatchGetOrCreateTags 获取新 tagID（可能包含 redis + DB 多次操作）
// 6. 批量创建新的 SipScoreTag 关联 -> 1 次 DB
// 7. 事务提交后，异步更新 Redis ZSet 统计
// todo 优化的话就是放在消息队列里
func (s *PostService) UpdateSipScoreInfo(_ context.Context, req *pb.UpdateSipScoreInfoRequest, _ *pb.Response) error {
	logger.Info("PostService UpdateSipScoreInfo")

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
		"domain":      req.GetDomain(),
		"category":    req.GetCategory(),
	}

	update := map[string]interface{}{
		"last_modified_by": lastModifiedBy,
	}
	var uniqueTags []string
	isTagsUpdate := false

	for _, path := range req.UpdateMask.Paths {
		if strings.Compare(path, "tags") == 0 {
			isTagsUpdate = true
			uniqueTags = unique.UniqueStrings(req.GetTags())
			for _, content := range uniqueTags {
				if content == "" {
					return errno.ServerErr(errno.ErrBadRequest, "tag content cannot be empty")
				}
			}
			continue
		}

		if val, ok := fieldMap[path]; ok {
			update[path] = val
		} else {
			return errno.ServerErr(errno.ErrBadRequest, "invalid update_mask path: "+path)
		}
	}

	if isTagsUpdate && len(uniqueTags) == 0 {
		return errno.ServerErr(errno.ErrBadRequest, "tags required")
	}

	var (
		oldTagIDs   []uint32
		newTagIDs   []uint32
		oldCategory string
		newCategory string
	)

	// 开启事务
	fc := func(tx *gorm.DB) error {
		id := req.GetId()

		// 更新主体基本信息
		err := s.Dao.UpdateSipScore(id, update, tx)
		if err != nil {
			return err
		}

		if !isTagsUpdate {
			return nil
		}

		sipScore, err := s.Dao.GetSipScore(id, tx)
		if err != nil {
			return err
		}
		oldCategory = sipScore.Category
		// 如果 category 没有更新，保持原来的 category
		if _, ok := update["category"]; ok {
			newCategory = req.GetCategory()
		} else {
			newCategory = oldCategory
		}

		// 获取旧 tagIDs
		oldTagIDs, err = s.Dao.ListTagIDsBySipScoreId(id, tx)
		if err != nil {
			return err
		}

		// 删除旧 sipScore tags
		if err = s.Dao.DeleteSipScoreTagsBySipScoreId(id, tx); err != nil {
			return err
		}

		// 获取新 tagID
		tagsModel, err := s.Dao.BatchGetOrCreateTags(uniqueTags)
		if err != nil {
			return err
		}

		// 顺序一样，直接构建
		sipScoreTags := make([]*dao.SipScoreTagModel, 0, len(uniqueTags))
		newTagIDs = make([]uint32, 0, len(tagsModel))

		for _, tag := range tagsModel {
			newTagIDs = append(newTagIDs, tag.Id)

			sipScoreTags = append(sipScoreTags, &dao.SipScoreTagModel{
				TagID:      tag.Id,
				SipScoreID: id,
			})
		}

		if err = s.Dao.BatchCreateSipScoreTags(sipScoreTags, tx); err != nil {
			return err
		}

		return nil
	}

	err := s.Dao.Transaction(fc)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ServerErr(errno.ErrItemNotFound, "sip score not found")
		}
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if isTagsUpdate {
		go func(oldTagIDs, newTagIDs []uint32, oldCategory, newCategory string) {
			if len(oldTagIDs) != 0 {
				_ = s.Dao.BatchRemoveTagsFromSortedSet(oldTagIDs, oldCategory)
			}
			if len(newTagIDs) != 0 {
				_ = s.Dao.BatchAddTagsToSortedSet(newTagIDs, newCategory)
			}
		}(oldTagIDs, newTagIDs, oldCategory, newCategory)
	}

	return nil
}
