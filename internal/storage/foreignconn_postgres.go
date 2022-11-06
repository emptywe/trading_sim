package storage

import (
	"github.com/jmoiron/sqlx"
)

type ForeignConnPostgres struct {
	db *sqlx.DB
}

func NewForeignConnPostgres(db *sqlx.DB) *ForeignConnPostgres {
	return &ForeignConnPostgres{db: db}
}

func (r ForeignConnPostgres) UpdateCurrency(name string, value float64) error {
	_, err := r.db.Exec("UPDATE currencies set value=$1 WHERE name=$2", value, name)
	if err != nil {
		return err
	}
	return nil
}
