package model

type Publisher struct {
	ID   uint   `json:"id" form:"id" binding:"omitempty"`
	UUID string `json:"uuid" form:"uuid" binding:"omitempty"`
	URL  string `json:"url" form:"url" binding:"omitempty,max=1000"`
}
