package comment

import (
	"context"
	pbf "forum-feed/proto"
	. "forum-gateway/handler"
	"forum-gateway/util"
	pb "forum-post/proto"
	"forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"forum/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Create ... 创建评论/从帖
// @Summary 创建评论/从帖 api
// @Tags comment
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param object body CreateRequest true "create_comment_request"
// @Success 200 {object} Comment
// @Router /comment [post]
func (a *Api) Create(c *gin.Context) {
	log.Info("Comment Create function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	var req CreateRequest
	if err := c.BindJSON(&req); err != nil {
		SendError(c, errno.ErrBind, nil, err.Error(), GetLine())
		return
	}

	if req.TypeName != constvar.SubPost && req.TypeName != constvar.FirstLevelComment && req.TypeName != constvar.SecondLevelComment {
		SendError(c, errno.ErrBadRequest, nil, "type_name must be "+constvar.SubPost+" or "+constvar.FirstLevelComment+" or "+constvar.SecondLevelComment, GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, constvar.Post, req.PostId, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	if ok := a.Dao.AllowN(userId, 5); !ok {
		SendError(c, errno.ErrExceededTrafficLimit, nil, "Please try again later", GetLine())
		return
	}

	createReq := pb.CreateCommentRequest{
		PostId:    req.PostId,
		TypeName:  req.TypeName,
		FatherId:  req.FatherId,
		Content:   req.Content,
		CreatorId: userId,
		ImgUrl:    req.ImgUrl,
	}

	createResp, err := client.PostClient.CreateComment(context.TODO(), &createReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	// 向 feed 发送请求
	pushReq := &pbf.PushRequest{
		Action: "评论",
		UserId: userId,
		Source: &pbf.Source{
			Id:       req.PostId,
			TypeName: createResp.TypeName,
			Name:     createResp.FatherContent,
		},
		TargetUserId: createResp.UserId,
		Content:      req.Content,
	}
	_, err = client.FeedClient.Push(context.TODO(), pushReq)

	resp := &Comment{
		Id:            createResp.Id,
		Content:       req.Content,
		TypeName:      req.TypeName,
		FatherId:      req.FatherId,
		CreateTime:    createResp.CreateTime,
		CreatorId:     userId,
		CreatorName:   createResp.CreatorName,
		CreatorAvatar: createResp.CreatorAvatar,
		LikeNum:       0,
		IsLiked:       false,
		ImgUrl:        req.ImgUrl,
	}

	if req.TypeName == constvar.SecondLevelComment {
		resp.BeRepliedContent = createResp.FatherContent
		resp.BeRepliedUserId = createResp.FatherUserId
	}

	SendResponse(c, err, resp)
}
