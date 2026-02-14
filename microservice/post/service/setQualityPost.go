package service

import (
	"context"
	pb "forum-post/proto"
)

func (s *PostService) SetQualityPost(_ context.Context, req *pb.Request, resp *pb.Response) error {
	return s.Dao.ChangeQualityPost(req.Id, true)
}
