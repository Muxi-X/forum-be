package service

import (
	"context"
	"encoding/base64"
	"forum-post/dao"
	pb "forum-post/proto"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// todo 评分系统需要完成

func (s *PostService) ListSipScoreEntry(_ context.Context, req *pb.ListSipScoreEntryRequest, resp *pb.ListSipScoreEntryResponse) error {
	if req.GetSipScoreId() <= 0 {
		return errno.ServerErr(errno.ErrBadRequest, "sip_score_id required")
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = constvar.DefaultPageSize
	} else if pageSize > constvar.MaxPageSize {
		pageSize = constvar.MaxPageSize
	}
	limit := pageSize + 1

	pageToken, err := decodeSipScoreEntryPageToken(req.GetPageToken())
	if err != nil {
		return errno.ServerErr(errno.ErrBadRequest, "invalid page token")
	}

	switch req.GetSortType() {
	case constvar.SortByNewest:
		return s.listSipScoreEntryNewest(req.GetSipScoreId(), pageToken, limit, resp)
	case constvar.SortByHottest:
		return s.listSipScoreEntryHottest(req.GetSipScoreId(), pageToken, limit, resp)
	case constvar.SortByHighestScore:
		return s.listSipScoreEntryHighestScore(req.GetSipScoreId(), pageToken, limit, resp)
	case constvar.SortByLowestScore:
		return s.listSipScoreEntryLowestScore(req.GetSipScoreId(), pageToken, limit, resp)
	default:
		return errno.ServerErr(errno.ErrBadRequest, "sort type not legal")
	}
}

// NOTE: 很好 AI 成功的给我原来重复度高的代码抽象了
func (s *PostService) listSipScoreEntriesCommon(
	sipScoreID uint32,
	pageToken *pb.SipScoreEntryPageToken,
	limit uint32,
	resp *pb.ListSipScoreEntryResponse,
	fetch func(token *pb.SipScoreEntryPageToken, limit uint32) ([]*dao.SipScoreEntryModel, error),
	nextToken func(last *dao.SipScoreEntryModel) *pb.SipScoreEntryPageToken,
) error {
	entries, err := fetch(pageToken, limit)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.HasMore = false
	if uint32(len(entries)) > limit-1 {
		resp.HasMore = true
		entries = entries[:limit-1]
	}

	// 如果没有更多了，直接返回空的 page token 和 entries
	if len(entries) == 0 {
		resp.PageToken = ""
		resp.Entries = nil
		return nil
	}

	if resp.HasMore {
		tokenStr, err := encodeSipScoreEntryPageToken(nextToken(entries[len(entries)-1]))
		if err != nil {
			return errno.ServerErr(errno.InternalServerError, "failed to encode page token")
		}
		resp.PageToken = tokenStr
	} else {
		resp.PageToken = ""
	}

	resp.Entries = make([]*pb.SipScoreEntry, len(entries))
	for i, e := range entries {
		resp.Entries[i] = sipScoreEntryModelToPB(e)
	}
	return nil
}

func (s *PostService) listSipScoreEntryNewest(sipScoreID uint32, pageToken *pb.SipScoreEntryPageToken, limit uint32, resp *pb.ListSipScoreEntryResponse) error {
	return s.listSipScoreEntriesCommon(sipScoreID, pageToken, limit, resp,
		func(token *pb.SipScoreEntryPageToken, limit uint32) ([]*dao.SipScoreEntryModel, error) {
			if token == nil {
				return s.Dao.ListSipScoreEntriesNewest(sipScoreID, limit)
			}
			return s.Dao.ListSipScoreEntriesNewestWithCursor(
				sipScoreID,
				token.GetEntryId(),
				token.GetUpdatedAt().AsTime(),
				limit,
			)
		},
		func(last *dao.SipScoreEntryModel) *pb.SipScoreEntryPageToken {
			return &pb.SipScoreEntryPageToken{
				EntryId:   last.ID,
				UpdatedAt: timestamppb.New(last.UpdatedAt),
			}
		},
	)
}

func (s *PostService) listSipScoreEntryHottest(sipScoreID uint32, pageToken *pb.SipScoreEntryPageToken, limit uint32, resp *pb.ListSipScoreEntryResponse) error {
	return s.listSipScoreEntriesCommon(sipScoreID, pageToken, limit, resp,
		func(token *pb.SipScoreEntryPageToken, limit uint32) ([]*dao.SipScoreEntryModel, error) {
			if token == nil {
				return s.Dao.ListSipScoreEntriesHottest(sipScoreID, limit)
			}
			return s.Dao.ListSipScoreEntriesHottestWithCursor(
				sipScoreID,
				token.GetEntryId(),
				token.GetParticipantCount(),
				limit,
			)
		},
		func(last *dao.SipScoreEntryModel) *pb.SipScoreEntryPageToken {
			return &pb.SipScoreEntryPageToken{
				EntryId:          last.ID,
				ParticipantCount: last.ParticipantCount,
			}
		},
	)
}

func (s *PostService) listSipScoreEntryHighestScore(sipScoreID uint32, pageToken *pb.SipScoreEntryPageToken, limit uint32, resp *pb.ListSipScoreEntryResponse) error {
	return s.listSipScoreEntriesCommon(sipScoreID, pageToken, limit, resp,
		func(token *pb.SipScoreEntryPageToken, limit uint32) ([]*dao.SipScoreEntryModel, error) {
			if token == nil {
				return s.Dao.ListSipScoreEntriesHighestScore(sipScoreID, limit)
			}
			return s.Dao.ListSipScoreEntriesHighestScoreWithCursor(
				sipScoreID,
				token.GetEntryId(),
				token.GetScoreAvg(),
				limit,
			)
		},
		func(last *dao.SipScoreEntryModel) *pb.SipScoreEntryPageToken {
			return &pb.SipScoreEntryPageToken{
				EntryId:  last.ID,
				ScoreAvg: last.ScoreAvg,
			}
		},
	)
}

func (s *PostService) listSipScoreEntryLowestScore(sipScoreID uint32, pageToken *pb.SipScoreEntryPageToken, limit uint32, resp *pb.ListSipScoreEntryResponse) error {
	return s.listSipScoreEntriesCommon(sipScoreID, pageToken, limit, resp,
		func(token *pb.SipScoreEntryPageToken, limit uint32) ([]*dao.SipScoreEntryModel, error) {
			if token == nil {
				return s.Dao.ListSipScoreEntriesLowestScore(sipScoreID, limit)
			}
			return s.Dao.ListSipScoreEntriesLowestScoreWithCursor(
				sipScoreID,
				token.GetEntryId(),
				token.GetScoreAvg(),
				limit,
			)
		},
		func(last *dao.SipScoreEntryModel) *pb.SipScoreEntryPageToken {
			return &pb.SipScoreEntryPageToken{
				EntryId:  last.ID,
				ScoreAvg: last.ScoreAvg,
			}
		},
	)
}

func sipScoreEntryModelToPB(entry *dao.SipScoreEntryModel) *pb.SipScoreEntry {
	return &pb.SipScoreEntry{
		Id:               entry.ID,
		SipScoreId:       entry.SipScoreID,
		CreatedAt:        timestamppb.New(entry.CreatedAt),
		UpdatedAt:        timestamppb.New(entry.UpdatedAt),
		CreatorId:        entry.CreatorID,
		LastModifiedBy:   entry.LastModifiedBy,
		Name:             entry.Name,
		Description:      entry.Description,
		CoverImg:         entry.CoverImg,
		ParticipantCount: entry.ParticipantCount,
		CommentCount:     entry.CommentCount,
		ScoreTotal:       entry.ScoreTotal,
		ScoreAvg:         entry.ScoreAvg,
	}
}

func encodeSipScoreEntryPageToken(token *pb.SipScoreEntryPageToken) (string, error) {
	if token == nil {
		return "", nil
	}

	data, err := proto.Marshal(token)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func decodeSipScoreEntryPageToken(pageToken string) (*pb.SipScoreEntryPageToken, error) {
	if pageToken == "" {
		return nil, nil
	}

	data, err := base64.RawURLEncoding.DecodeString(pageToken)
	if err != nil {
		return nil, err
	}

	token := &pb.SipScoreEntryPageToken{}
	if err := proto.Unmarshal(data, token); err != nil {
		return nil, err
	}

	return token, nil
}
