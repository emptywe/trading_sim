package main

import (
	"github.com/emptywe/trading_sim/pkg/config"
	"github.com/emptywe/trading_sim/pkg/logger"
	"go.uber.org/zap"
)

func init() {
	config.InitConfig("config", "config")
}

func main() {
	logger.InitLogger(logger.DisableTrace)
	defer zap.S().Sync()
	zap.S().Debugw("Parser logger initialised")
	execute()
}
