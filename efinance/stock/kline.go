package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/T1anjiu/efinance-go/efinance/common"
	"github.com/T1anjiu/efinance-go/efinance/errors"
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
		params.Beg = common.DefaultBegDate
	}
	if params.End == "" {
		params.End = common.DefaultEndDate
	}
	if params.KlineType == 0 {
		params.KlineType = common.KlineDaily
	}

	secid := common.GetSecid(params.Code)

	queryParams := map[string]string{
		"fields1": "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13",
		"fields2": "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61",
		"beg":     params.Beg,
		"end":     params.End,
		"rtntype": "6",
		"secid":   secid,
		"klt":     strconv.Itoa(int(params.KlineType)),
		"fqt":     strconv.Itoa(int(params.AdjustType)),
	}

	raw, err := common.DefaultClient().GetJSON(ctx, common.EastMoneyKlineURL, queryParams, common.HTTPHeaders)
	if err != nil {
		return nil, err
	}

	return parseEastMoneyKline(raw, secid)
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
		return klines, fmt.Errorf("批量请求 %d 个失败: %v", len(errs), errs)
	}

	return klines, nil
}

// parseEastMoneyKline 解析东方财富K线响应
func parseEastMoneyKline(raw *json.RawMessage, secid string) (*KlineResult, error) {
	var resp struct {
		Data struct {
			Code   string   `json:"code"`
			Name   string   `json:"name"`
			Klines []string `json:"klines"`
		} `json:"data"`
	}

	if err := json.Unmarshal(*raw, &resp); err != nil {
		return nil, fmt.Errorf("%w: JSON解析失败 - %v", errors.ErrParse, err)
	}

	if len(resp.Data.Klines) == 0 {
		return nil, errors.ErrNoData
	}

	code := strings.Split(secid, ".")[1]
	items := make([]common.KlineItem, 0, len(resp.Data.Klines))

	for _, kline := range resp.Data.Klines {
		fields := strings.Split(kline, ",")
		if len(fields) < 11 {
			continue
		}

		open, _ := strconv.ParseFloat(fields[1], 64)
		close, _ := strconv.ParseFloat(fields[2], 64)
		high, _ := strconv.ParseFloat(fields[3], 64)
		low, _ := strconv.ParseFloat(fields[4], 64)
		volume, _ := strconv.ParseFloat(fields[5], 64)
		amount, _ := strconv.ParseFloat(fields[6], 64)
		amplitude, _ := strconv.ParseFloat(fields[7], 64)
		changePct, _ := strconv.ParseFloat(fields[8], 64)
		changeAmt, _ := strconv.ParseFloat(fields[9], 64)
		turnoverRate, _ := strconv.ParseFloat(fields[10], 64)

		items = append(items, common.KlineItem{
			Code:         code,
			Name:         resp.Data.Name,
			Date:         fields[0],
			Open:         open,
			Close:        close,
			High:         high,
			Low:          low,
			Volume:       volume,
			Amount:       amount,
			Amplitude:    amplitude,
			ChangePCT:    changePct,
			ChangeAmt:    changeAmt,
			TurnoverRate: turnoverRate,
		})
	}

	return &KlineResult{
		Code:  code,
		Name:  resp.Data.Name,
		Items: items,
	}, nil
}