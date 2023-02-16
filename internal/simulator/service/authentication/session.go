package authentication

import (
	"context"
	"net/http"

	"github.com/emptywe/trading_sim/entity"
	"github.com/emptywe/trading_sim/pkg/session"
)

func (s *Service) CreateSession(user *entity.User) (*session.Session, error) {
	ses, sUUID, err := session.NewSession(user.UserName, user.Uid)
	if err != nil {
		return nil, err
	}
	if err = s.cache.Create(context.TODO(), user.UserName, sUUID, ses.Token); err != nil {
		return nil, err
	}
	return ses, err
}

func (s *Service) ValidateSession(token string) error {
	claims, err := session.ParseToken(&session.TokenClaims{}, token)
	if err != nil {
		return err
	}
	if err = s.cache.Read(context.Background(), claims.(*session.TokenClaims).UserName, claims.(*session.TokenClaims).Sid, token); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateSession(c *http.Cookie) (int, bool) {

	//s.repo.ValidateSession(c, id, sid)
	return 1, false
}

func (s *Service) DeleteSession(c *http.Cookie) *http.Cookie {

	//s.repo.ExpireSession(c, id, sid)
	return nil
}
