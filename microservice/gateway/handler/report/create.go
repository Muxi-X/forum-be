package report

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"github.com/Muxi-X/forum-be/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Create ... 举报帖子
// @Summary 举报帖子 api
// @Tags report
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRequest true "create_report_request"
// @Success 200 {object} handler.Response
// @Router /report [post]
func (a *Api) Create(c *gin.Context) {
	log.Info("Report Create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if req.TypeName != constvar.Post && req.TypeName != constvar.Comment {
		SendError(c, errno.ErrBadRequest, nil, "type_name must be "+constvar.Post+" or "+constvar.Comment, GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, req.TypeName, req.Id, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	if ok := a.Dao.AllowN(userId, 20); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, nil, "Please try again later", GetLine())
		return
	}

	createReq := pb.CreateReportRequest{
		UserId:   userId,
		Id:       req.Id,
		TypeName: req.TypeName,
		Category: req.Category,
		Cause:    req.Cause,
	}

	_, err = client.PostClient.CreateReport(context.TODO(), &createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
