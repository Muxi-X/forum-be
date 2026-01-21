package report

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handle ... 处理举报
// @Summary 处理举报 api
// @Tags report
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body HandleRequest true "handle_report_request"
// @Success 200 {object} handler.Response
// @Router /report [put]
func (a *Api) Handle(c *gin.Context) {
	log.Info("Report Handle function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req HandleRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if req.Result != constvar.InvalidReport && req.Result != constvar.ValidReport {
		SendError(c, errno.ErrBadRequest, nil, "result must be "+constvar.InvalidReport+" or "+constvar.ValidReport, GetLine())
		return
	}

	handleReq := pb.HandleReportRequest{
		Id:     req.Id,
		Result: req.Result,
	}

	_, err := client.PostClient.HandleReport(context.TODO(), &handleReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, nil)
}
