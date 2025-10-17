package dto

import "time"

type UpdateTagRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Title    string `json:"title"`
	StatusID uint   `json:"status_id" binding:"oneof=1 0"`
}

type ListTagsInput struct {
	Start    int    `form:"start"`
	Limit    int    `form:"limit"`
	OrderBy  string `form:"order_by"`
	Sort     string `form:"sort"`
	Filters  map[string]interface{}
	Title    string `form:"title"`
	ID       uint   `form:"id"`
	StatusID uint   `form:"status_id"`
}

type CreateTagRequest struct {
	Title    string `json:"title" binding:"required"`
	StatusID uint   `json:"status_id" binding:"oneof=1 0"`
}

type AddTagToArticle struct {
	TagIDs    string `json:"tag_ids" binding:"required"`
	ArticleID uint   `json:"article_id" binding:"required"`
}

type FetchedTag struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	StatusID  uint      `json:"status_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	Message  string      `json:"message" `
	Response interface{} `json:"response"`
}
type LoginReponse struct {
	Message     string `json:"message"`
	Token       string `json:"token"`
	LevelManage int8   `json:"level_manage"`
	UserId      uint   `json:"user_id`
}

type GetResponse struct {
	Message  string      `json:"message" `
	Response interface{} `json:"response"`
	Count    int64       `json:"count"`
}
type ErrorGetResponse struct {
	Message  string      `json:"message" `
	Response interface{} `json:"response"`
	Count    int64       `json:"count"`
	Err      interface{} `json:"err"`
}