package service

import (
	"context"
	"fmt"
	"time"

	pb "forum-user/proto"
	"forum/pkg/errno"
	"forum/pkg/token"
)

// Login 用户登录
func (s *UserService) Login(_ context.Context, req *pb.LoginRequest, resp *pb.LoginResponse) error {
	user, err := s.Dao.GetUserByEmail(req.Username)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	if user == nil {
		return errno.ServerErr(errno.ErrUserNotExisted, fmt.Sprintf("user %s not found", req.Username))
	}

	if !user.CheckPassword(req.Password) {
		return errno.ServerErr(errno.ErrPasswordIncorrect, "")
	}

	tokenStr, err := token.GenerateToken(&token.TokenPayload{
		Id:      user.Id,
		Role:    1,
		Expired: 24 * time.Hour * 7,
	})
	if err != nil {
		return errno.ServerErr(errno.ErrAuthToken, err.Error())
	}

	resp.Token = tokenStr
	return nil
}