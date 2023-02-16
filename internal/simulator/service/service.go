package service

import (
	"net/http"

	"github.com/emptywe/trading_sim/entity"
	"github.com/emptywe/trading_sim/internal/simulator/service/authentication"
	"github.com/emptywe/trading_sim/internal/simulator/service/basket"
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
	"github.com/emptywe/trading_sim/internal/storage/redis/simulator_cache/session_cache"
	"github.com/emptywe/trading_sim/pkg/session"
)

type Authorization interface {
	CreateUser(user entity.SignUpRequest) (int, error)
	ReadUser(request entity.SignInRequest) (entity.User, error)
	UpdateUser(email, userName, password string) (int, error)
	DeleteUser(id int) error
	CreateSession(user *entity.User) (*session.Session, error)
	ValidateSession(token string) error
	UpdateSession(c *http.Cookie) (int, bool)
	DeleteSession(c *http.Cookie) *http.Cookie
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

type Service struct {
	Authorization
	Basket
}

func NewService(repos *simulator_repo.Repository, cache *session_cache.Cache) *Service {
	return &Service{
		Authorization: authentication.NewAuthService(repos.Authorization, cache),
		Basket:        basket.NewBasketService(repos.Basket),
	}
}
