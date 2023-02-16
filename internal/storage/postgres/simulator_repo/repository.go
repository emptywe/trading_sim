package simulator_repo

import (
	"github.com/emptywe/trading_sim/entity"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user entity.SignUpRequest) (int, error)
	ReadUser(username string) (entity.User, error)
	UpdateUser(email, userName, password string) (int, error)
	DeleteUser(id int) error
	ValidatePassword(request entity.SignInRequest) error
}

type Basket interface {
	GetCurrency(name string) (entity.Currency, error)
	CreateBasket(id int, c1, c2 string, v float64) (int, error)
	GetBasket(id int, c string) (entity.Basket, error)
	GetAllCurrenciesUSD() ([]entity.CurrencyOutput, error)
	GetAllBaskets(id int) ([]entity.BasketOutput, error)
	GetTopUsers() ([]entity.TUser, error)
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
