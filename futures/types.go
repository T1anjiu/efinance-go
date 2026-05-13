package futures

import "github.com/T1anjiu/efinance-go/common"

type FuturesBaseInfo struct {
	FuturesCode string `json:"futures_code"`
	FuturesName string `json:"futures_name"`
	QuoteID     string `json:"quote_id"`
	MarketType  string `json:"market_type"`
}

type FuturesQuote struct {
	common.RealtimeQuote
}

type FuturesKline struct {
	Code string           `json:"code"`
	Name string           `json:"name"`
	Bars []common.KlineBar `json:"bars"`
}
