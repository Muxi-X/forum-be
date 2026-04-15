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
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeleteSipScore ... 删除榜单
// @Summary 删除榜单 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param sip_score_id path int true "sip_score_id"
// @Success 200 {object} handler.Response
// @Router /sip-score/{sip_score_id} [delete]
func (a *Api) DeleteSipScore(c *gin.Context) {
	log.Info("Post DeleteSipScore function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userID := c.MustGet("userId").(uint32)

	sipScoreID, err := strconv.Atoi(c.Param("sip_score_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	ok, err := model.Enforce(userID, constvar.SipScore, sipScoreID, constvar.Write)
	if err != nil {
		SendError(c, errno.ErrCasbin, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, &EmptyResponse{}, "权限不足", GetLine())
		return
	}

	deleteReq := &pb.DeleteItemRequest{
		Id:       uint32(sipScoreID),
		TypeName: constvar.SipScore,
		UserId:   userID,
	}

	_, err = client.PostClient.DeleteItem(c.Request.Context(), deleteReq)
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}
	SendResponse(c, nil, &EmptyResponse{})
}
