package feedback

import "forum-gateway/dao"

type Api struct {
	Dao dao.Interface
}

func New(i dao.Interface) *Api {
	api := new(Api)
	api.Dao = i
	return api
}

type CreateRequest struct {
	Category string `json:"category"`
	Content  string `json:"content" binding:"required"`
	Contact  string `json:"contact"`
	ImgURL   string `json:"img_url"`
} // @name FeedbackCreateRequest
