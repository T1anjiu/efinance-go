package stock

import (
	"context"
	"encoding/json"

	"github.com/T1anjiu/efinance-go/efinance/common"
    "github.com/T1anjiu/efinance-go/efinance/errors"
)

// Search 搜索股票
func Search(ctx context.Context, keyword string, marketType common.MarketType) ([]common.SearchResult, error) {
	// 先检查缓存
	if result, ok := common.DefaultSearchCache().Get(keyword); ok {
		return []common.SearchResult{result}, nil
	}

	params := map[string]string{
		"InputTip": keyword,
		"Type":     "14",
		"Token":    common.EastMoneySearchToken,
		"MktNum":   "f1",
		"Idx":      "idx",
		"No":       "n1",
		"Src":      "010200000",
		"Count":    "20",
	}

	raw, err := common.DefaultClient().GetJSON(ctx, "https://searchapi.eastmoney.com/api/suggest/get", params, common.HTTPHeaders)
	if err != nil {
		return nil, err
	}

	var resp struct {
		QuotationList []struct {
			Code     string `json:"QuoteCode"`
			Name     string `json:"QuoteName"`
			MktType  string `json:"MktNum"`
			SecuType string `json:"SecuType"`
		} `json:"QuotationList"`
	}

	if err := json.Unmarshal(*raw, &resp); err != nil {
		return nil, errors.ErrParse
	}

	results := make([]common.SearchResult, 0, len(resp.QuotationList))
	for _, item := range resp.QuotationList {
		// 过滤市场类型
		if marketType != "" && marketType != common.MarketType(item.MktType) {
			continue
		}

		// 生成行情ID
		quoteID := getQuoteID(item.MktType, item.Code)

		result := common.SearchResult{
			Code:     item.Code,
			Name:     item.Name,
			QuoteID:  quoteID,
			MktType:  item.MktType,
			SecuType: item.SecuType,
		}
		results = append(results, result)

		// 缓存第一个结果（精确匹配）
		if len(results) == 1 {
			common.DefaultSearchCache().Set(keyword, result)
		}
	}

	return results, nil
}

// GetQuoteID 获取行情ID
func GetQuoteID(ctx context.Context, code string, marketType common.MarketType) (string, error) {
	// 先搜索
	results, err := Search(ctx, code, marketType)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", errors.ErrInvalidCode
	}

	return results[0].QuoteID, nil
}

// getQuoteID 生成行情ID
func getQuoteID(mktType, code string) string {
	switch mktType {
	case "1":
		return "1." + code // 沪A
	case "0":
		return "0." + code // 深A
	case "23":
		return "1." + code // 科创板
	case "116", "128":
		return mktType + "." + code // 港股
	case "105", "106", "107":
		return mktType + "." + code // 美股
	default:
		return mktType + "." + code
	}
}
