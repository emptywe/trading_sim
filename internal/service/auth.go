package service

import (
	"errors"
	"fmt"
	"github.com/emptywe/trading_sim/internal/storage"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenEXP = 720 * time.Hour
	signTKey = "sNKL213%md#4411jHKjHuh7*@1"
)

type AuthService struct {
	repo storage.Authorization
}

func NewAuthService(repo storage.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int    `json:"user_id"`
	Sid    string `json:"sid"`
}

func (s *AuthService) CreateUser(email, userName, password string) (int, error) {
	var err error
	password, err = generatePasswordHash(password)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateUser(email, userName, password)
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("Bad password: %s", err)
		return "", err
	}

	return fmt.Sprintf("%s", hash), err
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, password)
	if err != nil {
		return "", err
	}
	sUUID := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenEXP)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.Uid,
		sUUID,
	})

	return token.SignedString([]byte(signTKey))
}

func (s *AuthService) ParseToken(accesToken string) (int, string, error) {

	token, err := jwt.ParseWithClaims(accesToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signTKey), nil
	})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type valid tokenClaims")
	}

	return claims.UserId, claims.Sid, nil
}

func (s *AuthService) CreateSession(c *http.Cookie) (int, string, error) {
	id, sid, err := s.ParseToken(c.Value)
	if err != nil {
		return 0, "", errors.New("invalid session token ")
	}
	return s.repo.CreateSession(c, id, sid)

}

func (s *AuthService) ValidateSession(c *http.Cookie) (int, bool) {

	id, sid, err := s.ParseToken(c.Value)
	if err != nil {
		return 0, false
	}

	return s.repo.ValidateSession(c, id, sid)
}

func (s *AuthService) ExpireSession(c *http.Cookie) *http.Cookie {
	id, sid, err := s.ParseToken(c.Value)
	if err != nil {
		return nil
	}

	return s.repo.ExpireSession(c, id, sid)

}

func (s *AuthService) DropUser(id int) error {
	return s.repo.DropUser(id)
}

func (s *AuthService) SessionGarbageCollector() error {
	return s.repo.SessionGarbageCollector()
}
