package service

import (
	"context"
	pb "forum-chat/proto"
	"forum/pkg/errno"
)

// UserList ... 获取用户列表
func (s *ChatService) UserList(ctx context.Context, req *pb.UserListRequest, resp *pb.UserListResponse) error {
	usersId, err := s.Dao.GetUserList(req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resp.UserIds = removeDistinct(usersId)

	return nil
}

func removeDistinct(usersId []uint32) []uint32 {
	newIds := make([]uint32, 0)
	idMap := make(map[uint32]bool)

	for _, id := range usersId {
		if idMap[id] {
			continue
		}
		newIds = append(newIds, id)
		idMap[id] = true
	}

	return newIds
}
