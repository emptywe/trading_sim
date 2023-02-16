package main

import (
	"github.com/emptywe/trading_sim/internal/parser"
	"github.com/emptywe/trading_sim/internal/storage/postgres"
	"github.com/emptywe/trading_sim/internal/storage/postgres/parser_repo"
	"github.com/emptywe/trading_sim/pkg/wait"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

func execute() {
	db, err := postgres.NewDB(postgres.Config{
		Url: viper.GetString("postgres.url"),
	})
	if err != nil {
		zap.S().Fatal("failed to initialize db: " + err.Error())
	}
	defer db.Close()
	repo := parser_repo.NewRepository(db)
	zap.S().Infof("Parser start")
	parser.NewParser(repo, 5, strings.Split(viper.GetString("currencies"), ",")).InitParser()
	wait.Interrupt()
	zap.S().Infof("Parser shutdown")
}
