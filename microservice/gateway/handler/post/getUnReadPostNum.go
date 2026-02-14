package post

import (
	"context"
	. "forum-gateway/handler"
	pb "forum-post/proto"
	"forum/client"

	"github.com/gin-gonic/gin"
)

// GetUnReadPostNum ... 获取每个板块未读帖子数
// @Summary 获取每个板块未读帖子数 api
// @Description 注意要先调用这个接口获取后再调用/list/:domain接口！！！
// @Tags post
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "token 用户令牌"
// @Success 200 {object} UnReadNumResponse
// @Router /post/unread_num [get]
func (a *Api) GetUnReadPostNum(c *gin.Context) {
	userId := c.MustGet("userId").(uint32)

	resp, err := client.PostClient.GetUnReadPostNum(context.TODO(), &pb.Request{
		UserId: userId,
	})
	if err != nil {
		SendError(c, err, nil, "", GetLine())
		return
	}

	SendMicroServiceResponse(c, nil, resp, UnReadNumResponse{})
}
