package service

import (
	"context"
	"forum-user/dao"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// GetInfo ... 获取用户信息
func (s *UserService) GetInfo(_ context.Context, req *pb.GetInfoRequest, resp *pb.UserInfoResponse) error {
	logger.Info("UserService GetInfo")

	list, err := s.Dao.GetUserByIds(req.Ids)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	// 确保返回的用户信息顺序与请求的 ID 顺序一致
	idToUser := make(map[uint32]*dao.UserModel)
	for _, user := range list {
		idToUser[user.Id] = user
	}

	userInfos := make([]*pb.UserInfo, len(req.Ids))
	for i, id := range req.Ids {
		user, exists := idToUser[id]
		if !exists {
			// 如果用户不存在，这里返回一个未知的用户信息
			userInfos[i] = &pb.UserInfo{
				Id:                        0,
				Name:                      "Unknown",
				Email:                     "Unknown",
				AvatarUrl:                 "Unknown",
				Signature:                 "Unknown",
				IsPublicCollectionAndLike: false,
				IsPublicFeed:              false,
			}
			continue
		}
		userInfos[i] = &pb.UserInfo{
			Id:                        user.Id,
			Name:                      user.Name,
			Email:                     user.Email,
			AvatarUrl:                 user.Avatar,
			Signature:                 user.Signature,
			IsPublicCollectionAndLike: user.IsPublicCollectionAndLike,
			IsPublicFeed:              user.IsPublicFeed,
		}
	}

	resp.List = userInfos

	return nil
}
