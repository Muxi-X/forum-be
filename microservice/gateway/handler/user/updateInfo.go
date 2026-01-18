package user

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"
	"forum/pkg/errno"

	"github.com/Muxi-X/forum-be/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UpdateInfo ... 修改用户个人信息
// @Summary 修改用户个人信息 api
// @Tags user
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body UpdateInfoRequest  true "update_info_request"
// @Success 200 {object} handler.Response
// @Router /user [put]
func UpdateInfo(c *gin.Context) {
	log.Info("User UpdateInfo function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req UpdateInfoRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	// 构造请求给 getInfo
	updateInfoReq := &pb.UpdateInfoRequest{
		Id: userId,
		Info: &pb.UserInfo{
			Name:                      req.Name,
			AvatarUrl:                 req.AvatarURL,
			Signature:                 req.Signature,
			IsPublicCollectionAndLike: req.IsPublicCollectionAndLike,
			IsPublicFeed:              req.IsPublicFeed,
		},
	}

	_, err := client.UserClient.UpdateInfo(context.TODO(), updateInfoReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
