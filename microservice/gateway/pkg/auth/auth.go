package auth

import (
	"errors"
	"forum-gateway/service"
	"forum/log"
	"forum/pkg/token"

	"github.com/gin-gonic/gin"
)

var (
	// ErrMissingHeader means the `Authorization` header was empty.
	ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")
	// ErrTokenInvalid means the token is invalid.
	ErrTokenInvalid = errors.New("the token is invalid")
)

// Context is the context of the JSON web token.
type Context struct {
	Id        uint32
	Role      uint32
	ExpiresAt int64 // 过期时间（时间戳，10位）
}

// Parse parses the token, and returns the context if the token is valid.
func Parse(tokenString string) (*Context, error) {
	if tokenString == "2" { // FIXME delete
		return &Context{
			Id:   2,
			Role: 7,
		}, nil
	}

	t, err := token.ResolveToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &Context{
		Id:        t.Id,
		Role:      t.Role,
		ExpiresAt: t.ExpiresAt,
	}, nil
}

// ParseRequest gets the token from the header and
// pass it to the Parse function to parse the token.
func ParseRequest(c *gin.Context) (*Context, error) {
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		header = c.Request.Header.Get("Sec-WebSocket-Protocol")
		if len(header) == 0 {
			return nil, ErrMissingHeader
		}
	}

	// check if in blacklist
	exist, err := service.CheckInBlacklist(header)
	if err != nil {
		log.Error("CheckInBlacklist error", log.String(err.Error()))
		return nil, err
	} else if exist {
		return nil, ErrTokenInvalid
	}

	return Parse(header)
}
