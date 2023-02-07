package main

import (
	"github.com/emptywe/trading_sim/internal/parser"
	"github.com/emptywe/trading_sim/pkg/config"
	"github.com/emptywe/trading_sim/pkg/logger"
	"github.com/emptywe/trading_sim/pkg/wait"
	"go.uber.org/zap"
)

func init() {
	config.InitConfig("config", "config")
}

func main() {
	logger.InitLogger(logger.DisableTrace)
	defer zap.S().Sync()
	zap.S().Debugw("Parser logger initialised")
	zap.S().Infof("Parser start")
	parser.NewParser(initParserData(5)).InitParser()
	wait.WaitInterrupt()
	zap.S().Infof("Parser shutdown")
}
