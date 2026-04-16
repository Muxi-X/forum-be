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

// DeleteSipScoreEntries ... 批量删除榜单条目
// @Summary 批量删除榜单条目 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body DeleteSipScoreEntriesRequest true "delete_sip_score_entries_request"
// @Success 200 {object} handler.Response
// @Router /sip-score/entries [delete]
func (a *Api) DeleteSipScoreEntries(c *gin.Context) {
	log.Info("Post DeleteSipScoreEntries function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req DeleteSipScoreEntriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, &EmptyResponse{}, err.Error(), GetLine())
		return
	}
	if req.SipScoreID == 0 {
		SendError(c, errno.ErrBadRequest, &EmptyResponse{}, "sip_score_id not legal", GetLine())
		return
	}

	if len(req.EntryIDs) == 0 {
		SendError(c, errno.ErrBadRequest, &EmptyResponse{}, "entry ids not legal", GetLine())
	}

	userID := c.MustGet("userId").(uint32)
	ok, err := model.Enforce(userID, constvar.SipScore, req.SipScoreID, constvar.Write)
	if err != nil {
		SendError(c, errno.ErrCasbin, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, &EmptyResponse{}, "权限不足", GetLine())
		return
	}

	if ok = a.Dao.AllowN(userID, 3); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, &EmptyResponse{}, "Please try again later", GetLine())
		return
	}

	deleteReq := &pb.DeleteSipScoreEntriesRequest{
		SipScoreId:      req.SipScoreID,
		SipScoreEntryId: req.EntryIDs,
		UserId:          userID,
	}

	_, err = client.PostClient.DeleteSipScoreEntries(c.Request.Context(), deleteReq)
	if err != nil {
		SendError(c, errno.ErrCasbin, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, &EmptyResponse{})
}
