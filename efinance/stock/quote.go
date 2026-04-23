package stock

import (
	"context"
	"strconv"
	"strings"

	"github.com/T1anjiu/efinance-go/efinance/common"
    "github.com/T1anjiu/efinance-go/efinance/errors"
)

// QuoteParams 行情查询参数
type QuoteParams struct {
	Markets []string // 市场列表，如 "stock", "futures", "ETF"
}

// QuoteItem 实时行情数据项
type QuoteItem struct {
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	ChangePCT      float64 `json:"change_pct"`
	LatestPrice    float64 `json:"latest_price"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Open           float64 `json:"open"`
	ChangeAmt      float64 `json:"change_amt"`
	TurnoverRate   float64 `json:"turnover_rate"`
	VolumeRatio    float64 `json:"volume_ratio"`
	DynamicPE      float64 `json:"dynamic_pe"`
	Volume         int64   `json:"volume"`
	Amount         float64 `json:"amount"`
	YesterdayClose float64 `json:"yesterday_close"`
	TotalMarketCap float64 `json:"total_market_cap"`
	FlowMarketCap  float64 `json:"flow_market_cap"`
	QuoteID        string  `json:"quote_id"`
	MarketType     string  `json:"market_type"`
	UpdateTime     string  `json:"update_time"`
}

// GetRealtimeQuotes 获取实时行情
func GetRealtimeQuotes(ctx context.Context, params QuoteParams) ([]QuoteItem, error) {
	// 使用腾讯行情API
	// 默认获取沪深A股
	fs := "sh600519,sh600036,sz000001"

	if len(params.Markets) > 0 {
		// 根据市场类型构建代码列表
		codes := []string{}
		for _, m := range params.Markets {
			switch m {
			case "沪深A股", "stock":
				codes = append(codes, "sh600519", "sz000001", "sh600036")
			case "上证A股":
				codes = append(codes, "sh600519", "sh600036")
			case "深证A股":
				codes = append(codes, "sz000001")
			default:
				codes = append(codes, "sh600519")
			}
		}
		fs = strings.Join(codes, ",")
	}

	// 调用腾讯实时行情API
	raw, err := common.DefaultClient().GetRaw(ctx, common.TencentRealtimeURL+fs, map[string]string{
		"Referer": "https://finance.qq.com/",
	})
	if err != nil {
		return nil, err
	}

	return parseTencentRealtimeResponse(raw)
}

// GetLatestQuote 获取指定股票的实时行情
func GetLatestQuote(ctx context.Context, codes []string) ([]QuoteItem, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	// 转换为腾讯格式
	txCodes := make([]string, len(codes))
	for i, code := range codes {
		prefix := "sh"
		code = strings.ToUpper(code)
		if len(code) == 6 {
			if strings.HasPrefix(code, "6") {
				prefix = "sh"
			} else {
				prefix = "sz"
			}
		}
		txCodes[i] = prefix + code
	}

	fs := strings.Join(txCodes, ",")

	raw, err := common.DefaultClient().GetRaw(ctx, common.TencentRealtimeURL+fs, map[string]string{
		"Referer": "https://finance.qq.com/",
	})
	if err != nil {
		return nil, err
	}

	return parseTencentRealtimeResponse(raw)
}

// parseTencentRealtimeResponse 解析腾讯实时行情响应
// 响应格式: v_sh600519="1~名称~代码~现价~昨收~今开~成交量~外盘~内盘~...";
func parseTencentRealtimeResponse(raw []byte) ([]QuoteItem, error) {
	content := string(raw)

	// 按行分割
	lines := strings.Split(content, "\n")

	items := make([]QuoteItem, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 解析格式: v_sh600519="1~名称~代码~...";
		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			continue
		}

		codePart := strings.TrimPrefix(parts[0], "v_")
		dataPart := strings.Trim(parts[1], ";\""+" ")
		fields := strings.Split(dataPart, "~")

		if len(fields) < 10 {
			continue
		}

		latestPrice, _ := strconv.ParseFloat(fields[3], 64)
		changePCT, _ := strconv.ParseFloat(fields[32], 64)
		high, _ := strconv.ParseFloat(fields[33], 64)
		low, _ := strconv.ParseFloat(fields[34], 64)
		open, _ := strconv.ParseFloat(fields[5], 64)

		// 提取代码
		code := codePart
		if strings.HasPrefix(codePart, "sh") {
			code = strings.TrimPrefix(codePart, "sh")
		} else if strings.HasPrefix(codePart, "sz") {
			code = strings.TrimPrefix(codePart, "sz")
		}

		items = append(items, QuoteItem{
			Code:        code,
			Name:        fields[1],
			LatestPrice: latestPrice,
			ChangePCT:   changePCT,
			High:        high,
			Low:         low,
			Open:        open,
		})
	}

	if len(items) == 0 {
		return nil, errors.ErrNoData
	}

	return items, nil
}
