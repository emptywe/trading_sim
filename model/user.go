package model

type User struct {
	Uid       int     `json:"-" db:"id"`
	Email     string  `json:"email" binding:"required"`
	Password  string  `json:"password" binding:"required"`
	UserName  string  `json:"user_name" binding:"required"`
	FirstName string  `json:"firstName" `
	LastName  string  `json:"lastName" `
	Status    string  `json:"status" `
	Balance   float64 `json:"balance" `
}

type TUser struct {
	UserName string
	Balance  float64
}

type SignUpRequest struct {
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID int `json:"id"`
}

type SignInRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type SignInResponse struct {
	ID      int    `json:"id"`
	Session string `json:"session"`
	Token   string `json:"token"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
