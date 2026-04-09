package sipscore

import "forum-gateway/dao"

type Api struct {
	Dao dao.Interface
}

func New(i dao.Interface) *Api {
	api := new(Api)
	api.Dao = i
	return api
}

// ====================
// SipScore Domain
// ====================

// ---- model ----

// ---- request ----

type CreateSipScoreRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	CoverImg    string   `json:"cover_img" binding:"required"`
	Domain      string   `json:"domain" binding:"required"`
	Category    string   `json:"category" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
}

// =====================
// Other Common Response
// =====================

type IdResponse struct {
	Id uint32 `json:"id"`
}
