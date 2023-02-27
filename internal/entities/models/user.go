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
	Current   float64 `json:"current" binding:"required"`
	Withdrawn float64 `json:"withdrawn" binding:"required"`
}
