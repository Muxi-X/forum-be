package user

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"
	"forum/pkg/errno"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TeamLogin ... 团队登录
// @Summary 团队登录 api
// @Description login the team-forum
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param object body TeamLoginRequest true "login_request"
// @Success 200 {object} TeamLoginResponse
// @Router /auth/login/team [post]
func TeamLogin(c *gin.Context) {
	log.Info("User TeamLogin function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	// 从前端获取 oauth_code
	var req TeamLoginRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	// 构造请求给 login
	loginReq := &pb.TeamLoginRequest{
		OauthCode: req.OauthCode,
	}

	loginResp, err := client.UserClient.TeamLogin(context.TODO(), loginReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, loginResp, TeamLoginResponse{})
}
