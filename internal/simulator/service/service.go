package service

import (
	"github.com/emptywe/trading_sim/internal/simulator/service/authentication"
	"github.com/emptywe/trading_sim/internal/simulator/service/basket"
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
	"net/http"

	"github.com/emptywe/trading_sim/model"
)

type Authorization interface {
	CreateUser(email, userName, password string) (int, error)
	GenerateToken(username, password string) (string, error)
	CreateSession(c *http.Cookie) (int, string, error)
	ValidateSession(c *http.Cookie) (int, bool)
	ExpireSession(c *http.Cookie) *http.Cookie
	DropUser(id int) error
	SessionGarbageCollector() error
}

type Basket interface {
	GetCurrency(name string) (model.Currency, error)
	CreateBasket(id int, c1, c2 string, v float64) (int, error)
	GetBasket(id int, c string) (model.Basket, error)
	GetAllCurrenciesUSD() ([]model.CurrencyOutput, error)
	GetAllBaskets(id int) ([]model.BasketOutput, error)
	GetTopUsers() ([]model.TUser, error)
	CreateStartingBasket(id int) (int, error)
	UpdateBalance() (string, error)
	UpdateBasket(name string) error
	CreateBasketSell(id int, c string, v float64) (int, error)
}

type Service struct {
	Authorization
	Basket
}

func NewService(repos *simulator_repo.Repository) *Service {
	return &Service{
		Authorization: authentication.NewAuthService(repos.Authorization),
		Basket:        basket.NewBasketService(repos.Basket),
	}
}
