package postgres

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	_ "github.com/emptywe/trading_sim/internal/storage/postgres/migrations"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
}

func NewDB(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	db, err := sqlx.Connect("pgx", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if err = migrateDB(db.DB, "schema_migration"); err != nil {
		return nil, err
	}

	return db, nil

}

func migrateDB(db *sql.DB, table string) error {
	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{MigrationsTable: table})
	migrator, err := migrate.NewWithDatabaseInstance("embed://", table, driver)
	if err != nil {
		return err
	}

	err = migrator.Up()
	if err != nil && err.Error() == "no change" { // "no change" is not an error
		err = nil
	}
	return err
}
