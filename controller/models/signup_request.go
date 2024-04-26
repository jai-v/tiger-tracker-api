package models

type SignupRequest struct {
	Username string `json:"userName" binding:"required,gt=4,max=20,alphanum"`
	Email    string `json:"email" binding:"required,email,max=70"`
	Password string `json:"password" binding:"required,gt=7,max=70"`
}
