package authentication

import (
	"github.com/emptywe/trading_sim/entity"
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
	"github.com/emptywe/trading_sim/internal/storage/redis/simulator_cache/session_cache"
	"strings"
)

type Service struct {
	repo  simulator_repo.Authorization
	cache *session_cache.Cache
}

func NewAuthService(repo simulator_repo.Authorization, cache *session_cache.Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) CreateUser(request entity.SignUpRequest) (int, error) {
	var err error
	request.UserName = strings.ToLower(request.UserName)
	request.Password, err = generatePasswordHash(request.Password)
	if err != nil {
		return 0, err
	}
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
	var err error
	password, err = generatePasswordHash(password)
	if err != nil {
		return 0, err
	}
	// TODO: add logic
	return s.repo.CreateUser(entity.SignUpRequest{})
}

func (s *Service) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
