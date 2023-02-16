package simulator_repo

import (
	"errors"
	"fmt"
	"github.com/emptywe/trading_sim/entity"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(request entity.SignUpRequest) (int, error) {
	var id int
	ts := time.Now()
	row := r.db.QueryRow("INSERT INTO users (email,firstname,lastname, username, password_hash,status, created_at, updated_at, last_signed_at)  values ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", request.Email, "", "", request.UserName, request.Password, "partner", ts, ts, ts)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("db:%v", err)
	}
	return id, nil
}

func (r *AuthPostgres) ValidatePassword(request entity.SignInRequest) error {
	var pass string
	row := r.db.QueryRow("SELECT password_hash FROM users WHERE username = $1", request.UserName)
	if err := row.Scan(&pass); err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(request.Password)); err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) ReadUser(username string) (entity.User, error) {
	var user entity.User
	row := r.db.QueryRow(`SELECT users.id, users.firstname, users.lastname, users.username, users.email, users.status, SUM(currencies.value * basket.amount) AS balance FROM users INNER JOIN basket on users.id = basket.user_id  INNER JOIN currencies on basket.currency_id = currencies.id WHERE username=$1 GROUP BY users.id, users.firstname, users.lastname, users.username, users.email, users.status`,
		username)
	if err := row.Scan(&user.Uid, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Status, &user.Balance); err != nil {
		return entity.User{}, fmt.Errorf("invalid user: %v", err)
	}
	return user, nil
}

func (r *AuthPostgres) UpdateUser(email, userName, password string) (int, error) {
	var id int

	row := r.db.QueryRow("INSERT INTO users (email, username, password_hash)  values ($1, $2, $3, $4) RETURNING id", email, userName, password)
	if err := row.Scan(&id); err != nil {
		return 0, errors.New("username or email already used")
	}
	return id, nil
}

func (r *AuthPostgres) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE from users WHERE id=$1", id)
	return err
}
