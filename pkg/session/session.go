package session

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

func NewSession(username string, userId int) (*Session, string, error) {
	session := new(Session)
	sUUID, err := session.generateSessionTokens(username, userId)
	session.Valid = true
	session.Established = time.Now()
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

func UpdateToken(token, rToken string) (string, error) {
	if err := ValidateToken(rToken); err != nil {
		return "", err
	}
	claims, err := ParseToken(TokenClaims{}, token)
	if err != nil {
		return "", err
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenEXP)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		claims.(*TokenClaims).UserId,
		claims.(*TokenClaims).UserName,
		claims.(*TokenClaims).Sid,
	})

	token, err = newToken.SignedString([]byte(signTKey))
	if err != nil {
		return "", err
	}
	return token, nil
}
