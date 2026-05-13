package common

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/T1anjiu/efinance-go/internal/client"
	"github.com/T1anjiu/efinance-go/internal/util"
	"github.com/T1anjiu/efinance-go/quote"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
)

type KlineBar struct {
	Date       string  `json:"date"`
	Open       float64 `json:"open"`
	Close      float64 `json:"close"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Volume     int64   `json:"volume"`
	Amount     float64 `json:"amount"`
	Amplitude  float64 `json:"amplitude"`
	ChangeRate float64 `json:"change_rate"`
	ChangeAmt  float64 `json:"change_amt"`
	Turnover   float64 `json:"turnover"`
}

type KlineResult struct {
	Code string     `json:"code"`
	Name string     `json:"name"`
	Bars []KlineBar `json:"bars"`
}

type KlineOptions struct {
	Begin string
	End   string
	KLT   int
	FQT   int
}

func DefaultKlineOptions() KlineOptions {
	return KlineOptions{
		Begin: "19000101",
		End:   "20500101",
		KLT:   101,
		FQT:   1,
	}
}

type RealtimeQuote struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	ChangeRate   float64 `json:"change_rate"`
	ChangeAmt    float64 `json:"change_amt"`
	Volume       int64   `json:"volume"`
	Amount       float64 `json:"amount"`
	High         float64 `json:"high"`
	Low          float64 `json:"low"`
	Open         float64 `json:"open"`
	PreClose     float64 `json:"pre_close"`
	TurnoverRate float64 `json:"turnover_rate"`
	PERatio      float64 `json:"pe_ratio"`
	PBRatio      float64 `json:"pb_ratio"`
	MarketCap    float64 `json:"market_cap"`
	CircMarketCap float64 `json:"circ_market_cap"`
	QuoteID      string  `json:"quote_id"`
	MarketType   string  `json:"market_type"`
	UpdateTimestamp int64 `json:"update_timestamp"`
	UpdateDate   string  `json:"update_date"`
	LatestTradeDate string `json:"latest_trade_date"`
}

type BillRecord struct {
	Date            string  `json:"date"`
	MainNetInflow   float64 `json:"main_net_inflow"`
	SmallNetInflow  float64 `json:"small_net_inflow"`
	MidNetInflow    float64 `json:"mid_net_inflow"`
	LargeNetInflow  float64 `json:"large_net_inflow"`
	HugeNetInflow   float64 `json:"huge_net_inflow"`
	MainNetInflowPct float64 `json:"main_net_inflow_pct"`
	SmallInflowPct  float64 `json:"small_inflow_pct"`
	MidInflowPct    float64 `json:"mid_inflow_pct"`
	LargeInflowPct  float64 `json:"large_inflow_pct"`
	HugeInflowPct   float64 `json:"huge_inflow_pct"`
	ClosePrice      float64 `json:"close_price"`
	ChangeRate      float64 `json:"change_rate"`
}

type TodayBillRecord struct {
	Time           string  `json:"time"`
	MainNetInflow  float64 `json:"main_net_inflow"`
	SmallNetInflow float64 `json:"small_net_inflow"`
	MidNetInflow   float64 `json:"mid_net_inflow"`
	LargeNetInflow float64 `json:"large_net_inflow"`
	HugeNetInflow  float64 `json:"huge_net_inflow"`
}

type BaseInfo map[string]interface{}

type DealRecord struct {
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	Time      string  `json:"time"`
	PreClose  float64 `json:"pre_close"`
	Price     float64 `json:"price"`
	Volume    int64   `json:"volume"`
	OrderNum  int64   `json:"order_num"`
}

type NDaysKlineBar struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume int64   `json:"volume"`
	Amount float64 `json:"amount"`
}

func GetQuoteHistorySingle(secid string, opts KlineOptions) (*KlineResult, error) {
	if opts.Begin == "" {
		opts.Begin = "19000101"
	}
	if opts.End == "" {
		opts.End = "20500101"
	}
	if opts.KLT == 0 {
		opts.KLT = 101
	}
	if opts.FQT == 0 {
		opts.FQT = 1
	}

	fields2 := strings.Join(KlineFields, ",")
	params := url.Values{
		"fields1": {"f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13"},
		"fields2": {fields2},
		"beg":     {opts.Begin},
		"end":     {opts.End},
		"rtntype": {"6"},
		"secid":   {secid},
		"klt":     {strconv.Itoa(opts.KLT)},
		"fqt":     {strconv.Itoa(opts.FQT)},
	}

	u := "https://push2his.eastmoney.com/api/qt/stock/kline/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	data := result.Get("data")
	if !data.Exists() {
		return &KlineResult{}, nil
	}

	name := data.Get("name").String()
	code := secid
	if idx := strings.LastIndex(secid, "."); idx >= 0 {
		code = secid[idx+1:]
	}

	klines := data.Get("klines")
	if !klines.Exists() || !klines.IsArray() {
		return &KlineResult{Code: code, Name: name}, nil
	}

	var bars []KlineBar
	klines.ForEach(func(_, v gjson.Result) bool {
		line := v.String()
		parts := strings.Split(line, ",")
		if len(parts) < 11 {
			return true
		}
		bars = append(bars, KlineBar{
			Date:       parts[0],
			Open:       util.ParseFloat(parts[1]),
			Close:      util.ParseFloat(parts[2]),
			High:       util.ParseFloat(parts[3]),
			Low:        util.ParseFloat(parts[4]),
			Volume:     util.ParseInt(parts[5]),
			Amount:     util.ParseFloat(parts[6]),
			Amplitude:  util.ParseFloat(parts[7]),
			ChangeRate: util.ParseFloat(parts[8]),
			ChangeAmt:  util.ParseFloat(parts[9]),
			Turnover:   util.ParseFloat(parts[10]),
		})
		return true
	})

	return &KlineResult{Code: code, Name: name, Bars: bars}, nil
}

func GetQuoteHistoryMulti(secids []string, opts KlineOptions) (map[string]*KlineResult, error) {
	results := make(map[string]*KlineResult, len(secids))
	var mu sync.Mutex

	eg, ctx := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, MaxConnections)

	for _, secid := range secids {
		sid := secid
		eg.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()

			r, err := GetQuoteHistorySingle(sid, opts)
			if err != nil {
				return nil
			}
			mu.Lock()
			results[sid] = r
			mu.Unlock()
			return nil
		})
	}

	_ = ctx
	if err := eg.Wait(); err != nil {
		return results, err
	}
	return results, nil
}

func GetRealtimeQuotesByFS(fs string) ([]RealtimeQuote, error) {
	fields := strings.Join(QuoteFieldKeys, ",")

	getByPage := func(pn, pz int) (gjson.Result, error) {
		params := url.Values{
			"pn":     {strconv.Itoa(pn)},
			"pz":     {strconv.Itoa(pz)},
			"po":     {"1"},
			"np":     {"1"},
			"fltt":   {"2"},
			"invt":   {"2"},
			"fid":    {"f12"},
			"fs":     {fs},
			"fields": {fields},
		}
		u := "http://push2.eastmoney.com/api/qt/clist/get"
		body, err := client.Get(u, params, nil)
		if err != nil {
			return gjson.Result{}, err
		}
		return gjson.ParseBytes(body), nil
	}

	firstResp, err := getByPage(1, 200)
	if err != nil {
		return nil, err
	}
	total := int(firstResp.Get("data.total").Int())
	diff := firstResp.Get("data.diff")
	if !diff.Exists() || !diff.IsArray() {
		return nil, nil
	}
	pz := len(diff.Array())
	if pz == 0 {
		pz = 50
	}
	pages := total / pz
	if total%pz != 0 {
		pages++
	}

	type pageResult struct {
		index int
		data  gjson.Result
	}

	results := make([]pageResult, pages)
	eg, _ := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, MaxConnections)

	for i := 1; i <= pages; i++ {
		idx := i
		eg.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()
			resp, err := getByPage(idx, pz)
			if err != nil {
				return err
			}
			results[idx-1] = pageResult{index: idx, data: resp}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	var quotes []RealtimeQuote
	for _, pr := range results {
		diff := pr.data.Get("data.diff")
		if !diff.Exists() {
			continue
		}
		diff.ForEach(func(_, item gjson.Result) bool {
			mktNum := item.Get("f13").String()
			code := item.Get("f12").String()
			quoteID := mktNum + "." + code
			marketType := MarketNumberDict[mktNum]

			q := RealtimeQuote{
				Code:       code,
				Name:       item.Get("f14").String(),
				Price:      util.ParseFloat(item.Get("f2").String()),
				ChangeRate: util.ParseFloat(item.Get("f3").String()),
				ChangeAmt:  util.ParseFloat(item.Get("f4").String()),
				Volume:     util.ParseInt(item.Get("f5").String()),
				Amount:     util.ParseFloat(item.Get("f6").String()),
				High:       util.ParseFloat(item.Get("f15").String()),
				Low:        util.ParseFloat(item.Get("f16").String()),
				Open:       util.ParseFloat(item.Get("f17").String()),
				PreClose:   util.ParseFloat(item.Get("f18").String()),
				TurnoverRate: util.ParseFloat(item.Get("f8").String()),
				PERatio:    util.ParseFloat(item.Get("f9").String()),
				MarketCap:  util.ParseFloat(item.Get("f20").String()),
				CircMarketCap: util.ParseFloat(item.Get("f21").String()),
				QuoteID:    quoteID,
				MarketType: marketType,
				UpdateTimestamp: item.Get("f124").Int(),
				LatestTradeDate: item.Get("f297").String(),
			}
			ts := q.UpdateTimestamp
			if ts > 0 {
				q.UpdateDate = time.Unix(ts, 0).Format("2006-01-02 15:04:05")
			}
			if q.LatestTradeDate != "" {
				if len(q.LatestTradeDate) == 8 {
					q.LatestTradeDate = q.LatestTradeDate[:4] + "-" + q.LatestTradeDate[4:6] + "-" + q.LatestTradeDate[6:8]
				}
			}
			quotes = append(quotes, q)
			return true
		})
	}

	sort.Slice(quotes, func(i, j int) bool {
		return quotes[i].ChangeRate > quotes[j].ChangeRate
	})

	return quotes, nil
}

func GetHistoryBill(secid string) (*KlineResult, []BillRecord, error) {
	fields2 := strings.Join(HistoryBillFields, ",")
	params := url.Values{
		"lmt":     {"100000"},
		"klt":     {"101"},
		"secid":   {secid},
		"fields1": {"f1,f2,f3,f7"},
		"fields2": {fields2},
	}
	u := "http://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, nil, err
	}

	result := gjson.ParseBytes(body)
	name := result.Get("data.name").String()
	code := secid
	if idx := strings.LastIndex(secid, "."); idx >= 0 {
		code = secid[idx+1:]
	}

	klines := result.Get("data.klines")
	if !klines.Exists() || !klines.IsArray() {
		return &KlineResult{Code: code, Name: name}, nil, nil
	}

	var records []BillRecord
	klines.ForEach(func(_, v gjson.Result) bool {
		parts := strings.Split(v.String(), ",")
		if len(parts) < 6 {
			return true
		}
		r := BillRecord{
			Date:           parts[0],
			MainNetInflow:  util.ParseFloat(parts[1]),
			SmallNetInflow: util.ParseFloat(parts[2]),
			MidNetInflow:   util.ParseFloat(parts[3]),
			LargeNetInflow: util.ParseFloat(parts[4]),
			HugeNetInflow:  util.ParseFloat(parts[5]),
		}
		if len(parts) > 6 {
			r.MainNetInflowPct = util.ParseFloat(parts[6])
		}
		if len(parts) > 7 {
			r.SmallInflowPct = util.ParseFloat(parts[7])
		}
		if len(parts) > 8 {
			r.MidInflowPct = util.ParseFloat(parts[8])
		}
		if len(parts) > 9 {
			r.LargeInflowPct = util.ParseFloat(parts[9])
		}
		if len(parts) > 10 {
			r.HugeInflowPct = util.ParseFloat(parts[10])
		}
		if len(parts) > 11 {
			r.ClosePrice = util.ParseFloat(parts[11])
		}
		if len(parts) > 12 {
			r.ChangeRate = util.ParseFloat(parts[12])
		}
		records = append(records, r)
		return true
	})

	return &KlineResult{Code: code, Name: name}, records, nil
}

func GetTodayBill(secid string) (*KlineResult, []TodayBillRecord, error) {
	params := url.Values{
		"lmt":     {"0"},
		"klt":     {"1"},
		"secid":   {secid},
		"fields1": {"f1,f2,f3,f7"},
		"fields2": {"f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63"},
	}
	u := "http://push2.eastmoney.com/api/qt/stock/fflow/kline/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, nil, err
	}

	result := gjson.ParseBytes(body)
	name := result.Get("data.name").String()
	code := secid
	if idx := strings.LastIndex(secid, "."); idx >= 0 {
		code = secid[idx+1:]
	}

	klines := result.Get("data.klines")
	if !klines.Exists() || !klines.IsArray() {
		return &KlineResult{Code: code, Name: name}, nil, nil
	}

	var records []TodayBillRecord
	klines.ForEach(func(_, v gjson.Result) bool {
		parts := strings.Split(v.String(), ",")
		if len(parts) < 6 {
			return true
		}
		records = append(records, TodayBillRecord{
			Time:           parts[0],
			MainNetInflow:  util.ParseFloat(parts[1]),
			SmallNetInflow: util.ParseFloat(parts[2]),
			MidNetInflow:   util.ParseFloat(parts[3]),
			LargeNetInflow: util.ParseFloat(parts[4]),
			HugeNetInflow:  util.ParseFloat(parts[5]),
		})
		return true
	})

	return &KlineResult{Code: code, Name: name}, records, nil
}

func GetBaseInfo(secid string) (BaseInfo, error) {
	fields := make([]string, 0, len(BaseInfoFields))
	for k := range BaseInfoFields {
		fields = append(fields, k)
	}
	params := url.Values{
		"ut":     {"fa5fd1943c7b386f172d6893dbfba10b"},
		"invt":   {"2"},
		"fltt":   {"2"},
		"fields": {strings.Join(fields, ",")},
		"secid":  {secid},
	}
	u := "http://push2.eastmoney.com/api/qt/stock/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	data := result.Get("data")
	if !data.Exists() {
		return BaseInfo{}, nil
	}

	info := make(BaseInfo)
	for apiKey, cnName := range BaseInfoFields {
		v := data.Get(apiKey).Value()
		info[cnName] = v
	}
	return info, nil
}

func GetLatestQuote(secids []string) ([]RealtimeQuote, error) {
	fields := strings.Join(QuoteFieldKeys, ",")
	params := url.Values{
		"OSVersion":     {"14.3"},
		"appVersion":    {"6.3.8"},
		"fields":        {fields},
		"fltt":          {"2"},
		"plat":          {"Iphone"},
		"product":       {"EFund"},
		"secids":        {strings.Join(secids, ",")},
		"serverVersion": {"6.3.6"},
		"version":       {"6.3.8"},
	}
	u := "https://push2.eastmoney.com/api/qt/ulist.np/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	diff := result.Get("data.diff")
	if !diff.Exists() || !diff.IsArray() {
		return nil, nil
	}

	var quotes []RealtimeQuote
	diff.ForEach(func(_, item gjson.Result) bool {
		mktNum := item.Get("f13").String()
		code := item.Get("f12").String()
		quoteID := mktNum + "." + code
		marketType := MarketNumberDict[mktNum]

		q := RealtimeQuote{
			Code:       code,
			Name:       item.Get("f14").String(),
			Price:      util.ParseFloat(item.Get("f2").String()),
			ChangeRate: util.ParseFloat(item.Get("f3").String()),
			ChangeAmt:  util.ParseFloat(item.Get("f4").String()),
			Volume:     util.ParseInt(item.Get("f5").String()),
			Amount:     util.ParseFloat(item.Get("f6").String()),
			High:       util.ParseFloat(item.Get("f15").String()),
			Low:        util.ParseFloat(item.Get("f16").String()),
			Open:       util.ParseFloat(item.Get("f17").String()),
			PreClose:   util.ParseFloat(item.Get("f18").String()),
			TurnoverRate: util.ParseFloat(item.Get("f8").String()),
			PERatio:    util.ParseFloat(item.Get("f9").String()),
			MarketCap:  util.ParseFloat(item.Get("f20").String()),
			CircMarketCap: util.ParseFloat(item.Get("f21").String()),
			QuoteID:    quoteID,
			MarketType: marketType,
			UpdateTimestamp: item.Get("f124").Int(),
			LatestTradeDate: item.Get("f297").String(),
		}
		ts := q.UpdateTimestamp
		if ts > 0 {
			q.UpdateDate = time.Unix(ts, 0).Format("2006-01-02 15:04:05")
		}
		if q.LatestTradeDate != "" && len(q.LatestTradeDate) == 8 {
			q.LatestTradeDate = q.LatestTradeDate[:4] + "-" + q.LatestTradeDate[4:6] + "-" + q.LatestTradeDate[6:8]
		}
		quotes = append(quotes, q)
		return true
	})
	return quotes, nil
}

func GetDealDetail(secid string, maxCount int) ([]DealRecord, error) {
	baseInfo, err := GetBaseInfo(secid)
	if err != nil {
		return nil, err
	}
	code, _ := baseInfo["代码"].(string)
	name, _ := baseInfo["名称"].(string)
	if code == "" {
		return nil, nil
	}

	params := url.Values{
		"secid":   {secid},
		"fields1": {"f1,f2,f3,f4,f5"},
		"fields2": {"f51,f52,f53,f54,f55"},
		"pos":     {fmt.Sprintf("-%d", maxCount)},
	}
	u := "https://push2.eastmoney.com/api/qt/stock/details/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	details := result.Get("data.details")
	prePrice := result.Get("data.prePrice").Float()

	if !details.Exists() || !details.IsArray() {
		return nil, nil
	}

	var records []DealRecord
	details.ForEach(func(_, v gjson.Result) bool {
		parts := strings.Split(v.String(), ",")
		if len(parts) < 4 {
			return true
		}
		records = append(records, DealRecord{
			Name:     name,
			Code:     code,
			Time:     parts[0],
			PreClose: prePrice,
			Price:    util.ParseFloat(parts[1]),
			Volume:   util.ParseInt(parts[2]),
			OrderNum: util.ParseInt(parts[3]),
		})
		return true
	})
	return records, nil
}

func GetLatestNDaysQuote(secid string, ndays int) (*KlineResult, []NDaysKlineBar, error) {
	fields2 := strings.Join(KlineNDaysFields, ",")
	params := url.Values{
		"fields1": {"f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13"},
		"fields2": {fields2},
		"ndays":   {strconv.Itoa(ndays)},
		"iscr":    {"0"},
		"iscca":   {"0"},
		"secid":   {secid},
	}
	u := "http://push2his.eastmoney.com/api/qt/stock/trends2/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, nil, err
	}

	result := gjson.ParseBytes(body)
	data := result.Get("data")
	name := data.Get("name").String()
	code := secid
	if idx := strings.LastIndex(secid, "."); idx >= 0 {
		code = secid[idx+1:]
	}

	trends := data.Get("trends")
	if !trends.Exists() || !trends.IsArray() {
		return &KlineResult{Code: code, Name: name}, nil, nil
	}

	var bars []NDaysKlineBar
	trends.ForEach(func(_, v gjson.Result) bool {
		parts := strings.Split(v.String(), ",")
		if len(parts) < 7 {
			return true
		}
		bars = append(bars, NDaysKlineBar{
			Date:   parts[0],
			Open:   util.ParseFloat(parts[1]),
			Close:  util.ParseFloat(parts[2]),
			High:   util.ParseFloat(parts[3]),
			Low:    util.ParseFloat(parts[4]),
			Volume: util.ParseInt(parts[5]),
			Amount: util.ParseFloat(parts[6]),
		})
		return true
	})

	return &KlineResult{Code: code, Name: name}, bars, nil
}

func ResolveSecid(code string, quoteIDMode bool) (string, error) {
	if quoteIDMode {
		return code, nil
	}
	return quote.GetQuoteID(code)
}
