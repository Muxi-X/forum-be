package service

import (
	"strings"

	"forum/pkg/errno"
)

const sensitiveContentMessage = "含敏感词无法发表"

func databaseErr(err error) error {
	if isSensitiveContentError(err) {
		return errno.ServerErr(errno.ErrSensitiveContent, "")
	}
	return errno.ServerErr(errno.ErrDatabase, err.Error())
}

func isSensitiveContentError(err error) bool {
	return err != nil && strings.Contains(err.Error(), sensitiveContentMessage)
}
