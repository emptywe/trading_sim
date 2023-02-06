package binance

const (
	wsUrl       = "wss://stream.parser.com:9443/ws"
	subscribe   = "SUBSCRIBE"
	unsubscribe = "UNSUBSCRIBE"
	trade       = "@trade"
	tradeUpd    = 500
)

type connRequest struct {
	method string   `json:"method"`
	params []string `json:"params"`
	id     uint     `json:"id"`
}

type DataPrice struct {
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
