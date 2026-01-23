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

// Get ... 获取帖子
// @Summary 获取帖子 api
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Param post_id path int true "post_id"
// @Success 200 {object} GetPostResponse
// @Router /post/{post_id} [get]
func (a *Api) Get(c *gin.Context) {
	log.Info("Post Get function called.", zap.String("X-Request-Id", util.GetReqID(c)))

	id, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		SendError(c, errno.ErrPathParam, nil, err.Error(), GetLine())
		return
	}

	userId := c.MustGet("userId").(uint32)

	ok, err := model.Enforce(userId, constvar.Post, id, constvar.Read)
	if err != nil {
		SendError(c, errno.ErrCasbin, nil, err.Error(), GetLine())
		return
	}

	if !ok {
		SendError(c, errno.ErrPermissionDenied, nil, "权限不足", GetLine())
		return
	}

	getReq := &pb.Request{
		UserId: userId,
		Id:     uint32(id),
	}

	getResp, err := client.PostClient.GetPost(context.TODO(), getReq)
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	// 三次遍历分别处理 SubPost FirstLevelComment SecondLevelComment
	var subPosts []*SubPost
	subPostCommentsMap := make(map[uint32]*[]*Comment)

	commentsMap := make(map[uint32]*pb.CommentInfo, len(getResp.Comments))

	for _, comment := range getResp.Comments {
		commentsMap[comment.Id] = comment

		if comment.TypeName == constvar.SubPost {
			subPost := &SubPost{
				ImgUrl: comment.ImgUrl,
			}

			subPostCommentsMap[comment.Id] = &subPost.Comments

			setInfo(&subPost.info, comment)

			subPosts = append(subPosts, subPost)
		}
	}

	for _, comment := range getResp.Comments {
		if comment.TypeName == constvar.FirstLevelComment {
			var subPostComment = new(Comment)

			subPostComments, ok := subPostCommentsMap[comment.FatherId]

			// 没有找到父级，数据库错误
			if !ok {
				log.Error(errno.ErrDatabase.Error(), log.String(constvar.SubPost+":"+strconv.Itoa(int(comment.FatherId))+" not found"))
				continue
			}

			setInfo(&subPostComment.info, comment)

			*subPostComments = append(*subPostComments, subPostComment)

		} else if comment.TypeName == constvar.SecondLevelComment {

			beRepliedComment := commentsMap[comment.FatherId]
			commentReply := &Comment{
				BeRepliedId:       comment.FatherId,
				BeRepliedContent:  beRepliedComment.Content,
				BeRepliedUserName: beRepliedComment.CreatorName,
				BeRepliedUserId:   beRepliedComment.CreatorId,
			}

			setInfo(&commentReply.info, comment)

			subPostComments, ok := subPostCommentsMap[beRepliedComment.FatherId]

			// 没有找到父级，数据库错误
			if !ok {
				log.Error(errno.ErrDatabase.Error(), log.String(constvar.SubPost+":"+strconv.Itoa(int(beRepliedComment.FatherId))+" not found"))
				continue
			}

			*subPostComments = append(*subPostComments, commentReply)

		}
	}

	for _, subPost := range subPosts {
		subPost.CommentNum = uint32(len(subPost.Comments))
	}

	resp := GetPostResponse{
		CommentNum:      getResp.CommentNum,
		Title:           getResp.Title,
		Category:        getResp.Category,
		IsCollection:    getResp.IsCollection,
		ContentType:     getResp.ContentType,
		SubPosts:        subPosts,
		Tags:            getResp.Tags,
		CompiledContent: getResp.CompiledContent,
		Summary:         getResp.Summary,
		CollectionNum:   getResp.CollectionNum,
	}

	setInfo(&resp.info, getResp)

	SendResponse(c, nil, resp)
}
