package simulator_repo

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/emptywe/trading_sim/entity"
)

type BasketPostgres struct {
	db *sqlx.DB
}

func NewBasketPostgres(db *sqlx.DB) *BasketPostgres {
	return &BasketPostgres{db: db}
}

func (r BasketPostgres) GetCurrency(name string) (entity.Currency, error) {
	var currency entity.Currency

	row := r.db.QueryRow("SELECT * FROM currencies WHERE name=$1", name)
	if err := row.Scan(&currency.Cid, &currency.Name, &currency.Value); err != nil {
		return entity.Currency{}, errors.New("no such currency")
	}
	return currency, nil
}

func (r BasketPostgres) GetBasket(id int, c string) (entity.Basket, error) {

	var basket entity.Basket

	row := r.db.QueryRow("SELECT id, transactiontype, value, amount ,status FROM basket WHERE user_id=$1 AND currency=$2", id, c)
	if err := row.Scan(&basket.Bid, &basket.TransactionType, &basket.Value, &basket.Amount, &basket.Status); err != nil {
		return entity.Basket{}, errors.New("asset doesn't exist")
	}
	basket.Currency = c
	return basket, nil
}

func (r BasketPostgres) CreateBasket(id int, c1, c2 string, v float64) (int, error) {

	var bId int

	cu1, err := r.GetCurrency(c1)
	if err != nil {
		return 0, err
	}
	cu2, err := r.GetCurrency(c2)
	if err != nil {
		return 0, err
	}
	b1, err := r.GetBasket(id, cu1.Name)
	if err != nil {
		return 0, errors.New("you don't have such currency")
	}
	if b1.Value*cu1.Value < v*cu2.Value {
		return 0, errors.New("you don't have enough currency")
	}

	res := b1.Value*cu1.Value - v*cu2.Value

	_, err = r.db.Exec("UPDATE basket set value=$1, amount=$2 where user_id=$3 AND currency=$4", res/cu1.Value, res, id, c1)
	if err != nil {
		return 0, err
	}

	b2, err := r.GetBasket(id, cu2.Name)
	if err == nil {

		row := r.db.QueryRow("UPDATE basket set value=$1, amount=$2 where user_id=$3 AND currency=$4 RETURNING id", b2.Value+v, (b2.Value+v)*cu2.Value, id, c2)
		if err := row.Scan(&bId); err != nil {
			_, _ = r.db.Exec("UPDATE basket set value=$1, amount=$2 where user_id=$3 AND currency=$4", b1.Value, b1.Amount, id, c1)
			return 0, errors.New("corrupted asset data")
		}

		return bId, nil
	}

	row := r.db.QueryRow("INSERT INTO basket (user_id, transactiontype, currency_id, currency, value, amount, status)  values ($1, $2, $3, $4, $5,$6, $7) RETURNING id", id, "buy", cu2.Cid, cu2.Name, v, v*cu2.Value, "exist")
	if err := row.Scan(&bId); err != nil {
		_, _ = r.db.Exec("UPDATE basket set value=$1, amount=$2 where user_id=$3 AND currency=$4", b1.Value, b1.Amount, id, c1)
		return 0, errors.New("corrupted new asset data")
	}

	return bId, nil

}

