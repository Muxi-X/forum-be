package service

import (
	"context"
	pb "forum-post/proto"
	"forum/util"
)

func (s *PostService) UpdateLastReadTime(_ context.Context, req *pb.UpdateLastReadRequest, resp *pb.Response) error {
	return s.Dao.UpdateLastRead(req.UserId, req.Category, util.GetCurrentTime())
}
