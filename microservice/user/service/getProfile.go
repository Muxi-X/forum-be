package service

import (
	"context"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// GetProfile ... 获取用户个人信息
func (s *UserService) GetProfile(_ context.Context, req *pb.GetRequest, resp *pb.UserProfile) error {
	logger.Info("UserService GetProfile")

	user, err := s.Dao.GetUser(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if user == nil {
		return errno.ServerErr(errno.ErrUserNotExisted, "")
	}

	resp.Id = user.Id
	resp.Name = user.Name
	resp.Avatar = user.Avatar
	resp.Email = user.Email
	resp.Role, err = resolveRoleByUserID(user.Id)
	if err != nil {
		return err
	}
	resp.Signature = user.Signature
	resp.IsPublicCollectionAndLike = user.IsPublicCollectionAndLike
	resp.IsPublicFeed = user.IsPublicFeed

	return nil
}
