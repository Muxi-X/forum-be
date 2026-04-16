package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"strconv"
	"time"
)

func (s *PostService) GetComment(_ context.Context, req *pb.Request, resp *pb.CommentInfo) error {
	logger.Info("PostService GetPostComment")

	comment, err := s.Dao.GetCommentInfo(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if comment == nil {
		return errno.NotFoundErr(errno.ErrItemNotFound, "comment-"+strconv.Itoa(int(req.Id)))
	}

	likeNum, err := s.Dao.GetLikedNum(dao.Item{
		Id:       req.Id,
		TypeName: constvar.Comment,
	})
	if err != nil {
		logger.Error(errno.ErrRedis.Error(), logger.String(err.Error()))
	}

	resp.LikeNum = comment.LikeNum
	if likeNum != 0 {
		resp.LikeNum = uint32(likeNum)
	}
	resp.TypeName = comment.TypeName
	resp.Id = comment.Id
	resp.Content = comment.Content
	resp.Time = comment.CreateTime.Format(time.DateTime)
	resp.CreatorId = comment.CreatorId
	resp.CreatorAvatar = comment.CreatorAvatar
	resp.CreatorName = comment.CreatorName
	resp.ImgUrl = comment.ImgUrl

	return nil
}
