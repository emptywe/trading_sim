package simulator_repo

import (
	"errors"
	"fmt"

	"github.com/emptywe/trading_sim/entity"
	"github.com/jmoiron/sqlx"
)

type BasketPostgres struct {
	db *sqlx.DB
}

func NewBasketPostgres(db *sqlx.DB) *BasketPostgres {
	return &BasketPostgres{db: db}
}

func (r BasketPostgres) CreateBasket(uid int, cur entity.Currency, amount float64) (id int, err error) {
	row := r.db.QueryRow("INSERT INTO basket (user_id,currency_id,currency,amount) VALUES ($1,$2,$3,$4) RETURNING id", uid, cur.Cid, cur.Name, amount)
	if err = row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r BasketPostgres) GetCurrency(name string) (entity.Currency, error) {
	var currency entity.Currency

	row := r.db.QueryRow("SELECT * FROM currencies WHERE name=$1", name)
	if err := row.Scan(&currency.Cid, &currency.Name, &currency.Value); err != nil {
		return entity.Currency{}, errors.New("no such currency")
	}
	return currency, nil
}

func (r BasketPostgres) GetBasket(cur string, uid int) (basket entity.Basket, err error) {
	row := r.db.QueryRow("SELECT basket.id,basket.currency,basket.amount,basket.amount*currencies.value FROM basket INNER JOIN currencies ON basket.currency_id = currencies.id WHERE user_id=$1 AND currency=$2", uid, cur)
	if err = row.Scan(&basket.Bid, &basket.Currency, &basket.Amount, &basket.ValueUSD); err != nil {
		return entity.Basket{}, fmt.Errorf("asset doesn't exist: %v", err)
	}
	return basket, nil
}

func (r BasketPostgres) UpdateBasketAmount(delta float64, bid int) error {
	_, err := r.db.Exec("UPDATE basket set amount = amount+$1 WHERE id = $2", delta, bid)
	return err
}

func (r BasketPostgres) GetAllBaskets(id int) ([]entity.BasketOutput, error) {
	var data []entity.BasketOutput

	rows, err := r.db.Query("SELECT  currency,amount,currencies.value ,basket.amount*currencies.value FROM basket INNER JOIN currencies ON basket.currency_id = currencies.id WHERE user_id=$1 AND amount > 0", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dat entity.BasketOutput
		if err := rows.Scan(&dat.Currency, &dat.Amount, &dat.Price, &dat.USDAmount); err != nil {
			return data, err
		}
		data = append(data, dat)
	}
	if err = rows.Err(); err != nil {
		return data, errors.New("corrupted data")
	}

	return data, nil
}
