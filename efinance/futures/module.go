package futures

import (
	"context"

	"github.com/efinance/efinance/efinance/common"
)

// Module 期货模块
type Module struct {
	GetBaseInfo func(ctx context.Context) ([]FuturesInfo, error)
	GetKline    func(ctx context.Context, quoteID string, beg, end string, klt common.KlineType, fqt common.AdjsType) ([]common.KlineItem, error)
}
