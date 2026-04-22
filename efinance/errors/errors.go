package errors

import (
	"errors"
	"fmt"
)

// 定义错误类型
var (
	// ErrNoData 无数据
	ErrNoData = errors.New("无数据")

	// ErrInvalidCode 无效的股票/基金代码
	ErrInvalidCode = errors.New("无效的代码")

	// ErrInvalidDate 无效的日期格式
	ErrInvalidDate = errors.New("无效的日期格式")

	// ErrNetwork 网络请求错误
	ErrNetwork = errors.New("网络请求错误")

	// ErrParse 解析数据错误
	ErrParse = errors.New("解析数据错误")

	// ErrCache 缓存错误
	ErrCache = errors.New("缓存错误")

	// ErrTimeout 超时
	ErrTimeout = errors.New("请求超时")

	// ErrRateLimit 限流
	ErrRateLimit = errors.New("请求过于频繁，请稍后重试")
)

// APIError API错误
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API错误 [%d]: %s (URL: %s)", e.Code, e.Message, e.URL)
}

// NewAPIError 创建API错误
func NewAPIError(code int, message, url string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		URL:     url,
	}
}

// IsRetryable 判断错误是否可重试
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// 网络错误、超时、限流可以重试
	return errors.Is(err, ErrNetwork) ||
		errors.Is(err, ErrTimeout) ||
		errors.Is(err, ErrRateLimit)
}

// IsNotRetryable 判断错误是否不可重试
func IsNotRetryable(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrInvalidCode) ||
		errors.Is(err, ErrInvalidDate) ||
		errors.Is(err, ErrNoData)
}
