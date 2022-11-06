package storage

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	_ "github.com/emptywe/trading_sim/internal/storage/migrations"
)

type Config struct {
	Url string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", cfg.Url)
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
