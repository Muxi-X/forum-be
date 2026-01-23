package chat

import (
	"context"
	pb "forum-chat/proto"
	. "forum-gateway/handler"
	"strconv"

	"forum/client"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
)

type UserStatus struct {
	Id     uint32
	Name   string
	Avatar string
}

// UserList ... 获取该用户的聊天列表
// @Summary 获取该用户的聊天列表
// @Description 获取该用户的聊天列表，包括每个聊天对象的名字和图片
// @Tags chat
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} []UserStatus "成功返回用户列表"
// @Router /chat/userList [get]
func UserList(c *gin.Context) {
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

	req := &pb.UserListRequest{
		UserId: userId,
		Page:   uint32(page),
		Limit:  uint32(limit),
	}

	resp, err := client.ChatClient.UserList(context.Background(), req)
	if err != nil {
		SendError(c, err, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, resp.UserLists)
}
