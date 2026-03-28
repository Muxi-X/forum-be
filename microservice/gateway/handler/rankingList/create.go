package rankingList

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

// CreateRankingList ... 创建帖子
// @Summary 创建帖子 api
// @Tags rankingList
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRankingListRequest true "create_ranking_list_request"
// @Success 200 {object} IdResponse
// @Router /ranking-list [post]
func (a *Api) CreateRankingList(c *gin.Context) {
	log.Info("Post CreateRankingList function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateRankingListRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if req.Domain != constvar.NormalDomain && req.Domain != constvar.MuxiDomain {
		SendError(c, errno.ErrBadRequest, nil, "domain must be "+constvar.NormalDomain+" or "+constvar.MuxiDomain, GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userID, constvar.RankingList, req.Domain, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	if ok := a.Dao.AllowN(userID, 30); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, nil, "Please try again later", GetLine())
		return
	}

	createReq := pb.CreateRankingListRequest{
		CreatorId:   userID,
		Name:        req.Name,
		Description: req.Description,
		CoverImg:    req.CoverImg,
		Tags:        req.Tags,
		Domain:      req.Domain,
		Category:    req.Category,
	}

	resp, err := client.PostClient.CreateRankingList(c.Request.Context(), &createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, &IdResponse{Id: resp.Id})
}
