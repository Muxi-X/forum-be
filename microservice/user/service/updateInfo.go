package service

import (
	"context"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// UpdateInfo ... 更新用户信息
func (s *UserService) UpdateInfo(_ context.Context, req *pb.UpdateInfoRequest, _ *pb.Response) error {
	logger.Info("UserService UpdateInfo")

	user, err := s.Dao.GetUser(req.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	if user == nil {
		return errno.ServerErr(errno.ErrUserNotExisted, "")
	}

	user.Name = req.Info.Name
	user.Avatar = req.Info.AvatarUrl
	user.Signature = req.Info.Signature
	user.IsPublicCollectionAndLike = req.Info.IsPublicCollectionAndLike
	user.IsPublicFeed = req.Info.IsPublicFeed

	if err := user.Update(); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
