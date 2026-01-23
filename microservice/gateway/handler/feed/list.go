package feed

import (
	"context"
	pb "forum-feed/proto"
	. "forum-gateway/handler"
	"forum-gateway/util"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"

	"forum/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// List ... list feeds.
// @Summary list 用户的动态 api
// @Tags feed
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param user_id path int true "user_id"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param last_id query int false "last_id"
// @Success 200 {object} FeedListResponse
// @Router /feed/list/{user_id} [get]
func (a *Api) List(c *gin.Context) {
	log.Info("feed List function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	targetUserId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	if int(userId) != targetUserId {
		ok, err := model.Enforce(userId, constvar.Feed, targetUserId, constvar.Read)
		if err != nil {
			SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
			return
		}

		if !ok {
			SendResponse(c, errno.ErrPrivacyInfo, nil)
			return
		}
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
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

	listReq := &pb.ListRequest{
		LastId:       uint32(lastId),
		Offset:       uint32(page * limit),
		Limit:        uint32(limit),
		UserId:       userId,
		Pagination:   limit != 0 || page != 0,
		TargetUserId: uint32(targetUserId),
	}

	listResp, err := client.FeedClient.List(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, listResp, FeedListResponse{})
}
