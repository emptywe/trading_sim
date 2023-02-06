package simulator_repo

import (
	"github.com/jmoiron/sqlx"
	"net/http"

	"github.com/emptywe/trading_sim/model"
)

type Authorization interface {
	CreateUser(email, userName, password string) (int, error)
	GetUser(username, password string) (model.User, error)
	CreateSession(c *http.Cookie, id int, sid string) (int, string, error)
	ValidateSession(c *http.Cookie, id int, sid string) (int, bool)
	ExpireSession(c *http.Cookie, id int, sid string) *http.Cookie
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

type Repository struct {
	Authorization
	Basket
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Basket:        NewBasketPostgres(db),
	}
}
