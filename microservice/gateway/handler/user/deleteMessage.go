package user

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeletePrivateMessage ... 删除 private message
// @Summary 删除 个人 message api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} handler.Response
// @Router /user/private_message [delete]
func DeletePrivateMessage(c *gin.Context) {
	log.Info("User DeletePrivateMessage function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)
	messageId := c.Query("id")
	listReq := &pb.DeletePrivateMessageRequest{
		UserId: userId,
		Id:     messageId,
	}

	_, err := client.UserClient.DeletePrivateMessage(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
