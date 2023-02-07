package main

import (
	"github.com/emptywe/trading_sim/internal/storage/postgres"
	"github.com/emptywe/trading_sim/internal/storage/postgres/parser_repo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

func initParserData(poolSize int) (*parser_repo.Repository, int, []string) {
	db, err := postgres.NewPostgresDB(postgres.Config{
		Url: viper.GetString("db.url"),
	})
	if err != nil {
		zap.S().Fatal("failed to initialize db: " + err.Error())
	}
	return parser_repo.NewRepository(db), 5, strings.Split(viper.GetString("currencies"), ",")
}
