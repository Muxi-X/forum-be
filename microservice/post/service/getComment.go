package service

import (
	"context"
	"forum-post/dao"
	pb "forum-post/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"strconv"
)

func (s *PostService) GetComment(ctx context.Context, req *pb.Request, resp *pb.CommentInfo) error {
	logger.Info("PostService GetComment")

	comment, err := s.Dao.GetCommentInfo(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if comment == nil {
		return errno.NotFoundErr(errno.ErrItemNotFound, "comment-"+strconv.Itoa(int(req.Id)))
	}

	likeNum, err := s.Dao.GetLikedNum(dao.Item{
		Id:     req.Id,
		TypeId: 2,
	})
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	resp.LikeNum = comment.LikeNum
	if likeNum != 0 {
		resp.LikeNum = uint32(likeNum)
	}
	resp.Id = comment.Id
	resp.Content = comment.Content
	resp.CreateTime = comment.CreateTime
	resp.CreatorId = comment.CreatorId
	resp.CreatorAvatar = comment.CreatorAvatar
	resp.CreatorName = comment.CreatorName

	return nil
}
