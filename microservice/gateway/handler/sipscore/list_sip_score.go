package sipscore

import (
	. "forum-gateway/handler"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/pkg/errno"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ListSipScores ... 批量获取榜单
// @Summary list 批量获取榜单 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param sort_type query int false "sort_type"
// @Param page_size query int false "page_size"
// @Param page_token query string false "page_token"
// @Success 200 {object} ListSipScoresResponse
// @Router /sip-score/list [get]
// todo 权限检测
func (a *Api) ListSipScores(c *gin.Context) {
	log.Info("Post ListSipScores function called.")

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil {
		SendError(c, errno.ErrQuery, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	sortType, err := strconv.Atoi(c.DefaultQuery("sort_type", "1"))
	if err != nil {
		SendError(c, errno.ErrQuery, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	pageToken := c.DefaultQuery("page_token", "")

	userID := c.MustGet("userId").(uint32)

	listReq := &pb.ListSipScoreRequest{
		PageToken: pageToken,
		PageSize:  uint32(pageSize),
		SortType:  uint32(sortType),
		UserId:    userID,
	}

	listResp, err := client.PostClient.ListSipScore(c.Request.Context(), listReq)
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}

	rpcSipScores := listResp.GetSipScores()
	httpSipScores := make([]*SipScoreWithEntries, len(rpcSipScores))

	for i, rpcSipScore := range rpcSipScores {
		if rpcSipScore == nil {
			rpcSipScore = &pb.SipScoreWithEntries{}
		}

		rpcEntries := rpcSipScore.Entries
		if len(rpcEntries) == 0 {
			rpcEntries = []*pb.SipScoreEntry{}
		}

		httpEntries := make([]*SipScoreEntry, len(rpcEntries))
		for j, rpcEntry := range rpcEntries {
			if rpcEntry == nil {
				rpcEntry = &pb.SipScoreEntry{}
			}
			httpEntries[j] = &SipScoreEntry{
				ID:               rpcEntry.Id,
				SipScoreID:       rpcEntry.SipScoreId,
				CreatedAt:        rpcEntry.CreatedAt.AsTime().Format(time.DateTime),
				UpdatedAt:        rpcEntry.UpdatedAt.AsTime().Format(time.DateTime),
				Creator:          &userInfo{ID: rpcEntry.CreatorId},
				LastModifiedBy:   &userInfo{ID: rpcEntry.LastModifiedBy},
				Name:             rpcEntry.Name,
				Description:      rpcEntry.Description,
				CoverImg:         rpcEntry.CoverImg,
				ParticipantCount: rpcEntry.ParticipantCount,
				CommentCount:     rpcEntry.CommentCount,
				ScoreTotal:       rpcEntry.ScoreTotal,
				ScoreAvg:         rpcEntry.ScoreAvg,
			}
		}

		meta := rpcSipScore.GetMeta()
		if meta == nil {
			meta = &pb.SipScore{}
		}

		httpSipScore := &SipScore{
			ID:               meta.Id,
			CreatedAt:        meta.CreatedAt.AsTime().Format(time.DateTime),
			UpdatedAt:        meta.UpdatedAt.AsTime().Format(time.DateTime),
			Creator:          &userInfo{ID: meta.CreatorId},
			LastModifiedBy:   &userInfo{ID: meta.LastModifiedBy},
			EntryCount:       meta.EntryCount,
			CollectCount:     meta.CollectCount,
			ParticipantCount: meta.ParticipantCount,
			Name:             meta.Name,
			Description:      meta.Description,
			CoverImg:         meta.CoverImg,
			Domain:           meta.Domain,
			Category:         meta.Category,
			Tags:             meta.Tags,
			IsCollected:      meta.IsCollected,
		}

		httpSipScores[i] = &SipScoreWithEntries{
			SipScore: httpSipScore,
			Entries:  httpEntries,
		}
	}

	resp := &ListSipScoresResponse{
		SipScores: httpSipScores,
		PageToken: listResp.PageToken,
		HasMore:   listResp.HasMore,
	}

	SendResponse(c, nil, resp)
}
