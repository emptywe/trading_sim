package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func NewSession(username string, userId int) (*Session, string, error) {
	session := new(Session)
	sUUID, err := session.generateSessionTokens(username, userId)
	return session, sUUID, err
}

func (s *Session) generateSessionTokens(username string, userId int) (sUUID string, err error) {

	sUUID = uuid.New().String()
	s.Token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenEXP)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userId,
		username,
		sUUID,
	}).SignedString([]byte(signTKey))
	s.RToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, &rTokenClaims{
		jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		sUUID,
	}).SignedString([]byte(signTKey))

	return sUUID, err
}

func UpdateToken(claims *TokenClaims, rToken string) (string, error) {
	rClaims, err := ParseToken(&rTokenClaims{}, rToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %v", err)
	}
	if rClaims.(*rTokenClaims).Sid != claims.Sid {
		return "", errors.New("invalid session id")
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenEXP)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		claims.UserId,
		claims.UserName,
		claims.Sid,
	})

	token, err := newToken.SignedString([]byte(signTKey))
	if err != nil {
		return "", err
	}
	return token, nil
}
