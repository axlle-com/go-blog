package models

type UpdatePostRequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
