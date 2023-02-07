package authentication

import (
	"errors"
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	tokenEXP = 5 * time.Minute
	signTKey = "sNKL213%md#4411jHKjHuh7*@1"
)

type AuthService struct {
	repo simulator_repo.Authorization
}

func NewAuthService(repo simulator_repo.Authorization) *AuthService {
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

func (s *AuthService) CreateSession(c *http.Cookie) (int, string, error) {
	id, sid, err := parseToken(c.Value)
	if err != nil {
		return 0, "", errors.New("invalid session token ")
	}
	return s.repo.CreateSession(c, id, sid)

}

func (s *AuthService) ValidateSession(c *http.Cookie) (int, bool) {

	id, sid, err := parseToken(c.Value)
	if err != nil {
		return 0, false
	}

	return s.repo.ValidateSession(c, id, sid)
}

func (s *AuthService) ExpireSession(c *http.Cookie) *http.Cookie {
	id, sid, err := parseToken(c.Value)
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
