package main

import (
	"github.com/emptywe/trading_sim/internal/storage/redis/simulator_cache/session_cache"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/emptywe/trading_sim/internal/simulator/api/v1/router"
	"github.com/emptywe/trading_sim/internal/simulator/api/v1/server"
	"github.com/emptywe/trading_sim/internal/simulator/service"
	"github.com/emptywe/trading_sim/internal/storage/postgres"
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
	myredis "github.com/emptywe/trading_sim/internal/storage/redis"
	"github.com/emptywe/trading_sim/pkg/logger"
	"github.com/emptywe/trading_sim/pkg/wait"
)

func initHandlers(pdb *sqlx.DB, rdb *redis.Client) *router.Handler {
	return router.NewHandler(
		service.NewService(simulator_repo.NewRepository(pdb), session_cache.NewCache(rdb)),
		strings.Split(viper.GetString("currencies"), ","),
	)
}

func execute() {
	srv := new(server.Server)

	pdb, err := postgres.NewDB(postgres.Config{
		Url: viper.GetString("postgres.url"),
	})
	if err != nil {
		zap.S().Fatalf("failed to initialize db: %v", err)
	}
	defer pdb.Close()
	rdb, err := myredis.NewDB(viper.GetString("redis.url"), viper.GetString("redis.password"), viper.GetInt("redis.db"))
	if err != nil {
		zap.S().Fatalf("failed to initialize cache: %v", err)
	}
	defer rdb.Close()

	go srv.StartNewServer(initHandlers(pdb, rdb).InitRoutes(), viper.GetInt("port"))
	srv.WaitServer()
	logger.InitLogger(logger.EnableTrace)
	wait.Interrupt()
	logger.InitLogger(logger.DisableTrace)
	srv.StopServer()
}
