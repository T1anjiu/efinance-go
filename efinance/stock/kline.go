package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/efinance/efinance/efinance/common"
	"github.com/efinance/efinance/efinance/errors"
)

// GetKlineParams K线查询参数
type GetKlineParams struct {
	Code        string            // 股票代码或名称
	Beg         string            // 开始日期 YYYYMMDD
	End         string            // 结束日期 YYYYMMDD
	KlineType   common.KlineType // K线周期
	AdjustType  common.AdjsType   // 复权类型
	MarketType  common.MarketType // 市场类型筛选
	SuppressErr bool             // 遇到错误是否静默
}

// KlineResult K线查询结果
type KlineResult struct {
	Code  string             // 股票代码
	Name  string             // 股票名称
	Items []common.KlineItem // K线数据
}

// GetKline 获取单只股票K线数据
func GetKline(ctx context.Context, params GetKlineParams) (*KlineResult, error) {
	if params.Beg == "" {
		params.Beg = "2024-01-01"
	}
	if params.End == "" {
		params.End = ""
	}
	if params.KlineType == 0 {
		params.KlineType = common.KlineDaily
	}

	code := params.Code

	// 如果不是市场前缀格式，需要添加
	marketPrefix := getMarketPrefix(code)
	apiCode := code
	if !strings.Contains(code, "_") && !strings.HasPrefix(code, "sh") && !strings.HasPrefix(code, "sz") {
		apiCode = marketPrefix + code
	}

	// 获取K线周期字符串
	klinePeriod := getKlinePeriod(params.KlineType)
	if klinePeriod == "" {
		klinePeriod = "day"
	}

	// 获取复权类型
	fq := "qfq"
	if params.AdjustType == common.AdjsNone {
		fq = ""
	} else if params.AdjustType == common.AdjsFront {
		fq = "qfq"
	} else if params.AdjustType == common.AdjsBack {
		fq = "hfq"
	}

	// 构建腾讯API URL
	url := fmt.Sprintf("%s?_var=kline_%s&param=%s,%s,%s,%s,100,%s",
		common.TencentKlineURL,
		klinePeriod+fq,
		strings.ToLower(apiCode),
		klinePeriod,
		params.Beg,
		params.End,
		fq,
	)

	raw, err := common.DefaultClient().GetRaw(ctx, url, map[string]string{
		"Referer": "https://finance.qq.com/",
	})
	if err != nil {
		return nil, err
	}

	return parseTencentKlineResponse(raw, code)
}

// GetKlineMulti 获取多只股票K线数据
func GetKlineMulti(ctx context.Context, params []GetKlineParams, workers int) (map[string]*KlineResult, error) {
	if workers <= 0 {
		workers = common.MaxConnections
	}

	type result struct {
		code  string
		kline *KlineResult
		err   error
	}

	results := make(chan result, len(params))

	var wg sync.WaitGroup
	sem := make(chan struct{}, workers)

	for _, p := range params {
		wg.Add(1)
		go func(p GetKlineParams) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			kline, err := GetKline(ctx, p)
			results <- result{code: p.Code, kline: kline, err: err}
		}(p)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	klines := make(map[string]*KlineResult)
	var errs []error

	for r := range results {
		if r.err != nil {
			errs = append(errs, r.err)
		} else {
			klines[r.code] = r.kline
		}
	}

	if len(errs) > 0 {
		return klines, errs[0]
	}

	return klines, nil
}

// parseTencentKlineResponse 解析腾讯K线响应
// 响应格式: var kline_dayqfq={...};
func parseTencentKlineResponse(raw []byte, code string) (*KlineResult, error) {
	// 去掉JSONP包装: var xxx={...};
	jsonStr := string(raw)

	// 找到第一个 { 的位置
	startIdx := strings.Index(jsonStr, "{")
	if startIdx == -1 {
		return nil, errors.ErrParse
	}

	// 找到最后一个 } 的位置
	endIdx := strings.LastIndex(jsonStr, "}")
	if endIdx == -1 {
		return nil, errors.ErrParse
	}

	jsonStr = jsonStr[startIdx : endIdx+1]

	// 解析为通用map结构
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return nil, errors.ErrParse
	}

	// 获取data字段
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		return nil, errors.ErrParse
	}

	// 获取第一个股票数据（key是动态的）
	var days [][]interface{}
	var name string
	for key, val := range data {
		if dayData, ok := val.(map[string]interface{}); ok {
			if dayArr, ok := dayData["day"].([]interface{}); ok {
				days = make([][]interface{}, len(dayArr))
				for i, d := range dayArr {
					if arr, ok := d.([]interface{}); ok {
						days[i] = arr
					}
				}
				name = key
				break
			}
		}
	}

	if len(days) == 0 {
		return nil, errors.ErrNoData
	}

	items := make([]common.KlineItem, 0, len(days))
	for _, d := range days {
		if len(d) < 6 {
			continue
		}

		dateStr := toString(d[0])
		open, _ := strconv.ParseFloat(toString(d[1]), 64)
		close, _ := strconv.ParseFloat(toString(d[2]), 64)
		high, _ := strconv.ParseFloat(toString(d[3]), 64)
		low, _ := strconv.ParseFloat(toString(d[4]), 64)
		volume, _ := strconv.ParseFloat(toString(d[5]), 64)

		var amount float64
		if len(d) > 6 {
			amount, _ = strconv.ParseFloat(toString(d[6]), 64)
		}

		items = append(items, common.KlineItem{
			Code:   code,
			Name:   name,
			Date:   dateStr,
			Open:   open,
			Close:  close,
			High:   high,
			Low:    low,
			Volume: volume,
			Amount: amount,
		})
	}

	return &KlineResult{
		Code:  code,
		Name:  name,
		Items: items,
	}, nil
}

// getMarketPrefix 根据股票代码获取市场前缀
func getMarketPrefix(code string) string {
	code = strings.ToUpper(code)
	if len(code) == 6 {
		switch {
		case strings.HasPrefix(code, "6"):
			return "sh"
		case strings.HasPrefix(code, "0"), strings.HasPrefix(code, "3"):
			return "sz"
		case strings.HasPrefix(code, "4"), strings.HasPrefix(code, "8"):
			return "bj"
		}
	}
	return "sh"
}

// getKlinePeriod 将K线类型转换为腾讯API的周期字符串
func getKlinePeriod(klt common.KlineType) string {
	switch klt {
	case common.KlineDaily:
		return "day"
	case common.KlineWeekly:
		return "week"
	case common.KlineMonthly:
		return "month"
	case common.Kline1Min:
		return "1min"
	case common.Kline5Min:
		return "5min"
	case common.Kline15Min:
		return "15min"
	case common.Kline30Min:
		return "30min"
	case common.Kline60Min:
		return "60min"
	default:
		return "day"
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", val)
	}
}
