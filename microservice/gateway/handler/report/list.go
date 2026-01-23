package report

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/pkg/errno"
	"strconv"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// List ... list举报
// @Summary list举报 api
// @Tags report
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param last_id query int false "last_id"
// @Success 200 {object} ListResponse
// @Router /report/list [get]
func (a *Api) List(c *gin.Context) {
	log.Info("Report List function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	lastId, err := strconv.Atoi(c.DefaultQuery("last_id", "0"))
	if err != nil {
		SendError(c, errno.ErrQuery, nil, err.Error(), GetLine())
		return
	}

	listReq := &pb.ListReportRequest{
		LastId:     uint32(lastId),
		Offset:     uint32(page * limit),
		Limit:      uint32(limit),
		Pagination: limit != 0 || page != 0,
	}

	listResp, err := client.PostClient.ListReport(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, listResp, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, listResp, ListResponse{})
}
