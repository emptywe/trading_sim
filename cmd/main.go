package main

import (
	"go.uber.org/zap"

	_ "github.com/emptywe/trading_sim/docs"
	"github.com/emptywe/trading_sim/pkg/config"
	"github.com/emptywe/trading_sim/pkg/logger"
)

func init() {
	config.Init("config", "config")
}

func main() {
	logger.InitLogger(logger.DisableTrace)
	defer zap.S().Sync()
	zap.S().Debugw("Simulator logger initialised")
	execute()
}
