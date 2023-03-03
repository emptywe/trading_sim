package authentication

import (
	"context"

	"github.com/emptywe/trading_sim/entity"
	"github.com/emptywe/trading_sim/pkg/session"
)

func (s *Service) CreateSession(user *entity.User) (*session.Session, error) {
	ses, sUUID, err := session.NewSession(user.UserName, user.Uid)
	if err != nil {
		return nil, err
	}
	if err = s.cache.Create(context.Background(), user.UserName, sUUID, ses.Token); err != nil {
		return nil, err
	}
	return ses, err
}

func (s *Service) ValidateSession(token string) (int, error) {
	if err := session.ValidateToken(token); err != nil {
		return 0, err
	}
	claims, err := session.ParseToken(&session.TokenClaims{}, token)
	if err != nil {
		return 0, err
	}
	if err = s.cache.Read(context.Background(), claims.(*session.TokenClaims).UserName, claims.(*session.TokenClaims).Sid, token); err != nil {
		return 0, err
	}
	return claims.(*session.TokenClaims).UserId, nil
}

func (s *Service) UpdateSession(token, rToken string) (string, error) {
	claims, err := session.ParseToken(&session.TokenClaims{}, token)
	if err != nil {
		return "", err
	}
	newToken, err := session.UpdateToken(claims.(*session.TokenClaims), rToken)
	if err != nil {
		return "", err
	}
	err = s.cache.Update(context.Background(), claims.(*session.TokenClaims).UserName, claims.(*session.TokenClaims).Sid, newToken)
	return newToken, err
}

func (s *Service) DeleteSession(token string) error {
	claims, err := session.ParseToken(&session.TokenClaims{}, token)
	if err != nil {
		return err
	}
	err = s.cache.Delete(context.Background(), claims.(*session.TokenClaims).UserName)
	if err != nil {
		return err
	}
	return nil
}
