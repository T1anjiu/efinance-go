package quote

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/T1anjiu/efinance-go/internal/cache"
	"github.com/T1anjiu/efinance-go/internal/client"
	"github.com/tidwall/gjson"
)

type Quote struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Pinyin    string `json:"pinyin"`
	ID        string `json:"id"`
	JYS       string `json:"jys"`
	Classify  string `json:"classify"`
	MarketType string `json:"market_type"`
	SecurityTypeName string `json:"security_type_name"`
	SecurityType    string `json:"security_type"`
	MktNum    string `json:"mkt_num"`
	TypeUS    string `json:"type_us"`
	QuoteID   string `json:"quote_id"`
	UnifiedCode string `json:"unified_code"`
	InnerCode string `json:"inner_code"`
}

type Option func(*searchOpts)

type searchOpts struct {
	marketType   string
	symbolMode   bool
	useCache     bool
	suppressError bool
}

func WithMarketType(mt string) Option {
	return func(o *searchOpts) { o.marketType = mt }
}

func WithSymbolMode(v bool) Option {
	return func(o *searchOpts) { o.symbolMode = v }
}

func WithCache(v bool) Option {
	return func(o *searchOpts) { o.useCache = v }
}

func WithSuppressError(v bool) Option {
	return func(o *searchOpts) { o.suppressError = v }
}

func GetQuoteID(code string, opts ...Option) (string, error) {
	o := &searchOpts{useCache: true}
	for _, opt := range opts {
		opt(o)
	}
	code = strings.TrimSpace(code)
	if len(code) == 0 {
		if o.suppressError {
			return "", nil
		}
		return "", fmt.Errorf("证券代码长度不应为0")
	}

	q, err := Search(code, opts...)
	if err != nil {
		if o.suppressError {
			return "", nil
		}
		return "", fmt.Errorf("证券代码 %q 可能有误: %w", code, err)
	}
	if q == nil {
		if o.suppressError {
			return "", nil
		}
		return "", fmt.Errorf("证券代码 %q 可能有误", code)
	}
	return q.QuoteID, nil
}

func Search(keyword string, opts ...Option) (*Quote, error) {
	o := &searchOpts{useCache: true}
	for _, opt := range opts {
		opt(o)
	}

	if o.useCache {
		if entry, ok := cache.Get(keyword); ok {
			if o.marketType == "" || entry.Classify == o.marketType {
				return &Quote{
					Code:     entry.Code,
					Name:     entry.Name,
					QuoteID:  entry.QuoteID,
					Classify: entry.Classify,
				}, nil
			}
		}
	}

	params := url.Values{
		"input": {keyword},
		"type":  {"14"},
		"token": {"D43BF722C8E33BDC906FB84D85E326E8"},
		"count": {"5"},
	}

	u := "https://searchapi.eastmoney.com/api/suggest/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	result := gjson.ParseBytes(body)
	table := result.Get("QuotationCodeTable.Data")
	if !table.Exists() || table.IsArray() == false {
		return nil, nil
	}

	var quotes []*Quote
	table.ForEach(func(_, item gjson.Result) bool {
		code := item.Get("Code").String()
		classify := item.Get("Classify").String()

		if o.symbolMode && keyword != code {
			return true
		}
		if o.marketType != "" && classify != o.marketType {
			return true
		}

		q := &Quote{
			Code:      code,
			Name:      item.Get("Name").String(),
			Pinyin:    item.Get("Pinyin").String(),
			ID:        item.Get("ID").String(),
			JYS:       item.Get("JYS").String(),
			Classify:  classify,
			MarketType: item.Get("MarketType").String(),
			SecurityTypeName: item.Get("SecurityTypeName").String(),
			SecurityType:     item.Get("SecurityType").String(),
			MktNum:    item.Get("MktNum").String(),
			TypeUS:    item.Get("Type_US").String(),
			QuoteID:   item.Get("MktNum").String() + "." + code,
			UnifiedCode: item.Get("UnifiedCode").String(),
			InnerCode: item.Get("InnerCode").String(),
		}
		quotes = append(quotes, q)
		return true
	})

	if len(quotes) == 0 {
		return nil, nil
	}

	first := quotes[0]
	cache.Set(keyword, cache.CacheEntry{
		QuoteID:  first.QuoteID,
		Name:     first.Name,
		Code:     first.Code,
		Classify: first.Classify,
	})

	return first, nil
}

func SearchMulti(keyword string, count int, opts ...Option) ([]*Quote, error) {
	o := &searchOpts{useCache: true}
	for _, opt := range opts {
		opt(o)
	}

	if count < 5 {
		count = 5
	}

	params := url.Values{
		"input": {keyword},
		"type":  {"14"},
		"token": {"D43BF722C8E33BDC906FB84D85E326E8"},
		"count": {fmt.Sprintf("%d", count)},
	}

	u := "https://searchapi.eastmoney.com/api/suggest/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	table := result.Get("QuotationCodeTable.Data")
	if !table.Exists() {
		return nil, nil
	}

	var quotes []*Quote
	table.ForEach(func(_, item gjson.Result) bool {
		code := item.Get("Code").String()
		classify := item.Get("Classify").String()

		if o.symbolMode && keyword != code {
			return true
		}
		if o.marketType != "" && classify != o.marketType {
			return true
		}

		q := &Quote{
			Code:      code,
			Name:      item.Get("Name").String(),
			Pinyin:    item.Get("Pinyin").String(),
			ID:        item.Get("ID").String(),
			JYS:       item.Get("JYS").String(),
			Classify:  classify,
			MarketType: item.Get("MarketType").String(),
			SecurityTypeName: item.Get("SecurityTypeName").String(),
			SecurityType:     item.Get("SecurityType").String(),
			MktNum:    item.Get("MktNum").String(),
			TypeUS:    item.Get("Type_US").String(),
			QuoteID:   item.Get("MktNum").String() + "." + code,
			UnifiedCode: item.Get("UnifiedCode").String(),
			InnerCode: item.Get("InnerCode").String(),
		}
		quotes = append(quotes, q)
		return true
	})

	if len(quotes) > 0 {
		cache.Set(keyword, cache.CacheEntry{
			QuoteID:  quotes[0].QuoteID,
			Name:     quotes[0].Name,
			Code:     quotes[0].Code,
			Classify: quotes[0].Classify,
		})
	}

	if len(quotes) > count {
		return quotes[:count], nil
	}
	return quotes, nil
}

func init() {
	_ = json.Unmarshal
	_ = cache.Load
}
