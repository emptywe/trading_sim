package parser_repo

import (
	"github.com/jmoiron/sqlx"
	"strings"
)

type ParserConnPostgres struct {
	db *sqlx.DB
}

func NewParserConnPostgres(db *sqlx.DB) *ParserConnPostgres {
	return &ParserConnPostgres{db: db}
}

func (r ParserConnPostgres) UpdateCurrency(name string, value float64) (err error) {
	_, err = r.db.Exec("UPDATE currencies set value=$1 WHERE name=$2", value, name)
	if err != nil {
		return err
	}
	return
}

func (r ParserConnPostgres) CreateNewCurrency(name string) error {
	_, err := r.db.Exec("INSERT into currencies (name,value) values ($1, $2)", name, 1)
	if err != nil {
		if !strings.Contains(err.Error(), "duplicate key value violates unique constraint ") {
			return err
		}
	}
	return nil
}
