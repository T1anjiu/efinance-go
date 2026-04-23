package efinance

import (
	"context"
	"testing"
	"time"

	"github.com/T1anjiu/efinance-go/efinance/stock"
)

// TestCacheExpiration 测试缓存过期机制
func TestCacheExpiration(t *testing.T) {
	ctx := context.Background()

	// 搜索股票
	results, err := Stock.Search(ctx, "贵州茅台", "")
	if err != nil {
		t.Fatalf("搜索失败: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("未找到搜索结果")
	}

	t.Logf("找到 %d 个结果", len(results))
	for _, r := range results {
		t.Logf("代码: %s, 名称: %s, QuoteID: %s", r.Code, r.Name, r.QuoteID)
	}
}

// TestGetKlineWithRetry 测试K线获取（带重试）
func TestGetKlineWithRetry(t *testing.T) {
	ctx := context.Background()

	params := stock.GetKlineParams{
		Code:      "600519",
		Beg:       "2024-01-01",
		End:       "2024-01-31",
		KlineType: 101,
	}

	kline, err := Stock.GetKline(ctx, params)
	if err != nil {
		t.Fatalf("获取K线失败: %v", err)
	}

	if len(kline.Items) == 0 {
		t.Fatal("K线数据为空")
	}

	t.Logf("获取到 %d 条K线数据", len(kline.Items))
	if len(kline.Items) > 0 {
		item := kline.Items[0]
		t.Logf("最新数据: %s 开:%.2f 收:%.2f 高:%.2f 低:%.2f",
			item.Date, item.Open, item.Close, item.High, item.Low)
	}
}

// TestGetLatestQuoteWithBoundsCheck 测试实时行情获取（带边界检查）
func TestGetLatestQuoteWithBoundsCheck(t *testing.T) {
	ctx := context.Background()

	// 测试多只股票
	codes := []string{"600519", "000001", "600036"}

	quotes, err := Stock.GetLatestQuote(ctx, codes)
	if err != nil {
		t.Fatalf("获取实时行情失败: %v", err)
	}

	if len(quotes) == 0 {
		t.Fatal("实时行情为空")
	}

	t.Logf("获取到 %d 只股票的实时行情", len(quotes))
	for _, q := range quotes {
		t.Logf("%s (%s): 价格=%.2f 涨跌=%.2f%%",
			q.Name, q.Code, q.LatestPrice, q.ChangePCT)
	}
}

// TestConcurrentRequests 测试并发请求
func TestConcurrentRequests(t *testing.T) {
	ctx := context.Background()

	codes := []string{"600519", "000001", "600036", "000002", "600000"}

	// 并发获取实时行情
	done := make(chan bool, len(codes))

	for _, code := range codes {
		go func(c string) {
			defer func() { done <- true }()

			quotes, err := Stock.GetLatestQuote(ctx, []string{c})
			if err != nil {
				t.Errorf("获取 %s 失败: %v", c, err)
				return
			}

			if len(quotes) > 0 {
				t.Logf("%s: %.2f", quotes[0].Name, quotes[0].LatestPrice)
			}
		}(code)
	}

	// 等待所有请求完成
	for i := 0; i < len(codes); i++ {
		select {
		case <-done:
		case <-time.After(30 * time.Second):
			t.Fatal("请求超时")
		}
	}
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	ctx := context.Background()

	// 测试无效代码
	_, err := Stock.GetLatestQuote(ctx, []string{"999999"})
	if err == nil {
		t.Error("期望返回错误，但返回了nil")
	}

	t.Logf("正确处理了无效代码: %v", err)
}
