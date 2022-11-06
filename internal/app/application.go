package app

import (
	"github.com/emptywe/trading_sim/internal/service"
	"github.com/sirupsen/logrus"
)

type Application struct {
	services     *service.Service
	currencyList []string
}

func NewApplication(services *service.Service, currencyList []string) *Application {
	return &Application{services: services, currencyList: currencyList}
}

func (a *Application) InitApp() {
	a.CreateCurrencies()
	a.CurrencyUpdater()
}

//func (a *Application) BasketUpdater(){
//
//	for {
//		for _, cur := range CurrencyList {
//		 	err :=  a.services.Basket.UpdateBasket(cur)
//			if err != nil {
//				logrus.Printf("can't' update basket, error: %s", err.Error())
//			}
//		}
//		time.Sleep(time.Second*2)
//	}
//}

func (a *Application) CurrencyUpdater() {

	for _, cur := range a.currencyList {
		go a.services.ForeignConn.WShandlerBinance(cur)

	}

}

func (a *Application) CreateCurrencies() {
	for _, cur := range a.currencyList {
		err := a.services.Basket.CreateCurrencyTable(cur)
		if err != nil {
			logrus.Printf("can't create currency table, maybe it's already exist, error: %s", err.Error())
		}
	}
}

//func (a *Application) UserBalanceUpdater() {
//	for {
//		_, err := a.services.Basket.UpdateBalance()
//		if err != nil {
//			logrus.Errorf("can't update user balance error: %s", err.Error())
//		}
//		time.Sleep(time.Second*2)
//	}
//}

//func (a *Application) SessionGarbageCollector() {
//	for {
//		err := a.services.Authorization.SessionGarbageCollector()
//		if err != nil {
//			fmt.Println(err)
//		}
//		time.Sleep(time.Hour * 2)
//	}
//}
