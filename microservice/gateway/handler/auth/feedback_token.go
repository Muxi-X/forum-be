package auth

import (
	. "forum-gateway/handler"
	"forum-gateway/service"
	"forum-gateway/util"
	pb "forum-user/proto"
	"forum/log"
	"forum/pkg/errno"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ExchangeFeedbackToken ... 为反馈中台换取访问令牌
// @Summary 换取反馈中台 Token
// @Description 使用当前 forum-be 登录态，为指定的反馈表换取短期反馈中台访问令牌。forum-faq 仅授予只读权限。
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "forum-be 用户令牌"
// @Param object body ExchangeFeedbackTokenRequest true "feedback_token_request"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Failure 500 {object} Response
// @Router /auth/feedback/token [post]
func (a *Api) ExchangeFeedbackToken(c *gin.Context) {
	log.Info("Auth ExchangeFeedbackToken function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req ExchangeFeedbackTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)

	if ok := a.Dao.AllowN(userID, 3); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, &EmptyResponse{}, "Please try again later", GetLine())
		return
	}

	resp, err := client.UserClient.GetUserIdentity(c.Request.Context(), &pb.GetUserIdentityRequest{UserId: userID})
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}

	studentID := resp.GetStudentId()
	if studentID == "" {
		SendError(c, errno.ErrItemNotFound, &EmptyResponse{}, "student_id not found", GetLine())
		return
	}

	feedbackToken, err := service.ExchangeFeedbackToken(c.Request.Context(), studentID, req.TableIdentity)
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}

	SendResponse(c, nil, feedbackToken)
}
