package efinance

import (
	"context"
	"fmt"
	"testing"

	"github.com/T1anjiu/efinance-go/efinance/stock"
)

// 测试获取K线数据
func TestGetKline(t *testing.T) {
	ctx := context.Background()

	params := stock.GetKlineParams{
		Code:      "600519",
		Beg:       "2024-01-01",
		End:       "2024-01-31",
		KlineType: 101,
	}

	kline, err := Stock.GetKline(ctx, params)
	if err != nil {
		t.Fatalf("GetKline failed: %v", err)
	}

	fmt.Printf("Kline: %d records\n", len(kline.Items))
	if len(kline.Items) > 0 {
		k := kline.Items[0]
		fmt.Printf("Latest: %s O:%.2f C:%.2f H:%.2f L:%.2f\n",
			k.Date, k.Open, k.Close, k.High, k.Low)
	}
}

// 测试获取实时行情
func TestGetRealtimeQuotes(t *testing.T) {
	ctx := context.Background()

	params := stock.QuoteParams{}
	quotes, err := Stock.GetRealtimeQuotes(ctx, params)
	if err != nil {
		t.Fatalf("GetRealtimeQuotes failed: %v", err)
	}

	fmt.Printf("Quotes: %d stocks\n", len(quotes))
	for i, q := range quotes {
		if i >= 3 {
			break
		}
		fmt.Printf("%s (%s): Price=%.2f Change=%.2f%%\n",
			q.Name, q.Code, q.LatestPrice, q.ChangePCT)
	}
}

// 测试获取指定股票实时行情
func TestGetLatestQuote(t *testing.T) {
	ctx := context.Background()

	quotes, err := Stock.GetLatestQuote(ctx, []string{"600519", "000001", "600036"})
	if err != nil {
		t.Fatalf("GetLatestQuote failed: %v", err)
	}

	fmt.Printf("Latest Quotes: %d stocks\n", len(quotes))
	for _, q := range quotes {
		fmt.Printf("%s (%s): Price=%.2f Change=%.2f%%\n",
			q.Name, q.Code, q.LatestPrice, q.ChangePCT)
	}
}
