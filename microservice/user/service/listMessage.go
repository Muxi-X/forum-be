package service

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"
)

// ListMessage ... 获取用户消息列表
func (s *UserService) ListMessage(_ context.Context, req *pb.ListMessageRequest, resp *pb.ListMessageResponse) error {
	logger.Info("UserService ListMessage")

	// DB 查询
	messages, err := s.Dao.ListMessage()
	if err != nil {
		return errno.ServerErr(errno.ErrRedis, err.Error())
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

	var wg sync.WaitGroup
	var mu sync.Mutex
	ch := make(chan struct{}, 20) // 限制并发数为 20
	resp.Messages = make([]*pb.Message, 0, len(messages))
	wg.Add(len(messages))

	for _, str := range messages {
		ch <- struct{}{}
		go func(str string) {
			defer func() {
				<-ch
				wg.Done()
			}()
			msg := &pb.Message{}
			if err := json.Unmarshal([]byte(str), msg); err != nil {
				logger.Error("json unmarshal failed", logger.String(err.Error()))
				return
			}

			senderId, _ := strconv.Atoi(msg.SendUserid)
			user, err := s.Dao.GetUser(uint32(senderId))
			if err != nil {
				logger.Error("GetUser failed", logger.String(err.Error()))
				return
			}
			msg.SenderName = user.Name
			msg.Avatar = user.Avatar

			mu.Lock()
			resp.Messages = append(resp.Messages, msg)
			mu.Unlock()
		}(str)
	}

	wg.Wait()
	return nil
}
