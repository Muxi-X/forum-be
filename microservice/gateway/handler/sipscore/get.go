package sipscore

import (
	"fmt"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	userpb "forum-user/proto"
	"forum/client"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// todo 根据 前端 zhy 的建议：
// todo 加一个 sse 通知前端，如果榜单 元信息更新了
// todo 直接全量刷新

// GetSipScore ... 获取榜单元数据
// @Summary 获取榜单元数据 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param post_id path int true "sip_score_id"
// @Success 200 {object} GetSipScoreResponse
// @Router /sip-score/{sip_score_id} [get]
func (a *Api) GetSipScore(c *gin.Context) {
	log.Info("GetSipScore called.", zap.String("X-Request-Id", util.GetReqID(c)))
	fmt.Println(1)
	sipScoreID, err := strconv.Atoi(c.Param("sip_score_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}
	if sipScoreID <= 0 {
		SendError(c, errno.ErrPathParam, nil, "sip_score_id must be positive", GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userID, constvar.SipScore, sipScoreID, constvar.Read)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}
	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	req := &pb.Request{
		UserId: userID,
		Id:     uint32(sipScoreID),
	}
	fmt.Println(2)

	sipScoreResp, err := client.PostClient.GetSipScore(c.Request.Context(), req)
	fmt.Println(2.5)
	fmt.Println(err)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}
	fmt.Println(2.6)

	userIDs := []uint32{
		sipScoreResp.CreatorId,
		sipScoreResp.LastModifiedBy,
	}

	// 确保顺序一致
	fmt.Println(3)

	userResp, err := client.UserClient.GetInfo(c.Request.Context(), &userpb.GetInfoRequest{Ids: userIDs})
	if err != nil {
		log.Error("Failed to get user info",
			zap.String("X-Request-Id", util.GetReqID(c)),
			zap.Uint32s("userIds", userIDs),
			zap.Error(err),
		)
	}

	var creator *userInfo
	var lastModifiedBy *userInfo

	build := func(u *userpb.UserInfo, fallbackID uint32) *userInfo {
		if u == nil {
			return &userInfo{ID: fallbackID}
		}
		return &userInfo{
			ID:     fallbackID, // 用请求 ID 保证一定有
			Name:   u.Name,
			Avatar: u.AvatarUrl,
		}
	}

	if userResp != nil {
		userList := userResp.GetList()

		if len(userList) > 0 {
			creator = build(userList[0], userIDs[0])
		}
		if len(userList) > 1 {
			lastModifiedBy = build(userList[1], userIDs[1])
		}
	}
	fmt.Println(4)

	sipScore := &SipScore{
		ID:               sipScoreResp.Id,
		CreatedAt:        sipScoreResp.CreatedAt.AsTime().Format(time.DateTime),
		UpdatedAt:        sipScoreResp.UpdatedAt.AsTime().Format(time.DateTime),
		Creator:          creator,
		LastModifiedBy:   lastModifiedBy,
		EntryCount:       sipScoreResp.EntryCount,
		CollectCount:     sipScoreResp.CollectCount,
		ParticipantCount: sipScoreResp.ParticipantCount,
		Name:             sipScoreResp.Name,
		Description:      sipScoreResp.Description,
		CoverImg:         sipScoreResp.CoverImg,
		Domain:           sipScoreResp.Domain,
		Category:         sipScoreResp.Category,
		Tags:             sipScoreResp.Tags,
		IsCollected:      sipScoreResp.IsCollected,
	}

	resp := &GetSipScoreResponse{SipScore: sipScore}

	SendResponse(c, nil, resp)
}
