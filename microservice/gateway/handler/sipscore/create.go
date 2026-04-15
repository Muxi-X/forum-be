package sipscore

import (
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateSipScore ... 创建榜单
// @Summary 创建榜单 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateSipScoreRequest true "create_sip_score_request"
// @Success 200 {object} IdResponse
// @Router /sip-score [post]
func (a *Api) CreateSipScore(c *gin.Context) {
	log.Info("Post CreateSipScore function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateSipScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	if req.Domain != constvar.NormalDomain && req.Domain != constvar.MuxiDomain {
		SendError(c, errno.ErrBadRequest, &EmptyResponse{}, "domain must be "+constvar.NormalDomain+" or "+constvar.MuxiDomain, GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)
	ok, err := model.Enforce(userID, constvar.SipScore, req.Domain, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, &EmptyResponse{}, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, &EmptyResponse{}, "权限不足", GetLine())
		return
	}

	if ok = a.Dao.AllowN(userID, 30); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, &EmptyResponse{}, "Please try again later", GetLine())
		return
	}

	createReq := pb.CreateSipScoreRequest{
		CreatorId:   userID,
		Name:        req.Name,
		Description: req.Description,
		CoverImg:    req.CoverImg,
		Tags:        req.Tags,
		Domain:      req.Domain,
		Category:    req.Category,
	}

	resp, err := client.PostClient.CreateSipScore(c.Request.Context(), &createReq)
	if err != nil {
		SendError(c, err, &EmptyResponse{}, "", GetLine())
		return
	}

	SendResponse(c, nil, &IdResponse{ID: resp.Id})
}
