package common

import (
	"net/http"
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

	// 请求超时时间
	RequestTimeout = 180 * time.Second

	// 最大并发数
	MaxConnections = 20

	// 重试次数
	MaxRetries = 3

	// 默认日期范围
	DefaultBegDate = "19000101"
	DefaultEndDate = "20500101"
)

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

// QuoteIDCache 行情ID缓存
type QuoteIDCache struct {
	mu    sync.RWMutex
	cache map[string]string // code -> quoteId
}

var quoteIDCache = &QuoteIDCache{
	cache: make(map[string]string),
}

// GetQuoteID 获取缓存的行情ID
func (c *QuoteIDCache) Get(code string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id, ok := c.cache[code]
	return id, ok
}

// SetQuoteID 设置行情ID缓存
func (c *QuoteIDCache) Set(code, quoteID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[code] = quoteID
}

// SearchCache 搜索结果缓存
type SearchCache struct {
	mu    sync.RWMutex
	cache map[string]SearchResult
}

// SearchResult 搜索结果
type SearchResult struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	QuoteID  string `json:"quote_id"`
	MktType  string `json:"mkt_type"`
	SecuType string `json:"secu_type"`
}

var searchCache = &SearchCache{
	cache: make(map[string]SearchResult),
}

// Get 获取缓存的搜索结果
func (c *SearchCache) Get(key string) (SearchResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result, ok := c.cache[key]
	return result, ok
}

// Set 设置搜索结果缓存
func (c *SearchCache) Set(key string, result SearchResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = result
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
