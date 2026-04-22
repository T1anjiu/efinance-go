package fund

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/efinance/efinance/efinance/common"
	"github.com/efinance/efinance/efinance/errors"
)

// FundHeaders 天天基金请求头
var FundHeaders = http.Header{
	"User-Agent":       []string{"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15"},
	"Referer":         []string{"https://fund.eastmoney.com"},
	"Accept":          []string{"*/*"},
	"Accept-Language": []string{"zh-CN,zh;q=0.9"},
}

// NetValueItem 基金净值数据项
type NetValueItem struct {
	Date        string  `json:"date"`         // 日期
	NAV         float64 `json:"nav"`          // 单位净值
	ACCNAV      float64 `json:"acc_nav"`      // 累计净值
	ChangePCT   float64 `json:"change_pct"`   // 涨跌幅
}

// GetQuoteHistory 获取基金历史净值
func GetQuoteHistory(ctx context.Context, fundCode string, pageSize int) ([]NetValueItem, error) {
	if pageSize <= 0 {
		pageSize = 40000 // 默认获取全部
	}

	data := url.Values{
		"FCODE":       []string{fundCode},
		"IsShareNet":  []string{"true"},
		"MobileKey":   []string{"1"},
		"appType":     []string{"ttjj"},
		"appVersion":  []string{"6.2.8"},
		"cToken":      []string{"1"},
		"deviceid":    []string{"1"},
		"pageIndex":   []string{"1"},
		"pageSize":    []string{strconv.Itoa(pageSize)},
		"plat":        []string{"Iphone"},
		"product":     []string{"EFund"},
		"serverVersion": []string{"6.2.8"},
		"uToken":      []string{"1"},
		"userId":      []string{"1"},
		"version":     []string{"6.2.8"},
	}

	raw, err := common.DefaultClient().PostForm(ctx, common.FundHistoryURL, data, FundHeaders)
	if err != nil {
		return nil, err
	}

	return parseNetValueResponse(raw)
}

// GetQuoteHistoryMulti 获取多只基金历史净值
func GetQuoteHistoryMulti(ctx context.Context, fundCodes []string, pageSize int) (map[string][]NetValueItem, error) {
	type result struct {
		code  string
		items []NetValueItem
		err   error
	}

	results := make(chan result, len(fundCodes))
	
	var wg sync.WaitGroup
	sem := make(chan struct{}, common.MaxConnections)

	for _, code := range fundCodes {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			items, err := GetQuoteHistory(ctx, code, pageSize)
			results <- result{code: code, items: items, err: err}
		}(code)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	data := make(map[string][]NetValueItem)
	var errs []error

	for r := range results {
		if r.err != nil {
			errs = append(errs, r.err)
		} else {
			data[r.code] = r.items
		}
	}

	if len(errs) > 0 {
		return data, errs[0]
	}

	return data, nil
}

// parseNetValueResponse 解析基金净值响应数据
func parseNetValueResponse(raw *json.RawMessage) ([]NetValueItem, error) {
	var resp struct {
		Datas []struct {
			FSRQ   string `json:"FSRQ"`   // 日期
			DWjz   string `json:"DWJZ"`   // 单位净值
			Ljjz   string `json:"LJJZ"`   // 累计净值
			Jzzzl  string `json:"JZZZL"`  // 涨跌幅
		} `json:"Datas"`
	}

	if err := json.Unmarshal(*raw, &resp); err != nil {
		return nil, errors.ErrParse
	}

	items := make([]NetValueItem, 0, len(resp.Datas))
	for _, d := range resp.Datas {
		nav, _ := strconv.ParseFloat(d.DWjz, 64)
		accNav, _ := strconv.ParseFloat(d.Ljjz, 64)
		changePCT, _ := strconv.ParseFloat(d.Jzzzl, 64)

		items = append(items, NetValueItem{
			Date:      d.FSRQ,
			NAV:       nav,
			ACCNAV:    accNav,
			ChangePCT: changePCT,
		})
	}

	return items, nil
}

// EstNAVItem 基金估算净值数据项
type EstNAVItem struct {
	Code         string  `json:"code"`          // 基金代码
	Name         string  `json:"name"`          // 基金名称
	NAV          float64 `json:"nav"`           // 最新净值
	PublishDate  string  `json:"publish_date"`  // 最新净值公开日期
	EstTime      string  `json:"est_time"`      // 估算时间
	EstChangePCT float64 `json:"est_change_pct"` // 估算涨跌幅
}

// GetRealtimeEstimate 获取基金实时估算涨跌
func GetRealtimeEstimate(ctx context.Context, fundCodes []string) ([]EstNAVItem, error) {
	if len(fundCodes) == 0 {
		return nil, nil
	}

	codesStr := strings.Join(fundCodes, ",")
	
	data := url.Values{
		"pageIndex":   []string{"1"},
		"pageSize":    []string{"300000"},
		"Sort":        []string{""},
		"Fcodes":      []string{codesStr},
		"SortColumn":  []string{""},
		"IsShowSE":    []string{"false"},
		"P":           []string{"F"},
		"deviceid":    []string{"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"plat":        []string{"Iphone"},
		"product":     []string{"EFund"},
		"version":     []string{"6.2.8"},
	}

	raw, err := common.DefaultClient().PostForm(ctx, common.FundRealtimeURL, data, FundHeaders)
	if err != nil {
		return nil, err
	}

	return parseEstNAVResponse(raw)
}

// parseEstNAVResponse 解析估算净值响应数据
func parseEstNAVResponse(raw *json.RawMessage) ([]EstNAVItem, error) {
	var resp struct {
		Datas []struct {
			FCODE    string `json:"FCODE"`
			SHORTNAME string `json:"SHORTNAME"`
			ACCNAV   string `json:"ACCNAV"`
			PDATE    string `json:"PDATE"`
			GZTIME   string `json:"GZTIME"`
			GSZZL    string `json:"GSZZL"`
		} `json:"Datas"`
	}

	if err := json.Unmarshal(*raw, &resp); err != nil {
		return nil, errors.ErrParse
	}

	items := make([]EstNAVItem, 0, len(resp.Datas))
	for _, d := range resp.Datas {
		accNav, _ := strconv.ParseFloat(d.ACCNAV, 64)
		changePCT, _ := strconv.ParseFloat(d.GSZZL, 64)

		items = append(items, EstNAVItem{
			Code:         d.FCODE,
			Name:         d.SHORTNAME,
			NAV:          accNav,
			PublishDate:  d.PDATE,
			EstTime:     d.GZTIME,
			EstChangePCT: changePCT,
		})
	}

	return items, nil
}
