package user

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"
	"strconv"

	"forum/client"

	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListMessage ... 获取 user message list
// @Summary list user message api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param limit query string false "limit"
// @Param page query string false "page"
// @Success 200 {object} ListMessageResponse
// @Router /user/message/list [get]
func ListMessage(c *gin.Context) {
	log.Info("User ListMessage function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	listReq := &pb.ListMessageRequest{
		UserId: userId,
		Limit:  uint32(limit),
		Offset: uint32(page * limit),
	}

	listResp, err := client.UserClient.ListMessage(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, listResp, ListMessageResponse{})
}

// ListPrivateMessage ... 获取 user private_message list
// @Summary list user private_message api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param limit query string false "limit"
// @Param page query string false "page"
// @Success 200 {object} ListMessageResponse
// @Router /user/private_message/list [get]
func ListPrivateMessage(c *gin.Context) {
	log.Info("User ListPrivateMessage function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	listReq := &pb.ListMessageRequest{
		UserId: userId,
		Limit:  uint32(limit),
		Offset: uint32(page * limit),
	}

	listResp, err := client.UserClient.ListPrivateMessage(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, listResp, ListMessageResponse{})
}
