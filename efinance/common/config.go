package common

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 常量配置
const (
	// 腾讯行情API URLs
	TencentRealtimeURL = "https://qt.gtimg.cn/q="
	TencentKlineURL    = "https://web.ifzq.gtimg.cn/appstock/app/fqkline/get"

	// 东方财富备用API (可能被防火墙拦截)
	EastMoneyKlineURL       = "https://push2his.eastmoney.com/api/qt/stock/kline/get"
	EastMoneyQuoteURL       = "https://push2.eastmoney.com/api/qt/ulist.np/get"
	EastMoneyRealTimeURL    = "https://push2.eastmoney.com/api/qt/ulist.np/get"
	EastMoneyDataCenterURL  = "http://datacenter-web.eastmoney.com/api/data/v1/get"
	EastMoneyEMH5URL       = "https://emh5.eastmoney.com/api"

	// 天天基金API
	FundHistoryURL   = "https://fundmobapi.eastmoney.com/FundMNewApi/FundMNHisNetList"
	FundRealtimeURL  = "https://fundmobapi.eastmoney.com/FundMNewApi/FundMNFInfo"
	FundListURL      = "https://fund.eastmoney.com/data/rankhandler"

	// 默认日期范围
	DefaultBegDate = "19000101"
	DefaultEndDate = "20500101"
)

// 从配置获取动态值
var (
	// RequestTimeout 请求超时时间
	RequestTimeout = getDurationEnv("REQUEST_TIMEOUT", 180*time.Second)

	// MaxConnections 最大并发数
	MaxConnections = getIntEnv("MAX_CONNECTIONS", 20)

	// MaxRetries 重试次数
	MaxRetries = getIntEnv("MAX_RETRIES", 3)

	// EastMoneySearchToken 东方财富搜索API Token
	EastMoneySearchToken = getEnv("EASTMONEY_SEARCH_TOKEN", "894050c76af8597a853f5b408b759f5d")

	// QuoteIDCacheTTL 行情ID缓存TTL
	QuoteIDCacheTTL = getDurationEnv("QUOTEID_CACHE_TTL", 24*time.Hour)

	// SearchCacheTTL 搜索缓存TTL
	SearchCacheTTL = getDurationEnv("SEARCH_CACHE_TTL", 1*time.Hour)
)

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv 获取整数环境变量
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getDurationEnv 获取时长环境变量
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// HTTPHeaders EastMoney请求头
var HTTPHeaders = http.Header{
	"User-Agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
	"Referer":         []string{"https://finance.eastmoney.com/"},
	"Accept":          []string{"*/*"},
	"Accept-Language": []string{"zh-CN,zh;q=0.9,en;q=0.8"},
	"Origin":          []string{"https://finance.eastmoney.com"},
}

// MarketNumberDict 市场编号映射
var MarketNumberDict = map[string]string{
	"0":   "深A",
	"1":   "沪A",
	"105": "美股",
	"106": "美股",
	"107": "美股",
	"116": "港股",
	"128": "港股",
	"113": "上期所",
	"114": "大商所",
	"115": "郑商所",
	"8":   "中金所",
	"142": "上海能源期货交易所",
	"155": "英股",
	"90":  "板块",
	"225": "广期所",
}

