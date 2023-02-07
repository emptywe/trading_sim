package basket

import (
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
	"github.com/emptywe/trading_sim/model"
)

type BasketService struct {
	repo simulator_repo.Basket
}

func NewBasketService(repo simulator_repo.Basket) *BasketService {
	return &BasketService{repo: repo}
}

func (s *BasketService) GetCurrency(name string) (model.Currency, error) {
	return s.repo.GetCurrency(name)
}

func (s *BasketService) CreateBasket(id int, c1, c2 string, v float64) (int, error) {
	return s.repo.CreateBasket(id, c1, c2, v)
}

func (s *BasketService) GetBasket(id int, c string) (model.Basket, error) {
	return s.repo.GetBasket(id, c)
}

func (s *BasketService) GetAllCurrenciesUSD() ([]model.CurrencyOutput, error) {
	return s.repo.GetAllCurrenciesUSD()
}

func (s *BasketService) GetAllBaskets(id int) ([]model.BasketOutput, error) {
	return s.repo.GetAllBaskets(id)
}

func (s *BasketService) GetTopUsers() ([]model.TUser, error) {
	return s.repo.GetTopUsers()
}

func (s *BasketService) CreateStartingBasket(id int) (int, error) {
	return s.repo.CreateStartingBasket(id)
}

func (s *BasketService) UpdateBalance() (string, error) {
	return s.repo.UpdateBalance()
}

func (s *BasketService) UpdateBasket(name string) error {
	return s.repo.UpdateBasket(name)
}

func (s *BasketService) CreateBasketSell(id int, c string, v float64) (int, error) {
	return s.repo.CreateBasketSell(id, c, v)
}
