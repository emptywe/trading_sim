package simulator_repo

import (
	"github.com/jmoiron/sqlx"

	"github.com/emptywe/trading_sim/entity"
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
	GetBasket(cur string, uid int) (entity.Basket, error)
	GetAllBaskets(id int) ([]entity.BasketOutput, error)
	CreateBasket(uid int, cur entity.Currency, amount float64) (int, error)
	UpdateBasketAmount(delta float64, bid int) error
}

type Info interface {
	GetAllCurrenciesUSD() ([]entity.CurrencyOutput, error)
	GetTopUsers() ([]entity.TUser, error)
}

type Repository struct {
	Authorization
	Basket
	Info
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Basket:        NewBasketPostgres(db),
		Info:          NewInfoPostgres(db),
	}
}
