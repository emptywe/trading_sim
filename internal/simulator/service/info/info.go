package info

import (
	"github.com/emptywe/trading_sim/entity"
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
)

type Service struct {
	repo simulator_repo.Info
}

func NewInfoService(repo simulator_repo.Info) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetTopUsers() ([]entity.TUser, error) {
	return s.repo.GetTopUsers()
}

func (s *Service) GetAllCurrenciesUSD() ([]entity.CurrencyOutput, error) {
	return s.repo.GetAllCurrenciesUSD()
}
