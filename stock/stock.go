package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/T1anjiu/efinance-go/common"
	"github.com/T1anjiu/efinance-go/internal/client"
	"github.com/T1anjiu/efinance-go/internal/util"
	"github.com/T1anjiu/efinance-go/quote"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
)

type StockKline struct {
	Code string           `json:"code"`
	Name string           `json:"name"`
	Bars []common.KlineBar `json:"bars"`
}

type StockQuote struct {
	common.RealtimeQuote
}

type StockBill struct {
	Code   string               `json:"code"`
	Name   string               `json:"name"`
	Record common.BillRecord     `json:"record"`
}

type StockTodayBill struct {
	Code   string                  `json:"code"`
	Name   string                  `json:"name"`
	Record common.TodayBillRecord  `json:"record"`
}

type StockBaseInfo struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	PERatio     float64 `json:"pe_ratio"`
	PBRatio     float64 `json:"pb_ratio"`
	Industry    string  `json:"industry"`
	TotalMV     float64 `json:"total_mv"`
	CircMV      float64 `json:"circ_mv"`
	BoardCode   string  `json:"board_code"`
	ROE         float64 `json:"roe"`
	NetMargin   float64 `json:"net_margin"`
	NetProfit   float64 `json:"net_profit"`
	GrossMargin float64 `json:"gross_margin"`
}

type HolderInfo struct {
	StockCode   string  `json:"stock_code"`
	Date        string  `json:"date"`
	HolderCode  string  `json:"holder_code"`
	HolderName  string  `json:"holder_name"`
	HoldCount   float64 `json:"hold_count"`
	HoldRatio   float64 `json:"hold_ratio"`
	Change      float64 `json:"change"`
	ChangeRatio float64 `json:"change_ratio"`
}

type BillboardRecord struct {
	StockCode       string  `json:"stock_code"`
	StockName       string  `json:"stock_name"`
	TradeDate       string  `json:"trade_date"`
	Explain         string  `json:"explain"`
	ClosePrice      float64 `json:"close_price"`
	ChangeRate      float64 `json:"change_rate"`
	TurnoverRate    float64 `json:"turnover_rate"`
	NetBuyAmt       float64 `json:"net_buy_amt"`
	BuyAmt          float64 `json:"buy_amt"`
	SellAmt         float64 `json:"sell_amt"`
	DealAmt         float64 `json:"deal_amt"`
	TotalAmt        float64 `json:"total_amt"`
	NetBuyRatio     float64 `json:"net_buy_ratio"`
	DealRatio       float64 `json:"deal_ratio"`
	FreeMarketCap   float64 `json:"free_market_cap"`
	Reason          string  `json:"reason"`
}

type PerformanceRecord = map[string]interface{}

type HolderNumberRecord = map[string]interface{}

type IPORecord = map[string]interface{}

type IndexMember struct {
	IndexCode string  `json:"index_code"`
	IndexName string  `json:"index_name"`
	StockCode string  `json:"stock_code"`
	StockName string  `json:"stock_name"`
	Weight    float64 `json:"weight"`
}

