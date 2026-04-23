package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/T1anjiu/efinance-go/efinance/errors"
)

// HTTPClient HTTP客户端
type HTTPClient struct {
	client  *http.Client
	mu      sync.Mutex
	retries int
}

var defaultClient *HTTPClient

func init() {
	defaultClient = NewHTTPClient(MaxConnections, MaxRetries)
}

// NewHTTPClient 创建HTTP客户端
func NewHTTPClient(maxConn, retries int) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        maxConn,
				MaxIdleConnsPerHost: maxConn,
				IdleConnTimeout:     90 * time.Second,
			},
			Timeout: RequestTimeout,
		},
		retries: retries,
	}
}

// GetJSON GET请求JSON
func (c *HTTPClient) GetJSON(ctx context.Context, url string, params map[string]string, headers http.Header) (*json.RawMessage, error) {
	c.mu.Lock()
	c.client.Timeout = RequestTimeout
	c.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加查询参数
	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// 添加请求头
	if headers != nil {
		for k, values := range headers {
			for _, v := range values {
				req.Header.Add(k, v)
			}
		}
	}

	// 添加默认请求头
	for k, values := range HTTPHeaders {
		if req.Header.Get(k) == "" {
			for _, v := range values {
				req.Header.Add(k, v)
			}
		}
	}

	// 发送请求并重试
	var lastErr error
	for i := 0; i <= c.retries; i++ {
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("请求失败: %w", err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("读取响应失败: %w", err)
			continue
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			lastErr = fmt.Errorf("解析JSON失败: %w", err)
			continue
		}

		// 检查API返回的错误
		if code, ok := result["code"].(float64); ok && code != 200 {
			if msg, ok := result["message"].(string); ok {
				return nil, fmt.Errorf("API错误: %s", msg)
			}
		}

		data, ok := result["data"]
		if !ok {
			return nil, errors.ErrNoData
		}

		raw, ok := data.(json.RawMessage)
		if !ok {
			return nil, errors.ErrNoData
		}

		return &raw, nil
	}

	return nil, lastErr
}

// PostJSON POST请求JSON
func (c *HTTPClient) PostJSON(ctx context.Context, url string, data interface{}, headers http.Header) (*json.RawMessage, error) {
	c.mu.Lock()
	c.client.Timeout = RequestTimeout
	c.mu.Unlock()

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("序列化请求数据失败: %w", err)
		}
		body = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 添加请求头
	if headers != nil {
		for k, values := range headers {
			for _, v := range values {
				req.Header.Add(k, v)
			}
		}
	}

	// 添加默认请求头
	for k, values := range HTTPHeaders {
		if req.Header.Get(k) == "" {
			for _, v := range values {
				req.Header.Add(k, v)
			}
		}
	}

	// 发送请求并重试
	var lastErr error
	for i := 0; i <= c.retries; i++ {
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("请求失败: %w", err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("读取响应失败: %w", err)
			continue
		}

		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			lastErr = fmt.Errorf("解析JSON失败: %w", err)
			continue
		}

		data, ok := result["data"]
		if !ok {
			return nil, errors.ErrNoData
		}

		raw, ok := data.(json.RawMessage)
		if !ok {
			return nil, errors.ErrNoData
		}

		return &raw, nil
	}

	return nil, lastErr
}

// PostForm postForm请求
func (c *HTTPClient) PostForm(ctx context.Context, url string, data url.Values, headers http.Header) (*json.RawMessage, error) {
	c.mu.Lock()
	c.client.Timeout = RequestTimeout
	c.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 添加请求头
	if headers != nil {
		for k, values := range headers {
			for _, v := range values {
				req.Header.Add(k, v)
			}
		}
	}

	// 添加默认请求头
	for k, values := range HTTPHeaders {
		if req.Header.Get(k) == "" {
			for _, v := range values {
				req.Header.Add(k, v)
			}
		}
	}

	// 发送请求并重试
	var lastErr error
	for i := 0; i <= c.retries; i++ {
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("请求失败: %w", err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("读取响应失败: %w", err)
			continue
		}

		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			lastErr = fmt.Errorf("解析JSON失败: %w", err)
			continue
		}

		data, ok := result["data"]
		if !ok {
			return nil, errors.ErrNoData
		}

		raw, ok := data.(json.RawMessage)
		if !ok {
			return nil, errors.ErrNoData
		}

		return &raw, nil
	}

	return nil, lastErr
}

// DefaultClient 获取默认HTTP客户端
func DefaultClient() *HTTPClient {
	return defaultClient
}

// GetRaw GET请求返回原始字节数据（不解析JSON）
func (c *HTTPClient) GetRaw(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	c.mu.Lock()
	c.client.Timeout = RequestTimeout
	c.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加默认请求头
	for k, values := range HTTPHeaders {
		req.Header.Set(k, values[0])
	}
	
	// 添加自定义请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// Do 直接发送请求（用于调试）
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// NewRequest 创建新的请求（用于调试）
func (c *HTTPClient) NewRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}
