package chat

import (
	"context"
	pb "forum-chat/proto"
	. "forum-gateway/handler"

	"forum/client"

	"github.com/gin-gonic/gin"
)

// UserList ... 获取该用户的聊天列表
// @Summary 获取该用户的聊天列表
// @Tags chat
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {array} uint32 "成功返回用户ID列表"
// @Success 200 "示例: {\"code\":0,\"message\":ok,\"data\":[123,456,789]}"
// @Router /chat/userList [get]
func UserList(c *gin.Context) {
	userId := c.MustGet("userId").(uint32)

	req := &pb.UserListRequest{
		UserId: userId,
	}

	resp, err := client.ChatClient.UserList(context.Background(), req)
	if err != nil {
		SendError(c, err, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, resp.UserIds)
}