// FSMarketDict 市场筛选参数
var FSMarketDict = map[string]string{
	"bond":       "b:MK0354",
	"可转债":      "b:MK0354",
	"stock":      "m:0 t:6,m:0 t:80,m:1 t:2,m:1 t:23,m:0 t:81 s:2048",
	"沪深A股":     "m:0 t:6,m:0 t:80,m:1 t:2,m:1 t:23",
	"沪深京A股":   "m:0 t:6,m:0 t:80,m:1 t:2,m:1 t:23,m:0 t:81 s:2048",
	"北证A股":     "m:0 t:81 s:2048",
	"北A":        "m:0 t:81 s:2048",
	"futures":    "m:113,m:114,m:115,m:8,m:142,m:225",
	"期货":        "m:113,m:114,m:115,m:8,m:142,m:225",
	"上证A股":     "m:1 t:2,m:1 t:23",
	"沪A":        "m:1 t:2,m:1 t:23",
	"深证A股":     "m:0 t:6,m:0 t:80",
	"深A":        "m:0 t:6,m:0 t:80",
	"创业板":      "m:0 t:80",
	"科创板":      "m:1 t:23",
	"美股":        "m:105,m:106,m:107",
	"港股":        "m:116,m:128",
	"ETF":        "m:0 t:80,m:1 t:23",
	"LOF":        "m:0 t:80,m:1 t:23",
	"中概股":      "m:105,m:106,m:107",
	"新股":        "m:0 t:6,m:1 t:2",
	"沪股通":      "m:1 t:2",
	"深股通":      "m:0 t:80",
	"行业板块":    "m:90 t:2",
	"概念板块":    "m:90 t:10",
}

// SearchResult 搜索结果
type SearchResult struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	QuoteID  string `json:"quote_id"`
	MktType  string `json:"mkt_type"`
	SecuType string `json:"secu_type"`
}

// QuoteIDCache 行情ID缓存
type QuoteIDCache struct {
	mu       sync.RWMutex
	cache    map[string]*cacheEntry
	ttl      time.Duration
}

type cacheEntry struct {
	value      string
	expireTime time.Time
}

var quoteIDCache = &QuoteIDCache{
	cache: make(map[string]*cacheEntry),
	ttl:   QuoteIDCacheTTL,
}

// GetQuoteID 获取缓存的行情ID
func (c *QuoteIDCache) Get(code string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cache[code]
	if !ok {
		return "", false
	}

	// 检查是否过期
	if time.Now().After(entry.expireTime) {
		return "", false
	}

	return entry.value, true
}

// SetQuoteID 设置行情ID缓存
func (c *QuoteIDCache) Set(code, quoteID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[code] = &cacheEntry{
		value:      quoteID,
		expireTime: time.Now().Add(c.ttl),
	}
}

// Clear 清空缓存
func (c *QuoteIDCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]*cacheEntry)
}

// Cleanup 清理过期缓存
func (c *QuoteIDCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.cache {
		if now.After(entry.expireTime) {
			delete(c.cache, key)
		}
	}
}

// SearchCache 搜索结果缓存
type SearchCache struct {
	mu    sync.RWMutex
	cache map[string]*searchCacheEntry
	ttl   time.Duration
}

type searchCacheEntry struct {
	result     SearchResult
	expireTime time.Time
}

var searchCache = &SearchCache{
	cache: make(map[string]*searchCacheEntry),
	ttl:   SearchCacheTTL,
}

// Get 获取缓存的搜索结果
func (c *SearchCache) Get(key string) (SearchResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cache[key]
	if !ok {
		return SearchResult{}, false
	}

	// 检查是否过期
	if time.Now().After(entry.expireTime) {
		return SearchResult{}, false
	}

	return entry.result, true
}

// Set 设置搜索结果缓存
func (c *SearchCache) Set(key string, result SearchResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &searchCacheEntry{
		result:     result,
		expireTime: time.Now().Add(c.ttl),
	}
}

// Clear 清空缓存
func (c *SearchCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]*searchCacheEntry)
}

// Cleanup 清理过期缓存
func (c *SearchCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.cache {
		if now.After(entry.expireTime) {
			delete(c.cache, key)
		}
	}
}

// DefaultSearchCache 获取默认搜索缓存
func DefaultSearchCache() *SearchCache {
	return searchCache
}

// GetSecid 根据股票代码获取secid
func GetSecid(code string) string {
	code = strings.ToUpper(code)
	if len(code) == 6 {
		switch {
		case strings.HasPrefix(code, "6"):
			return "1." + code // 上交所
		case strings.HasPrefix(code, "0"), strings.HasPrefix(code, "3"):
			return "0." + code // 深交所
		case strings.HasPrefix(code, "4"), strings.HasPrefix(code, "8"):
			return "0." + code // 北交所
		}
	}
	return code
}
