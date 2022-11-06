package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/emptywe/trading_sim/docs"
	"github.com/emptywe/trading_sim/internal/app"
	"github.com/emptywe/trading_sim/internal/router"
	"github.com/emptywe/trading_sim/internal/service"
	"github.com/emptywe/trading_sim/internal/storage"
)

func init() {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	db, err := storage.NewPostgresDB(storage.Config{
		Url: viper.GetString("db.url"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	currencyList := strings.Split(viper.GetString("currencies"), ",")
	repos := storage.NewRepository(db)
	services := service.NewService(repos)
	handlers := router.NewHandler(services, currencyList)
	applications := app.NewApplication(services, currencyList)

	srv := new(router.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			if err.Error() != "http: Server closed" {
				logrus.Fatalf("Can't run http server: %s", err)
			}
		}
	}()

	logrus.Println("Trading Sim Start")
	time.Sleep(time.Second)
	applications.InitApp()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Trading Sim Shout Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("can't shout down %s", err.Error())
	}

}
