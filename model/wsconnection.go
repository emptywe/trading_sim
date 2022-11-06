package model

type ConnRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     uint     `json:"id"`
}

type ConnResponse struct {
	EventType      string `json:"e"`
	EventTime      uint64 `json:"E"`
	Symbol         string `json:"s"`
	TradeId        uint64 `json:"t"`
	Price          string `json:"p"`
	Quantity       string `json:"q"`
	BuyerOrderId   uint64 `json:"b"`
	SellerOrderID  uint64 `json:"a"`
	TradeTime      uint64 `json:"T"`
	BuyerMarktMake bool   `json:"m"`
	Ignore         bool   `json:"M"`
}
