package bond

import (
	"context"
)

// Module 债券模块
type Module struct {
	GetBaseInfo     func(ctx context.Context, bondCode string) (*BondInfo, error)
	GetBaseInfoMulti func(ctx context.Context, bondCodes []string) ([]BondInfo, error)
	GetAllBaseInfo   func(ctx context.Context) ([]BondInfo, error)
}
