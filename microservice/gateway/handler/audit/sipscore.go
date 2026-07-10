package audit

import (
	"encoding/json"
	"strconv"

	"forum-gateway/dao"
	. "forum-gateway/handler"
	"forum-gateway/handler/sipscore"
	"forum/pkg/audit"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
)

// CreateSipScore 送审评分表创建
func (a *Api) CreateSipScore(c *gin.Context) {
	var req sipscore.CreateSipScoreRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	pendingID, err := a.Dao.NextPendingID()
	if err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	rawReq, _ := json.Marshal(req)
	pendingData := &dao.PendingData{
		ResourceType: "sipscore:create",
		UserId:       userId,
		RawRequest:   rawReq,
	}
	if err := a.Dao.SavePending(dao.PendingPrefixSipScore, pendingID, pendingData); err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	if err := audit.AuditClient.SubmitToAudit(pendingID, strconv.Itoa(int(userId)), req.Name, req.Description, []string{req.CoverImg}); err != nil {
		_ = a.Dao.DeletePending(dao.PendingPrefixSipScore, pendingID)
		SendError(c, errno.ErrAuditService, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, pendingID)
}

// UpdateSipScore 送审评分表修改
func (a *Api) UpdateSipScore(c *gin.Context) {
	var req sipscore.UpdateSipScoreRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)
	content := ""
	if req.Description != nil {
		content = *req.Description
	}

	pendingID, err := a.Dao.NextPendingID()
	if err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	rawReq, _ := json.Marshal(req)
	pendingData := &dao.PendingData{
		ResourceType: "sipscore:update",
		UserId:       userId,
		RawRequest:   rawReq,
	}
	if err := a.Dao.SavePending(dao.PendingPrefixSipScore, pendingID, pendingData); err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	if err := audit.AuditClient.SubmitToAudit(pendingID, strconv.Itoa(int(userId)), req.Name, content, []string{req.CoverImg}); err != nil {
		_ = a.Dao.DeletePending(dao.PendingPrefixSipScore, pendingID)
		SendError(c, errno.ErrAuditService, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, pendingID)
}
