package main

import (
	"flag"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/emptywe/trading_sim/internal/parser"
	"github.com/emptywe/trading_sim/internal/simulator/api/v1/router"
	"github.com/emptywe/trading_sim/internal/simulator/api/v1/server"
	"github.com/emptywe/trading_sim/internal/simulator/service"
	"github.com/emptywe/trading_sim/internal/storage/postgres"
	"github.com/emptywe/trading_sim/internal/storage/postgres/parser_repo"
	"github.com/emptywe/trading_sim/internal/storage/postgres/simulator_repo"
	myredis "github.com/emptywe/trading_sim/internal/storage/redis"
	"github.com/emptywe/trading_sim/internal/storage/redis/simulator_cache/session_cache"
	"github.com/emptywe/trading_sim/pkg/logger"
	"github.com/emptywe/trading_sim/pkg/wait"
)

func initHandlers(pdb *sqlx.DB, rdb *redis.Client) *router.Handler {
	return router.NewHandler(
		service.NewService(simulator_repo.NewRepository(pdb), session_cache.NewCache(rdb)),
		strings.Split(viper.GetString("currencies"), ","),
	)
}

func setupConfigs(local bool) (pgConf postgres.Config, redisConf myredis.Config) {
	if local {
		pgConf = postgres.Config{
			Username: viper.GetString("postgres.username"),
			Password: viper.GetString("postgres.password"),
			Host:     viper.GetString("postgres.host"),
			Port:     viper.GetString("postgres.port"),
			DbName:   viper.GetString("postgres.dbName"),
		}
		redisConf = myredis.Config{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetString("redis.port"),
			Password: viper.GetString("redis.password"),
			Db:       viper.GetString("redis.db"),
		}
	} else {
		pgConf = postgres.Config{
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			DbName:   os.Getenv("DB_NAME"),
		}
		redisConf = myredis.Config{
			Host:     os.Getenv("CACHE_HOST"),
			Port:     os.Getenv("CACHE_PORT"),
			Password: os.Getenv("CACHE_PASSWORD"),
			Db:       os.Getenv("CACHE_DB"),
		}
	}
	return
}

func execute() {
	local := flag.Bool("local", false, "set local db and cache settings")
	flag.Parse()
	pgConf, redisConf := setupConfigs(*local)
	srv := new(server.Server)
	pdb, err := postgres.NewDB(pgConf)
	if err != nil {
		zap.S().Fatalf("failed to initialize postgres db: %v", err)
	}
	defer pdb.Close()
	rdb, err := myredis.NewDB(redisConf)
	if err != nil {
		zap.S().Fatalf("failed to initialize redis db: %v", err)
	}
	defer rdb.Close()

	go srv.StartNewServer(initHandlers(pdb, rdb).InitRoutes(), viper.GetInt("port"))
	srv.WaitServer()
	zap.S().Infof("Parser start")
	go parser.NewParser(parser_repo.NewRepository(pdb), 5, strings.Split(viper.GetString("currencies"), ",")).InitParser()
	logger.InitLogger(logger.EnableTrace)
	wait.Interrupt()
	logger.InitLogger(logger.DisableTrace)
	srv.StopServer()
}
