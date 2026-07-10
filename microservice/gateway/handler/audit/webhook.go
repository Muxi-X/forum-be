package audit

import (
	"encoding/json"

	"forum-gateway/dao"
	. "forum-gateway/handler"
	"forum-gateway/handler/post"
	"forum-gateway/handler/sipscore"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/api/request"
	"go.uber.org/zap"
)

// Webhook 审核结果回调
func (a *Api) Webhook(c *gin.Context) {
	var req *request.HookPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	pendingID := req.Data.Id

	switch req.Data.Status {
	case "Pass", "pass", "通过":
		a.handleItem(c, pendingID)

	case "Reject", "reject", "不通过":
		_, prefix, err := a.Dao.FindPending(pendingID)
		if err == nil {
			_ = a.Dao.DeletePending(prefix, pendingID)
		}
		log.Error("审核不通过",
			zap.Uint("pending_id", pendingID),
			zap.String("reason", req.Data.Msg),
		)
		SendResponse(c, nil, nil)

	default:
		log.Info("收到非终态回调",
			zap.Uint("pending_id", pendingID),
			zap.String("status", req.Data.Status),
		)
		SendResponse(c, nil, nil)
	}
}

// handleItem 审核通过后根据资源类型分发处理
func (a *Api) handleItem(c *gin.Context, pendingID uint) {
	pendingData, prefix, err := a.Dao.FindPending(pendingID)
	if err != nil {
		log.Error("获取待审核数据失败",
			zap.Uint("pending_id", pendingID),
			zap.Error(err),
		)
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	switch pendingData.ResourceType {
	case "post:create":
		a.createPost(c, pendingID, prefix, pendingData)
	case "post:update":
		a.updatePost(c, pendingID, prefix, pendingData)
	case "sipscore:create":
		a.createSipScore(c, pendingID, prefix, pendingData)
	case "sipscore:update":
		a.updateSipScore(c, pendingID, prefix, pendingData)
	default:
		log.Error("未知资源类型",
			zap.Uint("pending_id", pendingID),
			zap.String("resource_type", pendingData.ResourceType),
		)
		SendResponse(c, nil, nil)
	}
}

func (a *Api) createPost(c *gin.Context, pendingID uint, prefix string, d *dao.PendingData) {
	var req post.CreateRequest
	if err := json.Unmarshal(d.RawRequest, &req); err != nil {
		log.Error("反序列化请求失败", zap.Error(err))
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	createReq := &pb.CreatePostRequest{
		UserId:          d.UserId,
		Content:         req.Content,
		Domain:          req.Domain,
		Title:           req.Title,
		Category:        req.Category,
		ContentType:     req.ContentType,
		Tags:            req.Tags,
		CompiledContent: req.CompiledContent,
		Summary:         req.Summary,
	}

	resp, err := client.PostClient.CreatePost(c.Request.Context(), createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	_ = a.Dao.DeletePending(prefix, pendingID)
	SendResponse(c, nil, resp)
}

func (a *Api) updatePost(c *gin.Context, pendingID uint, prefix string, d *dao.PendingData) {
	var req post.UpdateInfoRequest
	if err := json.Unmarshal(d.RawRequest, &req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	updateReq := &pb.UpdatePostInfoRequest{
		Id:       req.Id,
		Content:  req.Content,
		Title:    req.Title,
		Domain:   req.Domain,
		UserId:   d.UserId,
		Category: req.Category,
		Tags:     req.Tags,
		Summary:  req.Summary,
	}

	_, err := client.PostClient.UpdatePostInfo(c.Request.Context(), updateReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	_ = a.Dao.DeletePending(prefix, pendingID)
	SendResponse(c, nil, nil)
}

func (a *Api) createSipScore(c *gin.Context, pendingID uint, prefix string, d *dao.PendingData) {
	var req sipscore.CreateSipScoreRequest
	if err := json.Unmarshal(d.RawRequest, &req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	createReq := &pb.CreateSipScoreRequest{
		Name:        req.Name,
		Description: req.Description,
		CoverImg:    req.CoverImg,
		Domain:      req.Domain,
		Category:    req.Category,
		Tags:        req.Tags,
		CreatorId:   d.UserId,
	}

	resp, err := client.PostClient.CreateSipScore(c.Request.Context(), createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	_ = a.Dao.DeletePending(prefix, pendingID)
	SendResponse(c, nil, resp)
}

func (a *Api) updateSipScore(c *gin.Context, pendingID uint, prefix string, d *dao.PendingData) {
	var req sipscore.UpdateSipScoreRequest
	if err := json.Unmarshal(d.RawRequest, &req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	updateReq := &pb.UpdateSipScoreInfoRequest{
		Id:          req.Id,
		Name:        req.Name,
		Description: *req.Description,
		CoverImg:    req.CoverImg,
		Domain:      req.Domain,
		Category:    req.Category,
		Tags:        req.Tags,
	}

	_, err := client.PostClient.UpdateSipScoreInfo(c.Request.Context(), updateReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	_ = a.Dao.DeletePending(prefix, pendingID)
	SendResponse(c, nil, nil)
}
