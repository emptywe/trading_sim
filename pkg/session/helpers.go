package session

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateToken(reqToken string) error {
	if token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signTKey), nil
	}); err != nil || !token.Valid {
		return fmt.Errorf("invalid token: %v", err)
	}
	return nil
}

func ParseToken(claim jwt.Claims, token string) (interface{}, error) {
	exToken, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signTKey), nil
	})
	if err != nil {
		if !errors.Is(err, jwt.ErrTokenExpired) {
			return nil, err
		}
	}
	switch claim.(type) {
	case *TokenClaims:
		return exToken.Claims.(*TokenClaims), nil
	case *rTokenClaims:
		return exToken.Claims.(*rTokenClaims), nil
	default:
		return nil, errors.New("unknown claims type")
	}

}
