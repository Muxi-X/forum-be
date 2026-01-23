package collection

import (
	"context"
	pbf "forum-feed/proto"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"

	"forum/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateOrRemove ... 收藏/取消收藏帖子
// @Summary 收藏/取消收藏帖子 api
// @Tags collection
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param post_id path int true "post_id"
// @Success 200 {object} handler.Response
// @Router /collection/{post_id} [post]
func (a *Api) CreateOrRemove(c *gin.Context) {
	log.Info("Collection CreateOrRemove function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	postId, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, constvar.Post, postId, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	if ok := a.Dao.AllowN(userId, 2); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, nil, "Please try again later", GetLine())
		return
	}

	createReq := pb.Request{
		UserId: userId,
		Id:     uint32(postId),
	}

	resp, err := client.PostClient.CreateOrRemoveCollection(context.TODO(), &createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	// 向 feed 发送请求
	pushReq := &pbf.PushRequest{
		Action: "收藏",
		UserId: userId,
		Source: &pbf.Source{
			Id:       uint32(postId),
			TypeName: resp.TypeName,
			Name:     resp.Content,
		},
		TargetUserId: resp.UserId,
		Content:      "",
	}
	_, err = client.FeedClient.Push(context.TODO(), pushReq)

	SendResponse(c, err, nil)
}
