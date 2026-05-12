package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/T1anjiu/efinance-go/efinance/common"
	"github.com/T1anjiu/efinance-go/efinance/errors"
)

// QuoteParams 行情查询参数
type QuoteParams struct {
	Markets []string // 市场列表，如 "沪深A股", "上证A股", "深证A股"
}

// GetRealtimeQuotes 获取全市场实时行情
func GetRealtimeQuotes(ctx context.Context, params QuoteParams) ([]common.QuoteItem, error) {
	fs := common.FSMarketDict["沪深A股"]

	if len(params.Markets) > 0 {
		selected := make([]string, 0)
		for _, m := range params.Markets {
			if val, ok := common.FSMarketDict[m]; ok {
				selected = append(selected, val)
			}
		}
		if len(selected) > 0 {
			fs = strings.Join(selected, ",")
		}
	}

	return fetchQuotesByFS(ctx, fs)
}

// GetLatestQuote 获取指定股票的实时行情
func GetLatestQuote(ctx context.Context, codes []string) ([]common.QuoteItem, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	secids := make([]string, len(codes))
	for i, code := range codes {
		secids[i] = common.GetSecid(code)
	}

	fields := "f2,f3,f4,f5,f6,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f37"

	queryParams := map[string]string{
		"fltt":  "2",
		"invt":  "2",
		"fields": fields,
		"secids": strings.Join(secids, ","),
	}

	raw, err := common.DefaultClient().GetJSON(ctx, common.EastMoneyRealTimeURL, queryParams, common.HTTPHeaders)
	if err != nil {
		return nil, err
	}

	return parseQuoteResponse(raw)
}

// fetchQuotesByFS 通过fs参数获取行情
func fetchQuotesByFS(ctx context.Context, fs string) ([]common.QuoteItem, error) {
	fields := "f2,f3,f4,f5,f6,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f37"

	queryParams := map[string]string{
		"pn":    "1",
		"pz":    "10000",
		"po":    "1",
		"np":    "1",
		"fltt": "2",
		"invt": "2",
		"fid":  "f12",
		"fs":   fs,
		"fields": fields,
	}

	raw, err := common.DefaultClient().GetJSON(ctx, common.EastMoneyQuoteURL, queryParams, common.HTTPHeaders)
	if err != nil {
		return nil, err
	}

	return parseQuoteResponse(raw)
}

// parseQuoteResponse 解析东方财富行情响应
func parseQuoteResponse(raw *json.RawMessage) ([]common.QuoteItem, error) {
	var resp struct {
		Data struct {
			Diff []map[string]interface{} `json:"diff"`
		} `json:"data"`
	}

	// 东方财富 API 返回的可能带有 rc/rt 字段
	var rawResp struct {
		Data struct {
			Diff []map[string]interface{} `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(*raw, &rawResp); err != nil {
		return nil, fmt.Errorf("%w: JSON解析失败 - %v", errors.ErrParse, err)
	}

	_ = resp
	items := make([]common.QuoteItem, 0, len(rawResp.Data.Diff))

	for _, d := range rawResp.Data.Diff {
		code := getStringField(d, "f12")
		name := getStringField(d, "f14")
		marketNum := getStringField(d, "f13")

		marketType := "未知"
		if mt, ok := common.MarketNumberDict[marketNum]; ok {
			marketType = mt
		}

		quoteID := marketNum + "." + code

		items = append(items, common.QuoteItem{
			Code:           code,
			Name:           name,
			ChangePCT:      getFloatField(d, "f3"),
			LatestPrice:    getFloatField(d, "f2"),
			High:           getFloatField(d, "f15"),
			Low:            getFloatField(d, "f16"),
			Open:           getFloatField(d, "f17"),
			ChangeAmt:      getFloatField(d, "f4"),
			TurnoverRate:   getFloatField(d, "f8"),
			VolumeRatio:    getFloatField(d, "f10"),
			DynamicPE:      getFloatField(d, "f9"),
			Volume:         getInt64Field(d, "f5"),
			Amount:         getFloatField(d, "f6"),
			YesterdayClose: getFloatField(d, "f18"),
			TotalMarketCap: getFloatField(d, "f20"),
			FlowMarketCap:  getFloatField(d, "f21"),
			QuoteID:        quoteID,
			MarketType:     marketType,
		})
	}

	if len(items) == 0 {
		return nil, errors.ErrNoData
	}

	return items, nil
}

func getStringField(d map[string]interface{}, key string) string {
	if v, ok := d[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		if f, ok := v.(float64); ok {
			if f == float64(int64(f)) {
				return strconv.FormatInt(int64(f), 10)
			}
			return strconv.FormatFloat(f, 'f', -1, 64)
		}
	}
	return ""
}

func getFloatField(d map[string]interface{}, key string) float64 {
	if v, ok := d[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
		if s, ok := v.(string); ok {
			f, _ := strconv.ParseFloat(s, 64)
			return f
		}
		if n, ok := v.(json.Number); ok {
			f, _ := n.Float64()
			return f
		}
	}
	return 0
}

func getInt64Field(d map[string]interface{}, key string) int64 {
	if v, ok := d[key]; ok {
		if f, ok := v.(float64); ok {
			return int64(f)
		}
		if s, ok := v.(string); ok {
			n, _ := strconv.ParseInt(s, 10, 64)
			return n
		}
		if n, ok := v.(json.Number); ok {
			i, _ := n.Int64()
			return i
		}
	}
	return 0
}