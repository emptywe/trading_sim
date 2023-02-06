package simulator_repo

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/emptywe/trading_sim/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(email, userName, password string) (int, error) {
	var id int

	row := r.db.QueryRow("INSERT INTO users (email, username, password_hash, balance)  values ($1, $2, $3, $4) RETURNING id", email, userName, password, 0)
	if err := row.Scan(&id); err != nil {
		return 0, errors.New("username or email already used")
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (model.User, error) {
	var user model.User
	var dbp string

	row := r.db.QueryRow("SELECT password_hash FROM users WHERE username=$1", username)
	err := row.Scan(&dbp)
	if err != nil {
		return model.User{}, errors.New("invalid user")
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbp), []byte(password))
	if err != nil {
		return model.User{}, errors.New("wrong password")
	}

	row = r.db.QueryRow("SELECT id FROM users WHERE username=$1 AND password_hash=$2", username, dbp)
	err = row.Scan(&user.Uid)
	if err != nil {
		return model.User{}, errors.New("invalid user")
	}

	return user, err
}

func (r *AuthPostgres) CreateSession(c *http.Cookie, id int, sid string) (int, string, error) {
	err := r.SessionCleaner(id)
	if err != nil {
		zap.S().Errorf("sessions cleaned, error: %s", err.Error())
	}

	row := r.db.QueryRow("INSERT INTO session (user_id, sid, name, value, valid, established)  values ($1, $2, $3, $4, $5,$6) RETURNING user_id", id, sid, c.Name, c.Value, true, time.Now())
	if err := row.Scan(&id); err != nil {
		return 0, "", errors.New("can't create session> may be it's already exist")
	}
	return id, sid, nil
}

func (r *AuthPostgres) ValidateSession(c *http.Cookie, id int, sid string) (int, bool) {
	var valid bool

	row := r.db.QueryRow("SELECT valid FROM session WHERE sid=$1 AND user_id=$2", sid, id)
	if err := row.Scan(&valid); err != nil {
		return 0, false
	}
	_, err := r.db.Exec("UPDATE session set established=$1 where sid=$2", time.Now(), sid)
	if err != nil {
		zap.S().Errorf("can't update session, error: %s", err)
		return 0, false
	}
	return id, valid
}

func (r *AuthPostgres) ExpireSession(c *http.Cookie, id int, sid string) *http.Cookie {

	r.db.QueryRow("DELETE from session WHERE sid=$1 AND user_id=$2", sid, id)
	return c
}

func (r *AuthPostgres) DropUser(id int) error {

	_ = r.db.QueryRow("DELETE from users WHERE id=$1", id)
	return nil
}

func (r *AuthPostgres) SessionGarbageCollector() error {
	var cSess []model.CleanSession

	rows, err := r.db.Query("SELECT sid, established from session ")
	if err != nil {
		return err
		//return errors.New("basket with this value not exist yet")
	}

	for rows.Next() {
		var cSes model.CleanSession
		if err := rows.Scan(&cSes.Sid, &cSes.Established); err != nil {
			return errors.New("internal scan sessions problem")
		}

		cSess = append(cSess, cSes)
	}
	if err := rows.Err(); err != nil {
		return errors.New("corrupted data session")
	}

	for i, _ := range cSess {
		if time.Now().Sub(cSess[i].Established) > time.Hour*2 {
			_, err = r.db.Exec("DELETE from session where sid=$1", cSess[i].Sid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *AuthPostgres) SessionCleaner(id int) error {
	var cSess []model.CleanSession

	rows, err := r.db.Query("SELECT sid from session where user_id=$1", id)
	if err != nil {
		return err
		//return errors.New("basket with this value not exist yet")
	}

	for rows.Next() {
		var cSes model.CleanSession
		if err := rows.Scan(&cSes.Sid); err != nil {
			return errors.New("internal scan sessions problem")
		}

		cSess = append(cSess, cSes)
	}
	if err := rows.Err(); err != nil {
		return errors.New("corrupted data session")
	}

	for i, _ := range cSess {
		_, err = r.db.Exec("DELETE from session where sid=$1", cSess[i].Sid)
		if err != nil {
			return err
		}
	}

	return nil
}
