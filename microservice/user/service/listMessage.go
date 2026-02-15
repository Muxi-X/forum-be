package service

import (
	"context"
	"encoding/json"
	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"strconv"

	"github.com/samber/lo"
)

// ListMessage ... 获取用户消息列表
func (s *UserService) ListMessage(_ context.Context, req *pb.ListMessageRequest, resp *pb.ListMessageResponse) error {
	logger.Info("UserService ListMessage")

	// DB 查询
	messages, err := s.Dao.ListMessage()
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	if len(messages) < int(req.Offset) {
		return nil
	}

	if len(messages) > int(req.Offset+req.Limit) {
		messages = messages[req.Offset : req.Offset+req.Limit]
	} else {
		messages = messages[req.Offset:]
	}

	resp.Messages = messages

	return nil
}

func (s *UserService) ListPrivateMessage(_ context.Context, req *pb.ListMessageRequest, resp *pb.ListPrivateMessageResponse) error {
	logger.Info("UserService ListPrivateMessage")

	messages, err := s.Dao.ListPrivateMessage(req.UserId)
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
	}

	if len(messages) < int(req.Offset) {
		return nil
	}

	if len(messages) > int(req.Offset+req.Limit) {
		messages = messages[req.Offset : req.Offset+req.Limit]
	} else {
		messages = messages[req.Offset:]
	}

	ids := make([]uint32, 0, len(messages))
	resp.Messages = make([]*pb.Message, 0, len(messages))

	for _, str := range messages {
		msg := &pb.Message{}
		if err := json.Unmarshal([]byte(str), msg); err != nil {
			logger.Error("json unmarshal failed", logger.String(err.Error()))
			continue
		}
		senderId, _ := strconv.Atoi(msg.SendUserId)
		ids = append(ids, uint32(senderId))
		resp.Messages = append(resp.Messages, msg)
	}

	ids = lo.Uniq(ids)
	// 批量获取避免多次查询数据库
	user, err := s.Dao.BatchGetUser(ids)
	if err != nil {
		logger.Error("GetUser failed", logger.String(err.Error()))
		return nil
	}

	for _, msg := range resp.Messages {
		senderId, _ := strconv.Atoi(msg.SendUserId)
		msg.SenderName = user[uint32(senderId)].Name
		msg.Avatar = user[uint32(senderId)].Avatar
	}

	return nil
}
