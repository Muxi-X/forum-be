package feedback

import (
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	"forum/log"
	"forum/pkg/errno"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Create ... 反馈与建议
// @Summary 反馈与建议 api
// @Tags feedback
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRequest true "create_feedback_request"
// @Success 200 {object} Response
// @Router /feedback [post]
func (a *Api) Create(c *gin.Context) {
	log.Info("Feedback Create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		SendError(c, errno.ErrBind, nil, "反馈内容不能为空", GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)
	profile, err := service.GetUserProfile(c.Request.Context(), userID, userID)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	extraRecord := map[string]any{}
	if category := strings.TrimSpace(req.Category); category != "" {
		extraRecord["问题类型"] = category
	}
	if profile.GetEmail() != "" {
		extraRecord["用户邮箱"] = profile.GetEmail()
	}

	studentID := strings.TrimSpace(profile.GetStudentId())
	if len(studentID) != 10 {
		SendError(c, errno.ErrBadRequest, nil, "student_id not found", GetLine())
		return
	}

	err = service.CreateFeedbackRecord(c.Request.Context(), service.FeedbackRecordRequest{
		StudentID:   studentID,
		Content:     content,
		ContactInfo: firstNonEmpty(req.Contact, profile.GetEmail()),
		ExtraRecord: extraRecord,
		ImageURLs:   feedbackImageURLs(req.ImgURL),
	})
	if err != nil {
		SendError(c, errno.InternalServerError, nil, err.Error(), GetLine())
		return
	}

	SendResponse(c, nil, nil)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func feedbackImageURLs(imgURL string) []string {
	imgURL = strings.TrimSpace(imgURL)
	if imgURL == "" {
		return nil
	}
	return []string{imgURL}
}
