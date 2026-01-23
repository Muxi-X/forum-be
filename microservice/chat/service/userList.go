package service

import (
	"context"
	pb "forum-chat/proto"

	"github.com/samber/lo"
)

// UserList ... 获取用户列表
func (s *ChatService) UserList(ctx context.Context, req *pb.UserListRequest, resp *pb.UserListResponse) error {
	userList, err := s.Dao.GetUserList(req.UserId, int(req.Limit), int(req.Page))
	if err != nil {
		return err
	}
	userList = lo.Uniq(userList)

	resp.UserLists, err = s.Dao.GetUserById(userList)
	if err != nil {
		return err
	}

	return nil
}
