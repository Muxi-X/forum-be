package service

import (
	"context"

	"forum-user/dao"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// List ... 获取用户列表
func (s *UserService) List(ctx context.Context, req *pb.ListRequest, res *pb.ListResponse) error {
	logger.Info("UserService List")

	// 过滤条件
	filter := &dao.UserModel{Role: req.Role}

	// DB 查询
	list, err := s.Dao.ListUser(req.Offset, req.Limit, req.LastId, filter)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	var resList []*pb.User

	for _, item := range list {
		resList = append(resList, &pb.User{
			Id:     item.Id,
			Name:   item.Name,
			Avatar: item.Avatar,
			Role:   item.Role,
		})
	}

	res.Count = uint32(len(list))
	res.List = resList

	return nil
}