type QuoteSnapshot struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Time       string  `json:"time"`
	ChangeAmt  float64 `json:"change_amt"`
	ChangeRate float64 `json:"change_rate"`
	Price      float64 `json:"price"`
	PreClose   float64 `json:"pre_close"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	AvgPrice   float64 `json:"avg_price"`
	TopPrice   float64 `json:"top_price"`
	BottomPrice float64 `json:"bottom_price"`
	Turnover   float64 `json:"turnover"`
	Volume     int64   `json:"volume"`
	Amount     float64 `json:"amount"`
	Sell1      float64 `json:"sell1"`
	Sell2      float64 `json:"sell2"`
	Sell3      float64 `json:"sell3"`
	Sell4      float64 `json:"sell4"`
	Sell5      float64 `json:"sell5"`
	Buy1       float64 `json:"buy1"`
	Buy2       float64 `json:"buy2"`
	Buy3       float64 `json:"buy3"`
	Buy4       float64 `json:"buy4"`
	Buy5       float64 `json:"buy5"`
	Sell1Count float64 `json:"sell1_count"`
	Sell2Count float64 `json:"sell2_count"`
	Sell3Count float64 `json:"sell3_count"`
	Sell4Count float64 `json:"sell4_count"`
	Sell5Count float64 `json:"sell5_count"`
	Buy1Count  float64 `json:"buy1_count"`
	Buy2Count  float64 `json:"buy2_count"`
	Buy3Count  float64 `json:"buy3_count"`
	Buy4Count  float64 `json:"buy4_count"`
	Buy5Count  float64 `json:"buy5_count"`
}

type BoardInfo struct {
	StockName   string  `json:"stock_name"`
	StockCode   string  `json:"stock_code"`
	BoardCode   string  `json:"board_code"`
	BoardName   string  `json:"board_name"`
	BoardChangeRate float64 `json:"board_change_rate"`
}

type StockOption func(*stockOpts)

type stockOpts struct {
	marketType   *MarketType
	suppressError bool
	useIDCache   bool
	quoteIDMode  bool
}

func WithMarketType(mt MarketType) StockOption {
	return func(o *stockOpts) { o.marketType = &mt }
}

func WithSuppressError(v bool) StockOption {
	return func(o *stockOpts) { o.suppressError = v }
}

func WithIDCache(v bool) StockOption {
	return func(o *stockOpts) { o.useIDCache = v }
}

func WithQuoteIDMode(v bool) StockOption {
	return func(o *stockOpts) { o.quoteIDMode = v }
}

func GetQuoteHistory(codes []string, opts ...StockOption) (map[string]*StockKline, error) {
	o := &stockOpts{useIDCache: true}
	for _, opt := range opts {
		opt(o)
	}

	results := make(map[string]*StockKline)
	var mu sync.Mutex

	eg, _ := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, common.MaxConnections)

	for _, code := range codes {
		c := code
		eg.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()

			var secid string
			var err error
			if o.quoteIDMode {
				secid = c
			} else {
				qopts := []quote.Option{quote.WithCache(o.useIDCache), quote.WithSuppressError(o.suppressError)}
				if o.marketType != nil {
					qopts = append(qopts, quote.WithMarketType(o.marketType.String()))
				}
				secid, err = quote.GetQuoteID(c, qopts...)
				if err != nil {
					if o.suppressError {
						return nil
					}
					return err
				}
			}

			kOpts := common.DefaultKlineOptions()
			r, err := common.GetQuoteHistorySingle(secid, kOpts)
			if err != nil {
				return nil
			}
			mu.Lock()
			results[c] = &StockKline{Code: r.Code, Name: r.Name, Bars: r.Bars}
			mu.Unlock()
			return nil
		})
	}

	_ = context.Background
	if err := eg.Wait(); err != nil {
		return results, err
	}
	return results, nil
}

func GetRealtimeQuotes(markets ...string) ([]StockQuote, error) {
	fsList := make([]string, 0)
	if len(markets) == 0 {
		fsList = append(fsList, common.FSDict["stock"])
	} else {
		for _, m := range markets {
			fs, ok := common.FSDict[m]
			if !ok {
				return nil, fmt.Errorf("指定的行情参数 %q 不正确", m)
			}
			fsList = append(fsList, fs)
		}
	}
	fsStr := strings.Join(fsList, ",")

	rq, err := common.GetRealtimeQuotesByFS(fsStr)
	if err != nil {
		return nil, err
	}

	quotes := make([]StockQuote, len(rq))
	for i, q := range rq {
		quotes[i] = StockQuote{RealtimeQuote: q}
	}
	return quotes, nil
}

func GetHistoryBill(code string) ([]common.BillRecord, string, string, error) {
	secid, err := quote.GetQuoteID(code)
	if err != nil {
		return nil, "", "", err
	}
	info, records, err := common.GetHistoryBill(secid)
	if err != nil {
		return nil, "", "", err
	}
	return records, info.Code, info.Name, nil
}

func GetTodayBill(code string) ([]common.TodayBillRecord, string, string, error) {
	secid, err := quote.GetQuoteID(code)
	if err != nil {
		return nil, "", "", err
	}
	info, records, err := common.GetTodayBill(secid)
	if err != nil {
		return nil, "", "", err
	}
	return records, info.Code, info.Name, nil
}

func GetLatestQuote(codes []string, quoteIDMode bool) ([]StockQuote, error) {
	var secids []string
	if quoteIDMode {
		secids = codes
	} else {
		for _, c := range codes {
			secid, err := quote.GetQuoteID(c)
			if err != nil {
				continue
			}
			secids = append(secids, secid)
		}
	}
	rq, err := common.GetLatestQuote(secids)
	if err != nil {
		return nil, err
	}
	quotes := make([]StockQuote, len(rq))
	for i, q := range rq {
		quotes[i] = StockQuote{RealtimeQuote: q}
	}
	return quotes, nil
}

func GetBaseInfo(codes []string) ([]StockBaseInfo, error) {
	var results []StockBaseInfo
	var mu sync.Mutex

	eg, _ := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, common.MaxConnections)

	for _, code := range codes {
		c := code
		eg.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()

			secid, err := quote.GetQuoteID(c)
			if err != nil {
				return nil
			}
			info, err := common.GetBaseInfo(secid)
			if err != nil {
				return nil
			}

			si := StockBaseInfo{}
			if v, ok := info["股票代码"]; ok {
				si.Code, _ = v.(string)
			}
			if v, ok := info["股票名称"]; ok {
				si.Name, _ = v.(string)
			}
			if v, ok := info["市盈率(动)"]; ok {
				si.PERatio = toFloat(v)
			}
			if v, ok := info["市净率"]; ok {
				si.PBRatio = toFloat(v)
			}
			if v, ok := info["所处行业"]; ok {
				si.Industry, _ = v.(string)
			}
			if v, ok := info["总市值"]; ok {
				si.TotalMV = toFloat(v)
			}
			if v, ok := info["流通市值"]; ok {
				si.CircMV = toFloat(v)
			}
			if v, ok := info["板块编号"]; ok {
				si.BoardCode, _ = v.(string)
			}
			if v, ok := info["ROE"]; ok {
				si.ROE = toFloat(v)
			}
			if v, ok := info["净利率"]; ok {
				si.NetMargin = toFloat(v)
			}
			if v, ok := info["净利润"]; ok {
				si.NetProfit = toFloat(v)
			}
			if v, ok := info["毛利率"]; ok {
				si.GrossMargin = toFloat(v)
			}

			mu.Lock()
			results = append(results, si)
			mu.Unlock()
			return nil
		})
	}

	_ = context.Background
	eg.Wait()
	return results, nil
}

func GetTop10Holders(code string, top int) ([]HolderInfo, error) {
	if top <= 0 {
		top = 4
	}
	secid, err := quote.GetQuoteID(code)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(secid, ".", 2)
	stockCode := parts[1]

	var fc string
	marketNum := parts[0]
	if marketNum == "0" {
		fc = stockCode + "02"
	} else {
		fc = stockCode + "01"
	}

	dates, err := getPublicDates(fc)
	if err != nil {
		return nil, err
	}

	var holders []HolderInfo

	limit := top
	if limit > len(dates) {
		limit = len(dates)
	}

	for _, date := range dates[:limit] {
		data, _ := json.Marshal(map[string]string{"fc": fc, "BaoGaoQi": date})
		u := "https://emh5.eastmoney.com/api/GuBenGuDong/GetShiDaLiuTongGuDong"
		body, err := client.PostJSON(u, strings.NewReader(string(data)), nil)
		if err != nil {
			continue
		}

		result := gjson.ParseBytes(body)
		items := result.Get("data.ShiDaLiuTongGuDongList")
		if !items.Exists() || !items.IsArray() {
			continue
		}

		items.ForEach(func(_, item gjson.Result) bool {
			h := HolderInfo{
				StockCode:  stockCode,
				Date:       date,
				HolderCode: item.Get("GuDongDaiMa").String(),
				HolderName: item.Get("GuDongMingCheng").String(),
				HoldCount:  item.Get("ChiGuShu").Float(),
				HoldRatio:  item.Get("ChiGuBiLi").Float(),
				Change:     item.Get("ZengJian").Float(),
				ChangeRatio: item.Get("BianDongBiLi").Float(),
			}
			holders = append(holders, h)
			return true
		})
	}
	return holders, nil
}

func getPublicDates(fc string) ([]string, error) {
	data, _ := json.Marshal(map[string]string{"fc": fc})
	u := "https://emh5.eastmoney.com/api/GuBenGuDong/GetFirstRequest2Data"
	body, err := client.PostJSON(u, strings.NewReader(string(data)), nil)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	var dates []string
	result.Get("data.BaoGaoQi").ForEach(func(_, v gjson.Result) bool {
		dates = append(dates, v.String())
		return true
	})
	return dates, nil
}

func GetAllReportDates() ([]string, error) {
	fieldKeys := []string{"REPORT_DATE", "DATATYPE"}
	params := url.Values{
		"type": {"RPT_LICO_FN_CPD_BBBQ"},
		"sty":  {strings.Join(fieldKeys, ",")},
		"p":    {"1"},
		"ps":   {"2000"},
	}
	u := "https://datacenter.eastmoney.com/securities/api/data/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	var dates []string
	result.Get("result.data").ForEach(func(_, item gjson.Result) bool {
		d := item.Get("REPORT_DATE").String()
		if idx := strings.Index(d, " "); idx > 0 {
			d = d[:idx]
		}
		dates = append(dates, d)
		return true
	})
	return dates, nil
}

func GetAllCompanyPerformance(date string) ([]PerformanceRecord, error) {
	dates, err := GetAllReportDates()
	if err != nil {
		return nil, err
	}
	if date == "" && len(dates) > 0 {
		date = dates[0]
	}

	filter := fmt.Sprintf("(REPORTDATE='%s')", date)

	getByPage := func(pn, pz int) (gjson.Result, error) {
		params := url.Values{
			"st":    {"NOTICE_DATE,SECURITY_CODE"},
			"sr":    {"-1,-1"},
			"ps":    {strconv.Itoa(pz)},
			"p":     {strconv.Itoa(pn)},
			"type":  {"RPT_LICO_FN_CPD"},
			"sty":   {"ALL"},
			"token": {"894050c76af8597a853f5b408b759f5d"},
			"filter": {fmt.Sprintf(`(SECURITY_TYPE_CODE in ("058001001","058001008"))%s`, filter)},
		}
		u := "http://datacenter-web.eastmoney.com/api/data/get"
		body, err := client.Get(u, params, nil)
		if err != nil {
			return gjson.Result{}, err
		}
		return gjson.ParseBytes(body), nil
	}

	firstResp, err := getByPage(1, 500)
	if err != nil {
		return nil, err
	}

	total := int(firstResp.Get("result.count").Int())
	if total == 0 {
		return nil, nil
	}

	data := firstResp.Get("result.data")
	pz := len(data.Array())
	if pz == 0 {
		return nil, nil
	}
	pages := total / pz
	if total%pz != 0 {
		pages++
	}

	var allRecords []PerformanceRecord
	firstResp.Get("result.data").ForEach(func(_, item gjson.Result) bool {
		r := make(PerformanceRecord)
		for k, v := range CompanyPerformanceFields {
			r[v] = item.Get(k).Value()
		}
		allRecords = append(allRecords, r)
		return true
	})

	if pages > 1 {
		eg, _ := errgroup.WithContext(context.Background())
		sem := make(chan struct{}, common.MaxConnections)
		var mu sync.Mutex

		for i := 2; i <= pages; i++ {
			pn := i
			eg.Go(func() error {
				sem <- struct{}{}
				defer func() { <-sem }()
				resp, err := getByPage(pn, pz)
				if err != nil {
					return err
				}
				resp.Get("result.data").ForEach(func(_, item gjson.Result) bool {
					r := make(PerformanceRecord)
					for k, v := range CompanyPerformanceFields {
						r[v] = item.Get(k).Value()
					}
					mu.Lock()
					allRecords = append(allRecords, r)
					mu.Unlock()
					return true
				})
				return nil
			})
		}
		eg.Wait()
	}

	return allRecords, nil
}

func GetLatestHolderNumber(date string) ([]HolderNumberRecord, error) {
	reportName := "RPT_HOLDERNUMLATEST"
	var filter string

	if date != "" {
		t, err := time.Parse("2006-01-02", date)
		if err == nil {
			year, month := t.Year(), int(t.Month())
			if month%3 != 0 {
				month -= month % 3
			}
			if month < 3 {
				year--
				month = 12
			}
			dim := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.Local)
			date = dim.Format("2006-01-02")
		}
		filter = fmt.Sprintf("(END_DATE='%s')", date)
		reportName = "RPT_HOLDERNUM_DET"
	}

	getByPage := func(pn, pz int) (gjson.Result, error) {
		params := []struct{ k, v string }{
			{"sortColumns", "HOLD_NOTICE_DATE,SECURITY_CODE"},
			{"sortTypes", "-1,-1"},
			{"pageSize", strconv.Itoa(pz)},
			{"pageNumber", strconv.Itoa(pn)},
			{"columns", "SECURITY_CODE,SECURITY_NAME_ABBR,END_DATE,INTERVAL_CHRATE,AVG_MARKET_CAP,AVG_HOLD_NUM,TOTAL_MARKET_CAP,TOTAL_A_SHARES,HOLD_NOTICE_DATE,HOLDER_NUM,PRE_HOLDER_NUM,HOLDER_NUM_CHANGE,HOLDER_NUM_RATIO,END_DATE,PRE_END_DATE"},
			{"quoteColumns", "f2,f3"},
			{"source", "WEB"},
			{"client", "WEB"},
			{"reportName", reportName},
		}
		if filter != "" {
			params = append(params, struct{ k, v string }{"filter", filter})
		}

		v := url.Values{}
		for _, p := range params {
			v.Set(p.k, p.v)
		}
		u := "http://datacenter-web.eastmoney.com/api/data/v1/get"
		body, err := client.Get(u, v, nil)
		if err != nil {
			return gjson.Result{}, err
		}
		return gjson.ParseBytes(body), nil
	}

	firstResp, err := getByPage(1, 500)
	if err != nil {
		return nil, err
	}

	total := int(firstResp.Get("result.count").Int())
	if total == 0 {
		return nil, nil
	}

	var allRecords []HolderNumberRecord
	firstResp.Get("result.data").ForEach(func(_, item gjson.Result) bool {
		r := make(HolderNumberRecord)
		for k, v := range HolderNumberFields {
			r[v] = item.Get(k).Value()
		}
		allRecords = append(allRecords, r)
		return true
	})

	return allRecords, nil
}

func GetDailyBillboard(startDate, endDate string) ([]BillboardRecord, error) {
	today := time.Now().Format("2006-01-02")
	if startDate == "" {
		startDate = today
	}
	if endDate == "" {
		endDate = today
	}

	var records []BillboardRecord
	page := 1

	for {
		params := url.Values{
			"sortColumns": {"TRADE_DATE,SECURITY_CODE"},
			"sortTypes":   {"-1,1"},
			"pageSize":    {"500"},
			"pageNumber":  {strconv.Itoa(page)},
			"reportName":  {"RPT_DAILYBILLBOARD_DETAILS"},
			"columns":     {"ALL"},
			"source":      {"WEB"},
			"client":      {"WEB"},
			"filter":      {fmt.Sprintf("(TRADE_DATE<='%s')(TRADE_DATE>='%s')", endDate, startDate)},
		}
		u := "http://datacenter-web.eastmoney.com/api/data/v1/get"
		body, err := client.Get(u, params, nil)
		if err != nil {
			break
		}

		result := gjson.ParseBytes(body)
		data := result.Get("result.data")
		if !data.Exists() || !data.IsArray() || len(data.Array()) == 0 {
			break
		}

		data.ForEach(func(_, item gjson.Result) bool {
			r := BillboardRecord{
				StockCode:     item.Get("SECURITY_CODE").String(),
				StockName:     item.Get("SECURITY_NAME_ABBR").String(),
				TradeDate:     item.Get("TRADE_DATE").String(),
				Explain:       item.Get("EXPLAIN").String(),
				ClosePrice:    item.Get("CLOSE_PRICE").Float(),
				ChangeRate:    item.Get("CHANGE_RATE").Float(),
				TurnoverRate:  item.Get("TURNOVERRATE").Float(),
				NetBuyAmt:     item.Get("BILLBOARD_NET_AMT").Float(),
				BuyAmt:        item.Get("BILLBOARD_BUY_AMT").Float(),
				SellAmt:       item.Get("BILLBOARD_SELL_AMT").Float(),
				DealAmt:       item.Get("BILLBOARD_DEAL_AMT").Float(),
				TotalAmt:      item.Get("ACCUM_AMOUNT").Float(),
				NetBuyRatio:   item.Get("DEAL_NET_RATIO").Float(),
				DealRatio:     item.Get("DEAL_AMOUNT_RATIO").Float(),
				FreeMarketCap: item.Get("FREE_MARKET_CAP").Float(),
				Reason:        item.Get("EXPLANATION").String(),
			}
			if idx := strings.Index(r.TradeDate, " "); idx > 0 {
				r.TradeDate = r.TradeDate[:idx]
			}
			records = append(records, r)
			return true
		})

		page++
	}
	return records, nil
}

func GetIndexMembers(indexCode string) ([]IndexMember, error) {
	qs, err := quote.SearchMulti(indexCode, 10)
	if err != nil || len(qs) == 0 {
		return nil, fmt.Errorf("未找到指数 %q", indexCode)
	}

	for _, q := range qs {
		if q.SecurityTypeName != "指数" {
			continue
		}
		params := url.Values{
			"IndexCode":    {q.Code},
			"pageIndex":    {"1"},
			"pageSize":     {"10000"},
			"deviceid":     {"1234567890"},
			"version":      {"6.9.9"},
			"product":      {"EFund"},
			"plat":         {"Iphone"},
			"ServerVersion": {"6.9.9"},
		}
		u := "https://fundztapi.eastmoney.com/FundSpecialApiNew/FundSpecialZSB30ZSCFG"
		body, err := client.Get(u, params, nil)
		if err != nil {
			continue
		}

		result := gjson.ParseBytes(body)
		items := result.Get("Datas")
		if !items.Exists() || !items.IsArray() || len(items.Array()) == 0 {
			continue
		}

		var members []IndexMember
		items.ForEach(func(_, item gjson.Result) bool {
			members = append(members, IndexMember{
				IndexCode: item.Get("IndexCode").String(),
				IndexName: item.Get("IndexName").String(),
				StockCode: item.Get("StockCode").String(),
				StockName: item.Get("StockName").String(),
				Weight:    item.Get("MARKETCAPPCT").Float(),
			})
			return true
		})
		return members, nil
	}
	return nil, fmt.Errorf("未找到可用的指数成分股数据")
}

func GetLatestIPOInfo() ([]IPORecord, error) {
	var records []IPORecord
	page := 1

	for {
		params := url.Values{
			"st":    {"UPDATE_DATE,SECURITY_CODE"},
			"sr":    {"-1,-1"},
			"ps":    {"500"},
			"p":     {strconv.Itoa(page)},
			"type":  {"RPT_REGISTERED_INFO"},
			"sty":   {"ORG_CODE,ISSUER_NAME,CHECK_STATUS,CHECK_STATUS_CODE,REG_ADDRESS,CSRC_INDUSTRY,RECOMMEND_ORG,LAW_FIRM,ACCOUNT_FIRM,UPDATE_DATE,ACCEPT_DATE,TOLIST_MARKET,SECURITY_CODE"},
			"token": {"894050c76af8597a853f5b408b759f5d"},
			"client": {"WEB"},
		}
		u := "http://datacenter-web.eastmoney.com/api/data/get"
		body, err := client.Get(u, params, nil)
		if err != nil {
			break
		}

		result := gjson.ParseBytes(body)
		data := result.Get("result.data")
		if !data.Exists() || !data.IsArray() || len(data.Array()) == 0 {
			break
		}

		data.ForEach(func(_, item gjson.Result) bool {
			r := make(IPORecord)
			for k, v := range IPOFields {
				r[v] = item.Get(k).Value()
			}
			records = append(records, r)
			return true
		})

		page++
	}
	return records, nil
}

func GetQuoteSnapshot(code string) (*QuoteSnapshot, error) {
	params := url.Values{
		"id":       {code},
		"callback": {"jQuery183026310160411569883_1646052793441"},
	}
	u := "https://hsmarketwg.eastmoney.com/api/SHSZQuoteSnapshot"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, err
	}

	text := string(body)
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start < 0 || end < 0 || end <= start {
		return &QuoteSnapshot{}, nil
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(text[start:end+1]), &raw); err != nil {
		return &QuoteSnapshot{}, nil
	}

	snap := &QuoteSnapshot{}

	if fivequote, ok := raw["fivequote"].(map[string]interface{}); ok {
		if realtime, ok := raw["realtimequote"].(map[string]interface{}); ok {
			for k, v := range fivequote {
				raw[k] = v
			}
			for k, v := range realtime {
				raw[k] = v
			}
		}
	}

	snap.Code = toString(raw["code"])
	snap.Name = toString(raw["name"])
	snap.Time = toString(raw["time"])
	snap.ChangeAmt = toFloat(raw["zd"])
	snap.ChangeRate = toFloat(raw["zdf"])
	snap.Price = toFloat(raw["currentPrice"])
	snap.PreClose = toFloat(raw["yesClosePrice"])
	snap.Open = toFloat(raw["openPrice"])
	snap.High = toFloat(raw["high"])
	snap.Low = toFloat(raw["low"])
	snap.AvgPrice = toFloat(raw["avg"])
	snap.TopPrice = toFloat(raw["topprice"])
	snap.BottomPrice = toFloat(raw["bottomprice"])
	snap.Turnover = toFloat(raw["turnover"])
	snap.Volume = int64(toFloat(raw["volume"]))
	snap.Amount = toFloat(raw["amount"])
	snap.Sell1 = toFloat(raw["sale1"])
	snap.Sell2 = toFloat(raw["sale2"])
	snap.Sell3 = toFloat(raw["sale3"])
	snap.Sell4 = toFloat(raw["sale4"])
	snap.Sell5 = toFloat(raw["sale5"])
	snap.Buy1 = toFloat(raw["buy1"])
	snap.Buy2 = toFloat(raw["buy2"])
	snap.Buy3 = toFloat(raw["buy3"])
	snap.Buy4 = toFloat(raw["buy4"])
	snap.Buy5 = toFloat(raw["buy5"])
	snap.Sell1Count = toFloat(raw["sale1_count"])
	snap.Sell2Count = toFloat(raw["sale2_count"])
	snap.Sell3Count = toFloat(raw["sale3_count"])
	snap.Sell4Count = toFloat(raw["sale4_count"])
	snap.Sell5Count = toFloat(raw["sale5_count"])
	snap.Buy1Count = toFloat(raw["buy1_count"])
	snap.Buy2Count = toFloat(raw["buy2_count"])
	snap.Buy3Count = toFloat(raw["buy3_count"])
	snap.Buy4Count = toFloat(raw["buy4_count"])
	snap.Buy5Count = toFloat(raw["buy5_count"])

	return snap, nil
}

func GetDealDetail(code string, maxCount int, quoteIDMode bool) ([]common.DealRecord, error) {
	var secid string
	if quoteIDMode {
		secid = code
	} else {
		var err error
		secid, err = quote.GetQuoteID(code)
		if err != nil {
			return nil, err
		}
	}
	return common.GetDealDetail(secid, maxCount)
}

func GetBelongBoard(code string) ([]BoardInfo, error) {
	q, err := quote.Search(code)
	if err != nil || q == nil {
		return nil, fmt.Errorf("未找到证券 %q", code)
	}

	params := url.Values{
		"forcect": {"1"},
		"spt":     {"3"},
		"fields":  {"f1,f12,f152,f3,f14,f128,f136"},
		"pi":      {"0"},
		"pz":      {"1000"},
		"po":      {"1"},
		"fid":     {"f3"},
		"fid0":    {"f4003"},
		"invt":    {"2"},
		"secid":   {q.QuoteID},
	}
	u := "https://push2.eastmoney.com/api/qt/slist/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	diff := result.Get("data.diff")
	if !diff.Exists() || !diff.IsArray() {
		return nil, nil
	}

	var boards []BoardInfo
	diff.ForEach(func(_, item gjson.Result) bool {
		rate := item.Get("f3").Float()
		boards = append(boards, BoardInfo{
			StockName:     q.Name,
			StockCode:     q.Code,
			BoardCode:     item.Get("f12").String(),
			BoardName:     item.Get("f14").String(),
			BoardChangeRate: rate / 100,
		})
		return true
	})
	return boards, nil
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		s := strings.TrimSpace(val)
		s = strings.TrimSuffix(s, "%")
		return util.ParseFloat(s)
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprintf("%v", v)
	}
}
