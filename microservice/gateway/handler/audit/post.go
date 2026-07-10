package audit

import (
	"encoding/json"
	"strconv"

	"forum-gateway/dao"
	. "forum-gateway/handler"
	"forum-gateway/handler/post"
	"forum/model"
	"forum/pkg/audit"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
)

// Post 送审帖子（先审后发）
func (a *Api) Post(c *gin.Context) {
	var req post.CreateRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if req.Domain != constvar.NormalDomain && req.Domain != constvar.MuxiDomain {
		SendError(c, errno.ErrBadRequest, nil, "domain must be "+constvar.NormalDomain+" or "+constvar.MuxiDomain, GetLine())
		return
	}

	if req.ContentType != "md" && req.ContentType != "rtf" {
		SendError(c, errno.ErrBadRequest, nil, "content_type must be md or rtf", GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, constvar.Post, req.Domain, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}
	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	if ok := a.Dao.AllowN(userId, 30); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, nil, "Please try again later", GetLine())
		return
	}

	pendingID, err := a.Dao.NextPendingID()
	if err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	rawReq, _ := json.Marshal(req)
	pendingData := &dao.PendingData{
		ResourceType: "post:create",
		UserId:       userId,
		RawRequest:   rawReq,
	}
	if err := a.Dao.SavePending(dao.PendingPrefixPost, pendingID, pendingData); err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	if err := audit.AuditClient.SubmitToAudit(pendingID, strconv.Itoa(int(userId)), req.Title, req.Content, nil); err != nil {
		_ = a.Dao.DeletePending(dao.PendingPrefixPost, pendingID)
		SendError(c, errno.ErrAuditService, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, pendingID)
}

// PostUpdate 送审帖子修改
func (a *Api) PostUpdate(c *gin.Context) {
	var req post.UpdateInfoRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, constvar.Post, req.Id, constvar.Write)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}
	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	if ok := a.Dao.AllowN(userId, 3); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, nil, "Please try again later", GetLine())
		return
	}

	pendingID, err := a.Dao.NextPendingID()
	if err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	rawReq, _ := json.Marshal(req)
	pendingData := &dao.PendingData{
		ResourceType: "post:update",
		UserId:       userId,
		RawRequest:   rawReq,
	}
	if err := a.Dao.SavePending(dao.PendingPrefixPost, pendingID, pendingData); err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	if err := audit.AuditClient.SubmitToAudit(pendingID, strconv.Itoa(int(userId)), req.Title, req.Content, nil); err != nil {
		_ = a.Dao.DeletePending(dao.PendingPrefixPost, pendingID)
		SendError(c, errno.ErrAuditService, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, pendingID)
}
