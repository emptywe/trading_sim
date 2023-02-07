package parser_repo

import (
	"github.com/jmoiron/sqlx"
)

type Parser interface {
	UpdateCurrency(name string, value float64) error
	CreateNewCurrency(name string) error
}

type Repository struct {
	Parser
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Parser: NewParserConnPostgres(db),
	}
}
