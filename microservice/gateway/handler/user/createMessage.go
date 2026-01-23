package user

import (
	"context"
	"errors"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"
	"forum/pkg/errno"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ErrType = errors.New("type can only be (like/comment/collection/reply_comment)")

// CreateMessage ... 创建 public message
// @Summary 创建 公共 message api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateMessageRequest true "create_message_request"
// @Success 200 {object} handler.Response
// @Router /user/message [post]
func CreateMessage(c *gin.Context) {
	log.Info("User CreateMessage function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateMessageRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	listReq := &pb.CreateMessageRequest{
		Message: req.Message,
	}

	_, err := client.UserClient.CreateMessage(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}

// CreatePrivateMessage ... 创建 private message
// @Summary 创建 个人 message api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreatePrivateMessageRequest true "create_private_message_request"
// @Success 200 {object} handler.Response
// @Router /user/private_message [post]
func CreatePrivateMessage(c *gin.Context) {
	log.Info("User CreatePrivateMessage function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreatePrivateMessageRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}
	if req.Type != "comment" && req.Type != "collection" && req.Type != "like" && req.Type != "reply_comment" {
		SendError(c, errno.ErrBadRequest, nil, ErrType.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)
	listReq := &pb.CreatePrivateMessageRequest{
		ReceiveId: req.ReceiveUserid,
		SendId:    userId,
		Content:   req.Content,
		Type:      req.Type,
		PostId:    req.PostId,
		CommentId: req.CommentId,
	}

	_, err := client.UserClient.CreatePrivateMessage(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
