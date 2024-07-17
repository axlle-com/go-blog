package models

type AuthInput struct {
	Email    string `form:"email" json:"email" binding:"required,email" validate:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=32" validate:"required,min=3,max=32"`
}
