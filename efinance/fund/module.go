package fund

import (
	"context"
)

// Module 基金模块
type Module struct {
	GetQuoteHistory        func(ctx context.Context, fundCode string, pageSize int) ([]NetValueItem, error)
	GetQuoteHistoryMulti   func(ctx context.Context, fundCodes []string, pageSize int) (map[string][]NetValueItem, error)
	GetRealtimeEstimate   func(ctx context.Context, fundCodes []string) ([]EstNAVItem, error)
}
