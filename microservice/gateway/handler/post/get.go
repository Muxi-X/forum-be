package post

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/pkg/errno"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// Get ... 获取帖子
// @Summary 获取帖子 api
// @Description
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param post_id path int true "post_id"
// @Success 200 {object} Post
// @Router /post/{post_id} [get]
func (a *Api) Get(c *gin.Context) {
	log.Info("Post Get function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	id, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	getReq := &pb.Request{
		UserId: userId,
		Id:     uint32(id),
	}

	res, err := service.PostClient.GetPost(context.TODO(), getReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, errno.OK, res)
}
