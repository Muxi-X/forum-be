package post

import (
	. "forum-gateway/handler"
	pb "forum-post/proto"
	"forum/client"
	"strconv"

	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
)

// SetQualityPost ... 设置帖子为精华帖
// @Summary 设置帖子为精华帖 api
// @Description 设置帖子为精华帖
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} Response
// @Router /post/set_quality/{post_id} [patch]
func (a *Api) SetQualityPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	req := &pb.Request{
		Id: uint32(id),
	}
	_, err = client.PostClient.SetQualityPost(c.Request.Context(), req)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
