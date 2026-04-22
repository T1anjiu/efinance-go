package efinance

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/efinance/efinance/efinance/common"
)

// 调试URL生成和请求
func TestDebugKlineURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "短横线日期格式",
			url:  "https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayqfq&param=sh600519,day,2024-01-01,2024-01-31,100,qfq",
		},
		{
			name: "空日期格式",
			url:  "https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayqfq&param=sh600519,day,,,100,qfq",
		},
		{
			name: "只指定开始日期",
			url:  "https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayqfq&param=sh600519,day,2024-01-01,,100,qfq",
		},
		{
			name: "不带_var参数",
			url:  "https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?param=sh600519,day,2024-01-01,2024-01-31,100,qfq",
		},
		{
			name: "使用大写股票代码",
			url:  "https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayqfq&param=SH600519,day,2024-01-01,2024-01-31,100,qfq",
		},
	}

	client := common.DefaultClient()

	for _, test := range tests {
		fmt.Printf("\n=== [%s] ===\n", test.name)
		fmt.Printf("URL: %s\n", test.url)

		ctx, cancel := context.WithTimeout(context.Background(), 10*1000000000)
		req, err := http.NewRequestWithContext(ctx, "GET", test.url, nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			cancel()
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		req.Header.Set("Referer", "https://finance.qq.com/")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			cancel()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Printf("Error reading body: %v\n", err)
		} else {
			// 只显示前500字节
			content := string(body)
			if len(content) > 500 {
				content = content[:500] + "..."
			}
			fmt.Printf("Status: %d\nBody: %s\n", resp.StatusCode, content)
		}

		cancel()
	}
}

// 测试我的代码生成的URL
func TestMyCodeURL(t *testing.T) {
	client := common.DefaultClient()

	// 这应该是我的代码生成的URL
	url := "https://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayqfq&param=sh600519,day,2024-01-01,2024-01-31,100,qfq"
	
	fmt.Printf("Testing URL: %s\n", url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*1000000000)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://finance.qq.com/")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		cancel()
		return
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", string(body[:min(1000, len(body))]))
	
	cancel()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}