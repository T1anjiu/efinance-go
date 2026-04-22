package efinance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client HTTP客户端
type Client struct {
	httpClient *http.Client
	headers    map[string]string
}

// NewClient 创建客户端
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: map[string]string{
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":          "https://finance.eastmoney.com",
			"Accept":           "application/json, text/plain, */*",
			"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8",
		},
	}
}

// Get 发送GET请求
func (c *Client) Get(ctx context.Context, url string, params map[string]interface{}) ([]byte, error) {
	// 构建URL
	if len(params) > 0 {
		url += "?"
		for k, v := range params {
			url += fmt.Sprintf("%s=%v&", k, v)
		}
		url = url[:len(url)-1] // 去掉最后的&
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Post 发送POST请求
func (c *Client) Post(ctx context.Context, url string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// 默认客户端
var defaultClient = NewClient()

// SetTimeout 设置超时时间
func SetTimeout(d time.Duration) {
	defaultClient.httpClient.Timeout = d
}
