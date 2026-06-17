package service

import (
	"context"
	"fmt"

	pb "forum-user/proto"
	"forum-user/dao"
	"forum/pkg/errno"
)

// Register 用户注册
func (s *UserService) Register(_ context.Context, req *pb.RegisterRequest, resp *pb.RegisterResponse) error {
	existing, err := s.Dao.GetUserByEmail(req.Email)
	if err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}
	if existing != nil {
		return errno.ServerErr(errno.ErrBind, fmt.Sprintf("user %s already exists", req.Email))
	}

	if err := s.Dao.RegisterUser(&dao.RegisterInfo{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     "normal",
	}); err != nil {
		return errno.ServerErr(errno.ErrDatabase, err.Error())
	}

	return nil
}