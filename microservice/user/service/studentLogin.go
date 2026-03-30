package service

import (
	"context"
	"errors"
	"strings"

	"forum-user/dao"
	"forum-user/pkg/auth"
	pb "forum-user/proto"
	"forum-user/util"
	logger "forum/log"
	"forum/model"
	"forum/pkg/constvar"
	"forum/pkg/errno"
	"forum/pkg/token"
)

// StudentLogin handles the CCNU student login state machine and issues a forum token on success.
func (s *UserService) StudentLogin(_ context.Context, req *pb.StudentLoginRequest, resp *pb.LoginResponse) error {
	logger.Info("UserService StudentLogin")

	if strings.TrimSpace(req.GetAction()) == "" || strings.TrimSpace(req.GetAction()) == "start" {
		if strings.TrimSpace(req.GetStudentId()) == "" || strings.TrimSpace(req.GetPassword()) == "" {
			return errno.ServerErr(errno.ErrBind, "student_id and password are required")
		}
	}

	loginResult, err := auth.HandleStudentLogin(req)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrStudentCredentialInvalid):
			return errno.ServerErr(errno.ErrPasswordIncorrect, err.Error())
		case errors.Is(err, auth.ErrStudentLoginBadRequest), errors.Is(err, auth.ErrStudentLoginSessionExpired):
			return errno.ServerErr(errno.ErrBadRequest, err.Error())
		default:
			return errno.ServerErr(errno.ErrGetUserInfo, err.Error())
		}
	}

	fillStudentLoginState(resp, loginResult)
	if loginResult.State.Status != "logged_in" {
		return nil
	}

	Token, err := s.issueStudentToken(loginResult.StudentID, loginResult.Password)
	if err != nil {
		return err
	}
	resp.Token = Token
	return nil
}

func fillStudentLoginState(resp *pb.LoginResponse, result *auth.StudentLoginResult) {
	resp.SessionId = result.State.SessionID
	resp.Status = result.State.Status
	resp.Message = result.State.Message
	resp.CaptchaImageBase64 = result.State.CaptchaImageBase64
	resp.AvailableSecondAuthMethods = result.State.AvailableSecondAuthMethods
	resp.CurrentSecondAuthMethod = result.State.CurrentSecondAuthMethod
	resp.SecondAuthSmsTarget = result.State.SecondAuthSMSTarget
	resp.SecondAuthEmailTarget = result.State.SecondAuthEmailTarget
}

func (s *UserService) issueStudentToken(studentID string, password string) (string, error) {
	// 查询是否存在用户
	// 若已存在则刷新密码；否则自动注册。
	// NOTE: 这里仍沿用 forum 的本地用户体系与 token 发放逻辑。
	user, err := s.Dao.GetUserByStudentId(studentID)
	if err != nil {
		return "", err
	}

	// 如果用户为空
	if user == nil {
		info := &dao.RegisterInfo{
			StudentId: studentID,
			Password:  password,
			Role:      constvar.NormalRole,
			Name:      studentID,
		}

		// 用户未注册，自动注册
		if err := s.Dao.RegisterUser(info); err != nil {
			return "", errno.ServerErr(errno.ErrDatabase, err.Error())
		}

		// 注册后重新查询
		user, err = s.Dao.GetUserByStudentId(studentID)
		if err != nil {
			return "", errno.ServerErr(errno.ErrDatabase, err.Error())
		}

		if err := s.Dao.AddPublicPolicy(user.Role, user.Id); err != nil {
			return "", errno.ServerErr(errno.ErrCasbin, err.Error())
		}

		if err := model.AddRole("user", user.Id, constvar.NormalRole); err != nil {
			return "", errno.ServerErr(errno.ErrCasbin, err.Error())
		}
	} else {
		// 更新用户密码
		err := s.Dao.UpdatePassword(user.Id, password)
		if err != nil {
			return "", err
		}
	}

	user.Role, err = resolveRoleByUserID(user.Id)
	if err != nil {
		return "", err
	}

	// 根据权限生成 token
	role := uint32(constvar.Normal)
	if user.Role == constvar.NormalAdminRole || user.Role == constvar.MuxiAdminRole {
		role = constvar.Admin
	} else if user.Role == constvar.SuperAdminRole {
		role = constvar.SuperAdmin
	}

	// 生成 auth token
	Token, err := token.GenerateToken(&token.TokenPayload{
		Id:      user.Id,
		Role:    role,
		Expired: util.GetExpiredTime(),
	})
	if err != nil {
		return "", errno.ServerErr(errno.ErrAuthToken, err.Error())
	}

	return Token, nil
}
