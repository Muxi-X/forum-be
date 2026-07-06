package audit

import (
	"forum-gateway/dao"
	. "forum-gateway/handler"
	"forum-gateway/handler/post"
	"forum/model"
	"forum/pkg/audit"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Post 送审帖子（先审后发）
// @Summary 创建帖子并送审 api
// @Tags audit
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body post.CreateRequest true "create_post_request"
// @Success 200 {object} Response
// @Router /audit/post [post]
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

	// title 前加 post 前缀方便回调区分
	req.Title = "post:" + req.Title
	pendingData := &dao.PendingPost{
		UserId:          userId,
		Content:         req.Content,
		Domain:          req.Domain,
		Title:           req.Title,
		Category:        req.Category,
		ContentType:     req.ContentType,
		Tags:            req.Tags,
		CompiledContent: req.CompiledContent,
		Summary:         req.Summary,
	}
	if err := a.Dao.SavePendingPost(pendingID, pendingData); err != nil {
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	if err := audit.AuditClient.SubmitToAudit(pendingID, strconv.Itoa(int(userId)), req.Title, req.Content, nil); err != nil {
		// 送审失败，清理暂存数据
		_ = a.Dao.DeletePendingPost(pendingID)
		SendError(c, errno.ErrAuditService, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, &post.IdResponse{Id: uint32(pendingID)})
}
