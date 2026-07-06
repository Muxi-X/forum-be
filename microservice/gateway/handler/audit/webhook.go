package audit

import (
	. "forum-gateway/handler"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/pkg/errno"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/api/request"
	"go.uber.org/zap"
)

// Webhook 审核结果回调
// @Summary 审核结果回调接口
// @Tags audit
// @Accept application/json
// @Produce application/json
// @Param object body request.HookPayload true "审核回调数据"
// @Success 200 {object} Response
// @Router /audit/webhook [post]
func (a *Api) Webhook(c *gin.Context) {
	var req *request.HookPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	pendingID := req.Data.Id

	switch req.Data.Status {
	case "Pass", "pass", "通过":
		a.createItem(c, pendingID)

	case "Reject", "reject", "不通过":
		log.Error("审核不通过",
			zap.Uint("pending_id", pendingID),
			zap.String("reason", req.Data.Msg),
		)
		_ = a.Dao.DeletePendingPost(pendingID)
		SendResponse(c, nil, nil)

	default:
		log.Info("收到非终态回调",
			zap.Uint("pending_id", pendingID),
			zap.String("status", req.Data.Status),
		)
		SendResponse(c, nil, nil)
	}
}

// handlePass 审核通过后创建帖子
func (a *Api) createItem(c *gin.Context, pendingID uint) {
	pendingData, err := a.Dao.GetPendingPost(pendingID)
	if err != nil {
		log.Error("获取待审核数据失败",
			zap.Uint("pending_id", pendingID),
			zap.Error(err),
		)
		SendError(c, errno.ErrDatabase, nil, err.Error(), GetLine())
		return
	}

	parts := strings.SplitN(pendingData.Title, ":", 2)
	if len(parts) != 2 {
		parts = append(parts, "")
	}

	switch parts[0] {
	case "post":
		createReq := &pb.CreatePostRequest{
			UserId:          pendingData.UserId,
			Content:         pendingData.Content,
			Domain:          pendingData.Domain,
			Title:           parts[1],
			Category:        pendingData.Category,
			ContentType:     pendingData.ContentType,
			Tags:            pendingData.Tags,
			CompiledContent: pendingData.CompiledContent,
			Summary:         pendingData.Summary,
		}

		resp, err := client.PostClient.CreatePost(c.Request.Context(), createReq)
		if err != nil {
			log.Error("创建帖子失败",
				zap.Uint("pending_id", pendingID),
				zap.Error(err),
			)
			SendError(c, err, nil, "", GetLine())
			return
		}

		_ = a.Dao.DeletePendingPost(pendingID)

		log.Info("审核通过，帖子已创建",
			zap.Uint("pending_id", pendingID),
			zap.Uint32("post_id", resp.Id),
		)

		SendResponse(c, nil, resp)

	default:
		SendResponse(c, nil, nil)
	}
}
