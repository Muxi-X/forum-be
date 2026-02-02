package service

import (
	"context"
	"encoding/json"
	"forum-user/pkg/role"
	pbu "forum-user/proto"
	"forum/pkg/errno"
	"forum/util"

	"forum-feed/dao"
	pb "forum-feed/proto"
	logger "forum/log"

	"forum/client"
)

// Push ... 异步新增feed
func (s *FeedService) Push(_ context.Context, req *pb.PushRequest, _ *pb.Response) error {
	logger.Info("FeedService Push")

	getResp, err := client.UserClient.GetProfile(context.TODO(), &pbu.GetRequest{Id: req.UserId})
	if err != nil {
		return errno.ServerErr(errno.ErrRPC, err.Error())
	}

	feed := &dao.FeedModel{
		UserId:           req.UserId,
		UserName:         getResp.Name,
		UserAvatar:       getResp.Avatar,
		Action:           req.Action,
		SourceTypeName:   req.Source.TypeName,
		SourceObjectName: req.Source.Name,
		SourceObjectId:   req.Source.Id,
		TargetUserId:     req.TargetUserId,
		CreateTime:       util.GetCurrentTime(),
		Domain:           role.Role2Domain(getResp.Role),
	}

	msg, _ := json.Marshal(feed)

	if err := s.Dao.PublishMsg(msg); err != nil {
		return errno.ServerErr(errno.ErrPublishMsg, err.Error())
	}

	return nil
}
