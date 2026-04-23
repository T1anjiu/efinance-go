package stock

import (
	"context"

	"github.com/T1anjiu/efinance-go/efinance/common"
)

// Module 股票模块
type Module struct {
	GetKline           func(ctx context.Context, params GetKlineParams) (*KlineResult, error)
	GetKlineMulti      func(ctx context.Context, params []GetKlineParams, workers int) (map[string]*KlineResult, error)
	GetRealtimeQuotes  func(ctx context.Context, params QuoteParams) ([]QuoteItem, error)
	GetLatestQuote     func(ctx context.Context, codes []string) ([]QuoteItem, error)
	Search             func(ctx context.Context, keyword string, marketType common.MarketType) ([]common.SearchResult, error)
	GetQuoteID         func(ctx context.Context, code string, marketType common.MarketType) (string, error)
}
