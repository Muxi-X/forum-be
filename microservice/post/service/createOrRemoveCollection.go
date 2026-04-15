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

func (s *PostService) CreateOrRemoveCollection(_ context.Context, req *pb.ToggleTargetRequest, resp *pb.CreateOrRemoveCollectionResponse) error {
	logger.Info("PostService CreateOrRemoveCollection")

	collection := &dao.CollectionModel{
		UserID:      req.GetUserId(),
		ContentID:   req.GetTargetType(),
		ContentType: req.GetTargetType(),
	}

	switch req.GetTargetType() {
	case constvar.CollectionPost:
		return s.createOrRemovePostCollection(collection, req, resp)
	case constvar.CollectionSipScore:
		return s.createOrRemoveSipScoreCollection(collection, req, resp)
	default:
		return errno.ServerErr(errno.ErrBadRequest, "target_type not legal")
	}
}

func (s *PostService) createOrRemovePostCollection(collection *dao.CollectionModel, req *pb.ToggleTargetRequest, resp *pb.CreateOrRemoveCollectionResponse) error {
	var score int

	isCollection, err := s.Dao.IsUserCollected(req.GetUserId(), req.GetTargetType(), req.GetTargetId())
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if isCollection {
		err = s.Dao.DeleteCollection(collection)

		score = -constvar.CollectionScore
	} else {
		resp.Id, err = s.Dao.CreateCollection(collection)

		score = constvar.CollectionScore
	}
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	targetID := req.GetTargetId()
	scoreCopy := score

	go func() {
		if err := s.Dao.ChangePostScore(targetID, scoreCopy); err != nil {
			logger.Error(errno.ErrChangeScore.Error(), logger.String(err.Error()))
		}
	}()

	post, err := s.Dao.GetPost(req.GetTargetId())
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.UserId = post.CreatorId
	resp.Content = post.Title
	resp.TypeName = post.Domain

	return nil
}

func (s *PostService) createOrRemoveSipScoreCollection(collection *dao.CollectionModel, req *pb.ToggleTargetRequest, resp *pb.CreateOrRemoveCollectionResponse) error {
	return s.Dao.Transaction(func(tx *gorm.DB) error {
		// try delete
		deleted, err := s.Dao.TryDeleteCollection(collection, tx)
		if err != nil {
			return err
		}

		if deleted {
			// 取消收藏
			if err = s.Dao.DecrSipScoreCollectCount(req.GetTargetId(), tx); err != nil {
				return err
			}
		} else {
			// try create
			created, err := s.Dao.TryCreateCollection(collection, tx)
			if err != nil {
				return err
			}

			if created {
				// 收藏成功
				if err = s.Dao.IncrSipScoreCollectCount(req.GetTargetId(), tx); err != nil {
					return err
				}
			}
			// 如果 created == false 说明并发导致，忽略
		}

		sipScore, err := s.Dao.GetSipScore(req.GetTargetId(), tx)
		if err != nil {
			return err
		}

		resp.UserId = sipScore.CreatorID
		resp.Content = sipScore.Name
		resp.TypeName = sipScore.Domain

		return nil
	})
}
