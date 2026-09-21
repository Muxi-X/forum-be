package service

import (
	"context"
	"strings"

	pb "forum-user/proto"
	logger "forum/log"
	"forum/pkg/errno"

	"go.uber.org/zap"
)

func (s *UserService) GetUserIdentity(_ context.Context, req *pb.GetUserIdentityRequest, resp *pb.GetUserIdentityResponse) error {
	logger.Info("UserService GetUserIdentity: request received")

	if req == nil || req.GetUserId() == 0 {
		logger.Info("UserService GetUserIdentity: invalid user_id")
		return errno.ServerErr(errno.ErrBadRequest, "user_id is required")
	}

	userID := req.GetUserId()
	logger.Info("UserService GetUserIdentity: querying user", zap.Uint32("user_id", userID))

	user, err := s.Dao.GetUser(userID)
	if err != nil {
		logger.Error("UserService GetUserIdentity: database query failed",
			zap.Uint32("user_id", userID),
			logger.String(err.Error()),
		)
		return errno.ServerErr(errno.ErrDatabase, "get user error")
	}
	if user == nil {
		logger.Info("UserService GetUserIdentity: user not found or disabled",
			zap.Uint32("user_id", userID),
		)
		return errno.NotFoundErr(errno.ErrItemNotFound, "user not found")
	}

	// 一般不会为空
	studentID := strings.TrimSpace(user.StudentId)
	logger.Info("UserService GetUserIdentity: user loaded",
		zap.Uint32("user_id", user.Id),
		zap.Bool("deleted", user.Re),
		zap.Bool("student_id_present", studentID != ""),
		zap.Int("student_id_length", len(studentID)),
	)
	if studentID == "" {
		logger.Info("UserService GetUserIdentity: student_id is empty",
			zap.Uint32("user_id", userID),
		)
		return errno.ServerErr(errno.ErrItemNotFound, "student_id not found")
	}

	resp.UserId = user.Id
	resp.StudentId = studentID
	logger.Info("UserService GetUserIdentity: identity resolved",
		zap.Uint32("user_id", user.Id),
	)

	return nil
}
