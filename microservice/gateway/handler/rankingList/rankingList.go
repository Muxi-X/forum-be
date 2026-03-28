package rankingList

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
// RankingList Domain
// ====================

// ---- model ----

// ---- request ----

type CreateRankingListRequest struct {
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
