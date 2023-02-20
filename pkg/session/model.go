package session

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	tokenEXP      = time.Minute * 15
	ExpireSession = time.Hour * 24
	signTKey      = "sNKL213%md#4411jHKjHuh7*@1"
)

type Session struct {
	Token  string `json:"token"`
	RToken string `json:"refresh_token"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Sid      string `json:"sid"`
}

type rTokenClaims struct {
	jwt.RegisteredClaims
	Sid string `json:"sid"`
}
