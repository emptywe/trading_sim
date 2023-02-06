package parser

import (
	"github.com/emptywe/trading_sim/internal/storage/postgres/parser_repo"
	"github.com/emptywe/trading_sim/pkg/binance"
	"go.uber.org/zap"
)

type Parser struct {
	repo         *parser_repo.Repository
	currencyList []string
}

func NewParser(repo *parser_repo.Repository, currencyList []string) *Parser {
	return &Parser{repo: repo, currencyList: currencyList}
}

func (a *Parser) InitParser() {
	a.CreateCurrencies()
	a.CurrencyUpdater()
}

//func (a *Parser) BasketUpdater(){
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

func (a *Parser) CurrencyUpdater() {

	//for _, cur := range a.currencyList {
	//	go a.services.ForeignConn.WSHandlerBinance(cur)
	//
	//}

}

func (a *Parser) CreateCurrencies() {
	for _, cur := range a.currencyList {
		err := a.repo.CreateNewCurrency(cur)
		if err != nil {
			zap.S().Errorf("can't create currency table, maybe it's already exist, error: %v", err)
		}
	}
}

func (a *Parser) ParsePriceFromBinance() {
	client := binance.NewBinanceClient()
	client.WSHandlerBinance("BTCUSDT")
}

//func (a *Parser) UserBalanceUpdater() {
//	for {
//		_, err := a.services.Basket.UpdateBalance()
//		if err != nil {
//			logrus.Errorf("can't update user balance error: %s", err.Error())
//		}
//		time.Sleep(time.Second*2)
//	}
//}
