package simulator_repo

import (
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/emptywe/trading_sim/entity"
)

type InfoPostgres struct {
	db *sqlx.DB
}

func NewInfoPostgres(db *sqlx.DB) *InfoPostgres {
	return &InfoPostgres{db: db}
}

func (r InfoPostgres) GetAllCurrenciesUSD() ([]entity.CurrencyOutput, error) {
	var data []entity.CurrencyOutput

	rows, err := r.db.Query("SELECT name, value from currencies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dat entity.CurrencyOutput
		if err := rows.Scan(&dat.Name, &dat.Value); err != nil {
			return data, err
		}
		data = append(data, dat)
	}
	if err = rows.Err(); err != nil {
		return data, errors.New("corrupted data")
	}

	return data, nil
}

func (r InfoPostgres) GetTopUsers() ([]entity.TUser, error) {
	var data []entity.TUser

	rows, err := r.db.Query("SELECT username, balance from users ORDER BY balance ASC LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dat entity.TUser
		if err := rows.Scan(&dat.UserName, &dat.Balance); err != nil {
			return data, errors.New("can't find users")
		}

		data = append(data, dat)
	}
	if err = rows.Err(); err != nil {
		return data, err
	}

	return data, nil
}
