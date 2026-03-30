package service

import (
	"forum/model"
	"forum/pkg/errno"
)

func resolveRoleByUserID(userID uint32) (string, error) {
	roles, err := model.GetRole(userID)
	if err != nil {
		return "", errno.ServerErr(errno.ErrCasbin, err.Error())
	}

	return model.SelectRole(roles), nil
}
