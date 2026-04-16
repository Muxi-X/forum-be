package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/pkg/pagetoken"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// todo 需要针对这个 entry 再获取最热门的几条评价
// todo 还需要获取当前用户对每一个 entry 的评分
// todo 给每一个dao函数加一个ctx,可以取消

type sipScoreHydrate struct {
	collected map[uint32]bool                      // key: sipScoreID
	entries   map[uint32][]*dao.SipScoreEntryModel // key: sipScoreID
	// future:
	// ratings map[uint32]*dao.RatingModel        // key: entryID
	// topComments map[uint32][]*dao.CommentModel // key: entryID
}

func (s *PostService) ListSipScore(ctx context.Context, req *pb.ListSipScoreRequest, resp *pb.ListSipScoreResponse) error {
	logger.Info("PostService ListSipScore")

	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = constvar.DefaultPageSize
	} else if pageSize > constvar.MaxPageSize {
		pageSize = constvar.MaxPageSize
	}
	limit := pageSize + 1

	var pageToken pb.SipScorePageToken
	err := pagetoken.DecodePageToken(req.GetPageToken(), &pageToken)
	if err != nil {
		return errno.ServerErr(errno.ErrBadRequest, "invalid page token")
	}

	var pageTokenPtr *pb.SipScorePageToken
	if req.GetPageToken() != "" {
		pageTokenPtr = &pageToken

		if pageTokenPtr.GetSortType() != req.GetSortType() {
			return errno.ServerErr(errno.ErrBadRequest, "page token sort type mismatch")
		}
	}

	switch req.GetSortType() {
	case constvar.SortByNewest:
		return s.listSipScoreNewest(ctx, req.GetUserId(), pageTokenPtr, limit, resp)
	case constvar.SortByHottest:
		return s.listSipScoreHottest(ctx, req.GetUserId(), pageTokenPtr, limit, resp)
	default:
		return errno.ServerErr(errno.ErrBadRequest, "sort type not legal")
	}
}

func (s *PostService) listSipScoreCommon(
	ctx context.Context, userID uint32, pageToken *pb.SipScorePageToken, limit uint32, resp *pb.ListSipScoreResponse,
	fetch func(token *pb.SipScorePageToken, limit uint32) ([]*dao.SipScoreModel, error),
	nextToken func(last *dao.SipScoreModel) *pb.SipScorePageToken,
) error {

	sipScores, err := fetch(pageToken, limit)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.HasMore = uint32(len(sipScores)) > limit-1
	if resp.HasMore {
		sipScores = sipScores[:limit-1]
	}

	if len(sipScores) == 0 {
		resp.PageToken = ""
		resp.SipScores = []*pb.SipScoreWithEntries{}
		resp.HasMore = false
		return nil
	}

	ids := make([]uint32, len(sipScores))
	for i, sip := range sipScores {
		ids[i] = sip.ID
	}

	if resp.HasMore {
		tokenStr, err := pagetoken.EncodePageToken(nextToken(sipScores[len(sipScores)-1]))
		if err != nil {
			return errno.ServerErr(errno.InternalServerError, "encode token failed")
		}
		resp.PageToken = tokenStr
	} else {
		resp.PageToken = ""
	}

	h := &sipScoreHydrate{}

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		h.collected, err = s.Dao.ListIsUserCollected(userID, constvar.CollectionSipScore, ids)
		return err
	})

	g.Go(func() error {
		var err error
		h.entries, err = s.Dao.BatchListSipScoreEntriesHottest(ids, constvar.DefaultPageSize)
		return err
	})

	if err := g.Wait(); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.SipScores = make([]*pb.SipScoreWithEntries, len(sipScores))
	for i, sipScore := range sipScores {
		resp.SipScores[i] = &pb.SipScoreWithEntries{
			Meta:    sipScoreModelToPB(sipScore, h.collected[sipScore.ID]),
			Entries: sipScoreEntriesModelToPB(h.entries[sipScore.ID]),
		}
	}

	return nil
}

func (s *PostService) listSipScoreNewest(ctx context.Context, userID uint32, pageToken *pb.SipScorePageToken, limit uint32, resp *pb.ListSipScoreResponse) error {
	return s.listSipScoreCommon(
		ctx, userID, pageToken, limit, resp,
		func(token *pb.SipScorePageToken, limit uint32) ([]*dao.SipScoreModel, error) {
			if token == nil {
				return s.Dao.ListSipScoreNewest(limit)
			}
			return s.Dao.ListSipScoreNewestWithCursor(
				token.GetId(),
				token.GetUpdatedAt().AsTime(),
				limit,
			)
		},
		func(last *dao.SipScoreModel) *pb.SipScorePageToken {
			return &pb.SipScorePageToken{
				Id:        last.ID,
				UpdatedAt: timestamppb.New(last.UpdatedAt),
				SortType:  constvar.SortByNewest,
			}
		},
	)
}

func (s *PostService) listSipScoreHottest(ctx context.Context, userID uint32, pageToken *pb.SipScorePageToken, limit uint32, resp *pb.ListSipScoreResponse) error {
	return s.listSipScoreCommon(
		ctx, userID, pageToken, limit, resp,
		func(token *pb.SipScorePageToken, limit uint32) ([]*dao.SipScoreModel, error) {
			if token == nil {
				return s.Dao.ListSipScoreHottest(limit)
			}
			return s.Dao.ListSipScoreHottestWithCursor(
				token.GetId(),
				token.GetParticipantCount(),
				limit,
			)
		},
		func(last *dao.SipScoreModel) *pb.SipScorePageToken {
			return &pb.SipScorePageToken{
				Id:               last.ID,
				ParticipantCount: last.ParticipantCount,
				SortType:         constvar.SortByHottest,
			}
		},
	)
}

func sipScoreModelToPB(m *dao.SipScoreModel, isCollected bool) *pb.SipScore {
	return &pb.SipScore{
		Id:               m.ID,
		CreatedAt:        timestamppb.New(m.CreatedAt),
		UpdatedAt:        timestamppb.New(m.UpdatedAt),
		CreatorId:        m.CreatorID,
		LastModifiedBy:   m.LastModifiedBy,
		EntryCount:       m.EntryCount,
		CollectCount:     m.CollectCount,
		ParticipantCount: m.ParticipantCount,
		Name:             m.Name,
		Description:      m.Description,
		CoverImg:         m.CoverImg,
		Domain:           m.Domain,
		Category:         m.Category,
		IsCollected:      isCollected,
		// Tags:  tag就暂时不获取了
	}
}

func sipScoreEntriesModelToPB(list []*dao.SipScoreEntryModel) []*pb.SipScoreEntry {
	if len(list) == 0 {
		return nil
	}

	res := make([]*pb.SipScoreEntry, len(list))
	for i, e := range list {
		res[i] = sipScoreEntryModelToPB(e)
	}
	return res
}
