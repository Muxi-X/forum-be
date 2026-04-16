package collection

import (
	pbf "forum-feed/proto"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateOrRemove ... 收藏/取消收藏帖子
// @Summary 收藏/取消收藏帖子 api
// @Tags collection
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRequest true "create_request"
// @Success 200 {object} handler.Response
// @Router /collection [post]
func (a *Api) CreateOrRemove(c *gin.Context) {
	log.Info("Collection CreateOrRemove function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)

	targetType := ""
	switch req.TargetType {
	case constvar.CollectionPost:
		targetType = constvar.Post
	case constvar.CollectionSipScore:
		targetType = constvar.SipScore
	default:
		SendError(c, errno.ErrBadRequest, nil, "invalid target_type", GetLine())
		return
	}

	ok, err := model.Enforce(userID, targetType, int(req.TargetID), constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	if ok := a.Dao.AllowN(userID, 2); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, nil, "Please try again later", GetLine())
		return
	}

	// 调用 RPC
	createReq := &pb.ToggleTargetRequest{
		UserId:     userID,
		TargetId:   uint32(req.TargetID),
		TargetType: req.TargetType,
	}

	resp, err := client.PostClient.CreateOrRemoveCollection(c.Request.Context(), createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	// feed 推送
	pushReq := &pbf.PushRequest{
		Action: "收藏",
		UserId: userID,
		Source: &pbf.Source{
			Id:       uint32(req.TargetID),
			TypeName: resp.TypeName,
			Name:     resp.Content,
		},
		TargetUserId: resp.UserId,
		Content:      "",
	}

	_, _ = client.FeedClient.Push(c.Request.Context(), pushReq)

	SendResponse(c, nil, nil)
}
