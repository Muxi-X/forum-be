package comment

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"

	"github.com/Muxi-X/forum-be/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Delete ... 删除评论
// @Summary 删除评论 api
// @Tags comment
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param comment_id path int true "comment_id"
// @Success 200 {object} handler.Response
// @Router /comment/{comment_id} [delete]
func (a *Api) Delete(c *gin.Context) {
	log.Info("Comment Delete function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	ok, err := model.Enforce(userId, constvar.Comment, id, constvar.Write)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	deleteReq := &pb.DeleteItemRequest{
		Id:       uint32(id),
		TypeName: constvar.Comment,
		UserId:   userId,
	}

	_, err = client.PostClient.DeleteItem(context.TODO(), deleteReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
