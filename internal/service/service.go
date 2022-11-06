package service

import (
	"github.com/emptywe/trading_sim/internal/storage"
	"net/http"

	"github.com/emptywe/trading_sim/model"
)

type Authorization interface {
	CreateUser(email, userName, password string) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accesToken string) (int, string, error)
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
	CreateCurrencyTable(name string) error
	CreateBasketSell(id int, c string, v float64) (int, error)
}

type ForeignConn interface {
	WShandlerBinance(cur string)
}

type Service struct {
	Authorization
	Basket
	ForeignConn
}

func NewService(repos *storage.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Basket:        NewBasketService(repos.Basket),
		ForeignConn:   NewForeignConnService(repos.ForeignConn),
	}
}
