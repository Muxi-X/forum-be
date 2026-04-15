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
// Common
// ====================

// ---- model ----

type userInfo struct {
	ID     uint32 `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// ---- response ----

type IdResponse struct {
	ID uint32 `json:"id"`
}

type IdsResponse struct {
	IDs []uint32 `json:"ids"`
}

type EmptyResponse struct{}

// ====================
// SipScore Domain
// ====================

// ---- model ----

type SipScore struct {
	ID               uint32    `json:"id"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
	Creator          *userInfo `json:"creator"`
	LastModifiedBy   *userInfo `json:"last_modified_by"`
	EntryCount       uint32    `json:"entry_count"`
	CollectCount     uint32    `json:"collect_count"`
	ParticipantCount uint32    `json:"participant_count"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CoverImg         string    `json:"cover_img"`
	Domain           string    `json:"domain"`
	Category         string    `json:"category"`
	Tags             []string  `json:"tags"`
	IsCollected      bool      `json:"is_collected"`
}

// ---- request ----

type CreateSipScoreRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	CoverImg    string   `json:"cover_img" binding:"required"`
	Domain      string   `json:"domain" binding:"required"`
	Category    string   `json:"category" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
}

type UpdateSipScoreRequest struct {
	Id          uint32   `json:"id" binding:"required"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	CoverImg    string   `json:"cover_img"`
	Domain      string   `json:"domain"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
}

// ---- response ----

type GetSipScoreResponse struct {
	SipScore *SipScore `json:"sip_score"`
}

// ====================
// SipScoreEntry Domain
// ====================

// ---- model ----

type SipScoreEntryCreateInfo struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CoverImg    string `json:"cover_img" binding:"required"`
}

type SipScoreEntry struct {
	ID               uint32    `json:"id"`
	SipScoreID       uint32    `json:"sip_score_id"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
	Creator          *userInfo `json:"creator"`
	LastModifiedBy   *userInfo `json:"last_modified_by"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CoverImg         string    `json:"cover_img"`
	ParticipantCount uint32    `json:"participant_num"`
	CommentCount     uint32    `json:"comment_num"`
	ScoreTotal       uint32    `json:"score_total"`
	ScoreAvg         uint32    `json:"score_avg"`
}

// ---- request ----

type CreateSipScoreEntryRequest struct {
	SipScoreID uint32                     `json:"sip_score_id" binding:"required"`
	Entries    []*SipScoreEntryCreateInfo `json:"entries" binding:"required,dive"`
}

type UpdateSipScoreEntryRequest struct {
	SipScoreID  uint32  `json:"sip_score_id" binding:"required"`
	EntryID     uint32  `json:"entry_id" binding:"required"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CoverImg    string  `json:"cover_img"`
}

// ---- response ----

type ListSipScoreEntriesResponse struct {
	Entries   []*SipScoreEntry `json:"entries"`
	PageToken string           `json:"page_token"`
	HasMore   bool             `json:"has_more"`
}