func (r BasketPostgres) GetAllCurrenciesUSD() ([]entity.CurrencyOutput, error) {
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

func (r BasketPostgres) GetAllBaskets(id int) ([]entity.BasketOutput, error) {
	var data []entity.BasketOutput

	rows, err := r.db.Query("SELECT  currency,value, amount FROM basket WHERE user_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dat entity.BasketOutput
		if err := rows.Scan(&dat.Currency, &dat.Value, &dat.USDAmount); err != nil {
			return data, err
		}
		data = append(data, dat)
	}
	if err = rows.Err(); err != nil {
		return data, errors.New("corrupted data")
	}

	return data, nil
}

func (r BasketPostgres) GetTopUsers() ([]entity.TUser, error) {
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

func (r BasketPostgres) CreateStartingBasket(id int) (int, error) {
	var bId int

	cu, err := r.GetCurrency("usdt")
	if err != nil {
		return 0, err
	}

	row := r.db.QueryRow("INSERT INTO basket (user_id, currency_id, currency, amount)  values ($1, $2, $3, $4) RETURNING id", id, cu.Cid, cu.Name, 1000)
	err = row.Scan(&bId)
	if err != nil {
		return 0, fmt.Errorf("db: %v", err)
	}
	return bId, err
}

func (r BasketPostgres) GetAllUsers() ([]entity.User, error) {
	var data []entity.User

	rows, err := r.db.Query("SELECT id, username, balance from users")
	if err != nil {
		return nil, errors.New("internal query problem")
	}
	defer rows.Close()

	for rows.Next() {
		var dat entity.User
		if err := rows.Scan(&dat.Uid, &dat.UserName, &dat.Balance); err != nil {
			return data, errors.New("internal scan users problem")
		}

		data = append(data, dat)
	}
	if err = rows.Err(); err != nil {
		return data, errors.New("corrupted data users")
	}

	return data, nil
}

func (r BasketPostgres) UpdateBalance() (string, error) {
	uu, err := r.GetAllUsers()
	if err != nil {
		return "error", errors.New("can't get users from db")
	}
	for i, _ := range uu {
		bb, err := r.GetAllBaskets(uu[i].Uid)
		if err != nil {
			return "error", errors.New("can't get data about users from db")
		}

		var b float64
		for i, _ := range bb {
			b += bb[i].USDAmount
		}
		_, err = r.db.Exec("UPDATE users set balance=$1 where id=$2", b, uu[i].Uid)
		if err != nil {
			zap.S().Errorf("problem update balance %s", err.Error())
		}

	}

	return "Updated", nil

}

func (r BasketPostgres) UpdateBasket(name string) error {
	var basket []entity.Basket

	rows, err := r.db.Query("SELECT id, value from basket where currency=$1", name)
	if err != nil {
		return err
		//return errors.New("basket with this value not exist yet")
	}

	for rows.Next() {
		var bask entity.Basket
		if err := rows.Scan(&bask.Bid, &bask.Value); err != nil {
			return errors.New("internal scan basket problem")
		}

		basket = append(basket, bask)
	}
	if err := rows.Err(); err != nil {
		return errors.New("corrupted data basket")
	}

	c, err := r.GetCurrency(name)
	if err != nil {
		return err
	}

	for i, _ := range basket {
		newamount := basket[i].Value * c.Value
		_, err := r.db.Exec("UPDATE basket set amount=$1 where id=$2", newamount, basket[i].Bid)
		if err != nil {
			zap.S().Errorf("can't update amount on basket %v", basket[i].Bid)
		}
	}

	return nil
}

func (r BasketPostgres) CreateBasketSell(id int, c string, v float64) (int, error) {
	var bId int

	cu1, err := r.GetCurrency(c)
	if err != nil {
		return 0, err
	}
	cu2, err := r.GetCurrency("usdt")
	if err != nil {
		return 0, err
	}
	b1, err := r.GetBasket(id, cu1.Name)
	if err != nil {
		return 0, errors.New("you don't have such currency")
	}
	if b1.Value*cu1.Value < v*cu2.Value {
		return 0, errors.New("you don't have enough currency")
	}

	res := b1.Value*cu1.Value - v*cu2.Value

	//queryt:= fmt.Sprintf("UPDATE %s set value=$1, amount=$2 where user_id=$3 AND currency=$4", basketTable)
	_, err = r.db.Exec("UPDATE basket set value=$1, amount=$2 where user_id=$3 AND currency=$4", res/cu1.Value, res, id, c)
	if err != nil {
		return 0, err
	}

	b2, err := r.GetBasket(id, cu2.Name)
	if err == nil {
		//query:= fmt.Sprintf("UPDATE %s set value=$1, amount=$2 where user_id=$3 AND currency=$4 RETURNING id", basketTable)
		row := r.db.QueryRow("UPDATE basket set value=$1, amount=$2 where user_id=$3 AND currency=$4 RETURNING id", b2.Value+v, (b2.Value+v)*cu2.Value, id, "usdt")
		if err := row.Scan(&bId); err != nil {
			_, _ = r.db.Exec("UPDATE basket set value=$1, amount=$2 where user_id=$3 AND currency=$4", b1.Value, b1.Amount, id, c)
			return 0, errors.New("corrupted asset data")
		}

		return bId, nil
	}

	//query:= fmt.Sprintf("INSERT INTO %s (user_id, transactiontype, currency_id, currency, value, amount, status)  values ($1, $2, $3, $4, $5,$6, $7) RETURNING id", basketTable)
	row := r.db.QueryRow("INSERT INTO basket (user_id, transactiontype, currency_id, currency, value, amount, status)  values ($1, $2, $3, $4, $5,$6, $7) RETURNING id", id, "buy", cu2.Cid, cu2.Name, v, v*cu2.Value, "exist")
	if err := row.Scan(&bId); err != nil {
		_, _ = r.db.Exec("UPDATE basket set value=$1, amount=$2 where user_id=$3 AND currency=$4", b1.Value, b1.Amount, id, c)
		return 0, errors.New("corrupted new asset data")
	}

	return bId, nil
}
