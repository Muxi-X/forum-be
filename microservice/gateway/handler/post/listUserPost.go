package post

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/pkg/errno"
	"strconv"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListUserPost ... 获取用户发布的帖子
// @Summary list 用户发布的帖子 api
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param user_id path int true "user_id"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param last_id query int false "last_id"
// @Success 200 {object} PostPartInfoResponse
// @Router /post/published/{user_id} [get]
func (a *Api) ListUserPost(c *gin.Context) {
	log.Info("Post ListUserPost function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	targetUserId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	lastId, err := strconv.Atoi(c.DefaultQuery("last_id", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	listReq := &pb.ListPostPartInfoRequest{
		UserId:     uint32(targetUserId),
		LastId:     uint32(lastId),
		Offset:     uint32(page * limit),
		Limit:      uint32(limit),
		Pagination: limit != 0 || page != 0,
	}

	postResp, err := client.PostClient.ListUserCreatedPost(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, postResp, PostPartInfoResponse{})
}
