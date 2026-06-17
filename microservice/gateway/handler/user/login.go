package user

import (
	"forum-gateway/handler"
	"forum/pkg/errno"

	pb "forum-user/proto"
	"forum/client"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
// @Summary 用户登录
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登录信息"
// @Success 200 {object} handler.Response
// @Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var form LoginRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		handler.SendError(c, errno.ErrBind, nil, err.Error(), "")
		return
	}

	resp, err := client.UserClient.Login(c.Request.Context(), &pb.LoginRequest{
		Username: form.Username,
		Password: form.Password,
	})
	if err != nil {
		handler.SendError(c, err, nil, err.Error(), "")
		return
	}

	handler.SendResponse(c, errno.OK, gin.H{"token": resp.Token})
}