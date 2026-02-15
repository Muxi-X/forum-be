package post

import (
	"context"
	. "forum-gateway/handler"
	pb "forum-post/proto"
	"forum/client"
	"strconv"

	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
)

// SetQualityPost ... 设置帖子为精华帖
// @Summary 设置帖子为精华帖 api
// @Description 需要管理员权限才能设置帖子为精华帖
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} handler.Response
// @Router /post/set_quality/{post_id} [patch]
func (a *Api) SetQualityPost(c *gin.Context) {
	userId := c.MustGet("userId").(uint32)

	id, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	ok, err := model.HasRole(userId, constvar.NormalAdminRole)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	req := &pb.Request{
		Id: uint32(id),
	}
	_, err = client.PostClient.SetQualityPost(context.TODO(), req)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
