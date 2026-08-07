package auth

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

// ====================
// Common
// ====================

// ---- response ----

type EmptyResponse struct{}

type ExchangeFeedbackTokenRequest struct {
	TableIdentity string `json:"table_identity" binding:"required"`
}
