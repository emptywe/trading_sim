package entity

const (
	BaseCurrency    = "usd"
	StartingBalance = 1000
	DefaultBalance  = 0
)

type Basket struct {
	Bid      int
	Currency string
	Amount   float64
	ValueUSD float64
}

type BasketOutput struct {
	Currency  string
	Amount    float64
	Price     float64
	USDAmount float64
}

type BalanceResponse struct {
	Balance string           `json:"balance"`
	Baskets []CurrencyOutput `json:"baskets"`
}

type Transaction struct {
	BaseCurrency  string  `json:"base"`
	TradeCurrency string  `json:"trade"`
	TradeAmount   float64 `json:"amount"`
}
