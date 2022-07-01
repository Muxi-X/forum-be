package service

import (
	"context"
	"forum-chat/dao"
	pb "forum-chat/proto"
	logger "forum/log"
	"forum/pkg/errno"
	"time"
)

// Create 发送消息
func (s *ChatService) Create(ctx context.Context, req *pb.CreateRequest, resp *pb.Response) error {
	logger.Info("CharService Create")

	data := &dao.ChatData{
		Message:  req.Message,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Receiver: req.TargetUserId,
		Sender:   req.UserId,
	}

	err := s.Dao.Create(data)

	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}
