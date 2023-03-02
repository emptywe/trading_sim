package authentication

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"

	"github.com/emptywe/trading_sim/entity"
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
	"github.com/emptywe/trading_sim/internal/storage/redis/simulator_cache/session_cache"
)

type Service struct {
	repo  simulator_repo.Authorization
	cache *session_cache.Cache
}

func NewAuthService(repo simulator_repo.Authorization, cache *session_cache.Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) CreateUser(request entity.SignUpRequest) (int, error) {
	request.UserName = strings.ToLower(request.UserName)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	request.Password = fmt.Sprintf("%s", hashedPassword)
	return s.repo.CreateUser(request)
}

func (s *Service) ReadUser(request entity.SignInRequest) (entity.User, error) {
	request.UserName = strings.ToLower(request.UserName)
	if err := s.repo.ValidatePassword(request); err != nil {
		return entity.User{}, err
	}
	return s.repo.ReadUser(request.UserName)
}

func (s *Service) UpdateUser(email, userName, password string) (int, error) {

	// TODO: add logic
	return s.repo.CreateUser(entity.SignUpRequest{})
}

func (s *Service) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
