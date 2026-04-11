package sipscore

import (
	. "forum-gateway/handler"
	pb "forum-post/proto"
	"forum/client"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// UpdateSipScore ... 修改榜单信息
// @Summary 修改榜单信息 api
// @Tags sipscore
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body UpdateSipScoreRequest  true "update_sip_score_request"
// @Success 200 {object} handler.Response
// @Router /sip-score [put]
func (a *Api) UpdateSipScore(c *gin.Context) {
	log.Info("SipScore Update function called.")

	var req UpdateSipScoreRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	userID := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userID, constvar.SipScore, req.Id, constvar.Write)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	if ok = a.Dao.AllowN(userID, 3); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, nil, "Please try again later", GetLine())
		return
	}

	paths := buildPaths(&req)
	if len(paths) == 0 {
		SendError(c, errno.ErrBadRequest, nil, "no fields to update", GetLine())
		return
	}

	updateReq := pb.UpdateSipScoreInfoRequest{
		Id:             req.Id,
		Name:           req.Name,
		CoverImg:       req.CoverImg,
		Domain:         req.Domain,
		Category:       req.Category,
		Tags:           req.Tags,
		LastModifiedBy: userID,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
	}

	if req.Description != nil {
		updateReq.Description = *req.Description
	}

	_, err = client.PostClient.UpdateSipScoreInfo(c.Request.Context(), &updateReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendResponse(c, nil, &EmptyResponse{})
}

func buildPaths(req *UpdateSipScoreRequest) []string {
	paths := make([]string, 0)
	if req.Name != "" {
		paths = append(paths, "name")
	}
	if req.Description != nil {
		paths = append(paths, "description")
	}
	if req.CoverImg != "" {
		paths = append(paths, "cover_img")
	}
	if req.Domain != "" {
		paths = append(paths, "domain")
	}
	if req.Category != "" {
		paths = append(paths, "category")
	}
	if len(req.Tags) > 0 {
		paths = append(paths, "tags")
	}

	return paths
}
