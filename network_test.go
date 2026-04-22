package efinance

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// 测试网络连接 - 寻找可用的K线API
func TestNetwork(t *testing.T) {
	// 测试不同的API端点
	urls := []string{
		// 腾讯实时行情 - 成功
		"https://qt.gtimg.cn/q=sh600519",
		// 腾讯K线日K (带日期范围)
		"https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayfqk&param=sh600519,day,2024-01-01,2024-01-31,100,qfq",
		// 腾讯K线周K
		"https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_weekfqk&param=sh600519,week,,,100,qfq",
	}

	for _, url := range urls {
		fmt.Printf("\n测试URL: %s\n", url)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		req.Header.Set("Referer", "https://finance.qq.com/")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("  ❌ 失败: %v\n", err)
			cancel()
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("  ✅ 状态码: %d, Body长度: %d bytes\n", resp.StatusCode, len(body))
		if len(body) < 2000 {
			fmt.Printf("  内容: %s\n", string(body))
		}

		cancel()
	}
}
