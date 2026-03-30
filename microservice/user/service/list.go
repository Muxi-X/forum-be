package service

import (
	"context"
	"forum-user/dao"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// List ... 获取用户列表
func (s *UserService) List(_ context.Context, req *pb.ListRequest, resp *pb.ListResponse) error {
	logger.Info("UserService List")

	// 不再按 DB role 过滤，role 统一由 Casbin 获取
	filter := &dao.UserModel{}

	// DB 查询
	list, err := s.Dao.ListUser(req.Offset, req.Limit, req.LastId, filter)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	resList := make([]*pb.User, 0, len(list))

	for _, item := range list {
		role, err := resolveRoleByUserID(item.Id)
		if err != nil {
			return err
		}

		if req.Role != "" && role != req.Role {
			continue
		}

		resList = append(resList, &pb.User{
			Id:     item.Id,
			Name:   item.Name,
			Avatar: item.Avatar,
			Role:   role,
		})
	}

	resp.Count = uint32(len(resList))
	resp.List = resList

	return nil
}
