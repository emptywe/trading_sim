package service

import (
	"github.com/emptywe/trading_sim/entity"
	"github.com/emptywe/trading_sim/internal/simulator/service/authentication"
	"github.com/emptywe/trading_sim/internal/simulator/service/basket"
	"github.com/emptywe/trading_sim/internal/simulator/service/info"
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
	ValidateSession(token string) (int, error)
	UpdateSession(token, rToken string) (string, error)
	DeleteSession(token string) error
}

type Basket interface {
	GetCurrency(name string) (entity.Currency, error)
	GetBasket(cur string, uid int) (entity.Basket, error)
	GetAllBaskets(id int) ([]entity.BasketOutput, error)
	CreateStartingBasket(uid int) error
	ServeTrade(ts entity.Transaction, uid int) error
}

type Info interface {
	GetTopUsers() ([]entity.TUser, error)
	GetAllCurrenciesUSD() ([]entity.CurrencyOutput, error)
}

type Service struct {
	Authorization
	Basket
	Info
}

func NewService(repos *simulator_repo.Repository, cache *session_cache.Cache) *Service {
	return &Service{
		Authorization: authentication.NewAuthService(repos.Authorization, cache),
		Basket:        basket.NewBasketService(repos.Basket),
		Info:          info.NewInfoService(repos.Info),
	}
}
