package models

type User struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserBallance struct {
	Current   string `json:"current" binding:"required"`
	Withdrawn string `json:"withdrawn" binding:"required"`
}
