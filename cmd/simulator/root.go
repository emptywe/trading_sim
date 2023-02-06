package main

import (
	"github.com/emptywe/trading_sim/internal/simulator/api/v1/router"
	"github.com/emptywe/trading_sim/internal/simulator/api/v1/server"
	"github.com/emptywe/trading_sim/internal/simulator/service"
	"github.com/emptywe/trading_sim/internal/storage/postgres"
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
	"github.com/emptywe/trading_sim/pkg/logger"
	"github.com/emptywe/trading_sim/pkg/wait"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

func initHandlers() *router.Handler {
	db, err := postgres.NewPostgresDB(postgres.Config{
		Url: viper.GetString("db.url"),
	})
	if err != nil {
		zap.S().Fatal("failed to initialize db: " + err.Error())
	}
	return router.NewHandler(service.NewService(simulator_repo.NewRepository(db)), strings.Split(viper.GetString("currencies"), ","))
}

func Execute() {
	srv := new(server.Server)
	go srv.StartNewServer(initHandlers().InitRoutes(), viper.GetInt("port"))
	srv.WaitServer()
	logger.InitLogger(logger.EnableTrace)
	wait.WaitInterrupt()
	logger.InitLogger(logger.DisableTrace)
	srv.StopServer()
}
