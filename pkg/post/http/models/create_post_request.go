package models

type CreatePostRequest struct {
	Title   string `json:"title" form:"title" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
}
