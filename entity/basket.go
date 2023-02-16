package entity

type Basket struct {
	Bid             int
	TransactionType string
	Currency        string
	Value           float64
	Amount          float64
	Status          string
}

type BasketOutput struct {
	Currency  string
	Value     float64
	USDAmount float64
}

type BalanceResponse struct {
	Balance string           `json:"balance"`
	Baskets []CurrencyOutput `json:"baskets"`
}
