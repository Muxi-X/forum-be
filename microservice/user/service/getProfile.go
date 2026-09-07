package service

import (
	"context"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

const placeholderStudentID = "0000000000"

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
	resp.StudentId = user.StudentId
	if len(resp.StudentId) != 10 {
		resp.StudentId = placeholderStudentID
	}
	resp.Role, err = resolveRoleByUserID(user.Id)
	if err != nil {
		return err
	}
	resp.Signature = user.Signature
	resp.IsPublicCollectionAndLike = user.IsPublicCollectionAndLike
	resp.IsPublicFeed = user.IsPublicFeed
	resp.FollowingCount, err = s.Dao.CountFollowing(user.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	resp.FollowerCount, err = s.Dao.CountFollowers(user.Id)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	if req.ViewerId != 0 && req.ViewerId != user.Id {
		resp.IsFollowing, err = s.Dao.IsFollowing(req.ViewerId, user.Id)
		if err != nil {
			return errno.ServerErr(errno.ErrDatabase, err.Error())
		}
	}

	//若查询者非本人，则清除敏感信息
	if req.ViewerId != user.Id {
		resp.Email = ""
		resp.StudentId = placeholderStudentID
	}

	return nil
}
