package basket

import (
	"github.com/emptywe/trading_sim/entity"
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
)

type Service struct {
	repo simulator_repo.Basket
}

func NewBasketService(repo simulator_repo.Basket) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCurrency(name string) (entity.Currency, error) {
	return s.repo.GetCurrency(name)
}

func (s *Service) CreateBasket(id int, c1, c2 string, v float64) (int, error) {
	return s.repo.CreateBasket(id, c1, c2, v)
}

func (s *Service) GetBasket(id int, c string) (entity.Basket, error) {
	return s.repo.GetBasket(id, c)
}

func (s *Service) GetAllCurrenciesUSD() ([]entity.CurrencyOutput, error) {
	return s.repo.GetAllCurrenciesUSD()
}

func (s *Service) GetAllBaskets(id int) ([]entity.BasketOutput, error) {
	return s.repo.GetAllBaskets(id)
}

func (s *Service) GetTopUsers() ([]entity.TUser, error) {
	return s.repo.GetTopUsers()
}

func (s *Service) CreateStartingBasket(id int) (int, error) {
	return s.repo.CreateStartingBasket(id)
}

func (s *Service) UpdateBalance() (string, error) {
	return s.repo.UpdateBalance()
}

func (s *Service) UpdateBasket(name string) error {
	return s.repo.UpdateBasket(name)
}

func (s *Service) CreateBasketSell(id int, c string, v float64) (int, error) {
	return s.repo.CreateBasketSell(id, c, v)
}
