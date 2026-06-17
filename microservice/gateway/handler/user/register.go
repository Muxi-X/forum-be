package user

import (
	"forum-gateway/handler"
	"forum/pkg/errno"

	pb "forum-user/proto"
	"forum/client"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
// @Summary 用户注册
// @Tags auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "注册信息"
// @Success 200 {object} handler.Response
// @Router /api/v1/auth/register [post]
func Register(c *gin.Context) {
	var form RegisterRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		handler.SendError(c, errno.ErrBind, nil, err.Error(), "")
		return
	}

	_, err := client.UserClient.Register(c.Request.Context(), &pb.RegisterRequest{
		Email:    form.Username,
		Name:     form.Username,
		Password: form.Password,
	})
	if err != nil {
		handler.SendError(c, err, nil, err.Error(), "")
		return
	}

	handler.SendResponse(c, errno.OK, nil)
}
