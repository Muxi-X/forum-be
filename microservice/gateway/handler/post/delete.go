package post

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// Delete ... 删除帖子
// @Summary 删除帖子 api
// @Description
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param post_id path int true "post_id"
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} handler.Response
// @Router /post/{post_id} [delete]
func (a *Api) Delete(c *gin.Context) {
	log.Info("Post Delete function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	id, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	ok, err := a.Dao.Enforce(userId, "type_id", constvar.Write)
	if err != nil {
		SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	deleteReq := &pb.Item{
		Id:     uint32(id),
		TypeId: constvar.Post,
	}

	// 发送请求
	_, err = service.PostClient.DeleteItem(context.Background(), deleteReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, errno.OK, nil)
}