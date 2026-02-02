package service

import (
	"forum-feed/dao"
)

// FeedService ... 动态服务
type FeedService struct {
	Dao dao.Interface
}

func New(i dao.Interface) *FeedService {
	service := new(FeedService)
	service.Dao = i
	return service
}
