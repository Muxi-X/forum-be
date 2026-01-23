package post

import (
	"context"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"

	"forum/client"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListMainPost ... 获取主帖
// @Summary list 主帖 api
// @Description 根据category or tag 获取主帖list
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Param last_id query int false "last_id"
// @Param category query string false "category"
// @Param filter query string false "filter"
// @Param search_content query string false "search_content"
// @Param tag query string false "tag"
// @Param domain path string true "normal -> 团队外; muxi -> 团队内"
// @Success 200 {object} ListMainPostResponse
// @Router /post/list/{domain} [get]
func (a *Api) ListMainPost(c *gin.Context) {
	log.Info("Post ListMainPost function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	userId := c.MustGet("userId").(uint32)

	domain := c.Param("domain")
	if domain != constvar.NormalDomain && domain != constvar.MuxiDomain {
		SendError(c, errno.ErrPathParam, nil, "domain not legal", GetLine())
		return
	}

	ok, err := model.Enforce(userId, constvar.Post, domain, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	category := c.DefaultQuery("category", "")

	filter := c.DefaultQuery("filter", "")

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
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

	searchContent := c.DefaultQuery("search_content", "")

	tag := c.DefaultQuery("tag", "")

	if domain == constvar.NormalDomain {
		ok, err := model.Enforce(userId, constvar.Post, constvar.MuxiDomain, constvar.Read)
		if err != nil {
			SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
			return
		}

		// 团队用户normal默认获取所有帖子
		if ok {
			domain = constvar.AllDomain
		}
	}

	listReq := &pb.ListMainPostRequest{
		UserId:        userId,
		Category:      category,
		Domain:        domain,
		LastId:        uint32(lastId),
		Offset:        uint32(page * limit),
		Limit:         uint32(limit),
		Pagination:    limit != 0 || page != 0,
		SearchContent: searchContent,
		Filter:        filter,
		Tag:           tag,
	}

	postResp, err := client.PostClient.ListMainPost(context.TODO(), listReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, postResp, ListMainPostResponse{})
}
