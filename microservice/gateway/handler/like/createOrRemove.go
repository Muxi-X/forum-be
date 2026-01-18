package like

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"github.com/Muxi-X/forum-be/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateOrRemove ... 点赞/取消点赞
// @Summary 点赞/取消点赞 api
// @Tags like
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body Item true "like_request"
// @Success 200 {object} handler.Response
// @Router /like [post]
func (a *Api) CreateOrRemove(c *gin.Context) {
	log.Info("Like CreateOrRemove function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req Item
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, req.TypeName, req.TargetId, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	if ok := a.Dao.AllowN(userId, 1); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, nil, "Please try again later", GetLine())
		return
	}

	likeReq := &pb.LikeRequest{
		UserId: userId,
		Item: &pb.LikeItem{
			TargetId: req.TargetId,
			TypeName: req.TypeName,
		},
	}

	_, err = client.PostClient.CreateOrRemoveLike(context.TODO(), likeReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
