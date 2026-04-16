package collection

import (
	"forum-gateway/dao"
)

type Api struct {
	Dao dao.Interface
}

func New(i dao.Interface) *Api {
	api := new(Api)
	api.Dao = i
	return api
}

type CreateRequest struct {
	TargetID   uint64 `json:"target_id" binding:"required"`
	TargetType uint32 `json:"target_type" binding:"required"`
}
