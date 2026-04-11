package sipscore

import (
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateSipScoreEntries ... 创建榜单条目
// @Summary 创建榜单条目 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateSipScoreEntryRequest true "create_sip_score_entry_request"
// @Success 200 {object} IdsResponse
// @Router /sip-score/entries [post]
func (a *Api) CreateSipScoreEntries(c *gin.Context) {
	log.Info("Post CreateSipScoreEntries function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateSipScoreEntryRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)
	ok, err := model.Enforce(userID, constvar.SipScore, req.SipScoreID, constvar.Write)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}
	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}
	if ok = a.Dao.AllowN(userID, 3); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, nil, "Please try again later", GetLine())
		return
	}

	entries := make([]*pb.SipScoreEntryCreateInfo, 0, len(req.Entries))
	for _, entry := range req.Entries {
		entries = append(entries, &pb.SipScoreEntryCreateInfo{
			Name:        entry.Name,
			Description: entry.Description,
			CoverImg:    entry.CoverImg,
		})
	}

	createReq := pb.CreateSipScoreEntryRequest{
		SipScoreId: req.SipScoreID,
		Entries:    entries,
		CreatorId:  userID,
	}

	resp, err := client.PostClient.CreateSipScoreEntry(c.Request.Context(), &createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, &IdsResponse{IDs: resp.EntryIds})
}
