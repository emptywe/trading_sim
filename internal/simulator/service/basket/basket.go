package basket

import (
	"errors"
	"fmt"

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

func (s *Service) GetBasket(cur string, uid int) (entity.Basket, error) {
	return s.repo.GetBasket(cur, uid)
}

func (s *Service) GetAllBaskets(id int) ([]entity.BasketOutput, error) {
	return s.repo.GetAllBaskets(id)
}

func (s *Service) CreateStartingBasket(uid int) error {
	cur, err := s.repo.GetCurrency(entity.BaseCurrency)
	if err != nil {
		return err
	}
	_, err = s.repo.CreateBasket(uid, cur, entity.StartingBalance)
	return err
}

func (s *Service) ServeTrade(ts entity.Transaction, uid int) error {
	var bid2 int
	b1, err := s.repo.GetBasket(ts.BaseCurrency, uid)
	if err != nil {
		return fmt.Errorf("you don't have such currency yet: %v", err)
	}
	cur1, err := s.repo.GetCurrency(ts.BaseCurrency)
	if err != nil {
		return err
	}
	cur2, err := s.repo.GetCurrency(ts.TradeCurrency)
	if err != nil {
		return err
	}

	if b1.ValueUSD < ts.TradeAmount*cur2.Value {
		return errors.New("not enough base currency amount")
	}

	b2, err := s.repo.GetBasket(ts.TradeCurrency, uid)
	if err != nil {
		bid2, err = s.repo.CreateBasket(uid, cur2, entity.DefaultBalance)
		if err != nil {
			return err
		}
	} else {
		bid2 = b2.Bid
	}
	err = s.repo.UpdateBasketAmount(-ts.TradeAmount*cur2.Value/cur1.Value, b1.Bid)
	if err != nil {
		return err
	}
	err = s.repo.UpdateBasketAmount(ts.TradeAmount, bid2)
	if err != nil {
		return err
	}

	return nil
}
