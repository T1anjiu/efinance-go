package bond

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/T1anjiu/efinance-go/common"
	"github.com/T1anjiu/efinance-go/internal/client"
	"github.com/T1anjiu/efinance-go/internal/util"
	"github.com/T1anjiu/efinance-go/quote"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
)

func GetBaseInfo(bondCodes []string) ([]BondBaseInfo, error) {
	var results []BondBaseInfo
	var mu sync.Mutex

	eg, _ := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, common.MaxConnections)

	for _, code := range bondCodes {
		c := code
		eg.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()

			info, err := getBaseInfoSingle(c)
			if err != nil {
				return nil
			}
			mu.Lock()
			results = append(results, *info)
			mu.Unlock()
			return nil
		})
	}
	eg.Wait()
	return results, nil
}

func getBaseInfoSingle(bondCode string) (*BondBaseInfo, error) {
	params := url.Values{
		"reportName": {"RPT_BOND_CB_LIST"},
		"columns":    {"ALL"},
		"source":     {"WEB"},
		"client":     {"WEB"},
		"filter":     {fmt.Sprintf(`(SECURITY_CODE="%s")`, bondCode)},
	}
	u := "http://datacenter-web.eastmoney.com/api/data/v1/get"
	body, err := client.Get(u, params, nil)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	data := result.Get("result.data")
	if !data.Exists() || !data.IsArray() || len(data.Array()) == 0 {
		return &BondBaseInfo{}, nil
	}

	item := data.Array()[0]
	return &BondBaseInfo{
		BondCode:      item.Get("SECURITY_CODE").String(),
		BondName:      item.Get("SECURITY_NAME_ABBR").String(),
		StockCode:     item.Get("CONVERT_STOCK_CODE").String(),
		StockName:     item.Get("SECURITY_SHORT_NAME").String(),
		Rating:        item.Get("RATING").String(),
		SubscribeDate: item.Get("PUBLIC_START_DATE").String(),
		IssueScale:    item.Get("ACTUAL_ISSUE_SCALE").Float(),
		WinRate:       item.Get("ONLINE_GENERAL_LWR").Float(),
		ListingDate:   item.Get("LISTING_DATE").String(),
		ExpireDate:    item.Get("EXPIRE_DATE").String(),
		Term:          item.Get("BOND_EXPIRE").Float(),
		InterestDesc:  item.Get("INTEREST_RATE_EXPLAIN").String(),
	}, nil
}

func GetAllBaseInfo() ([]BondBaseInfo, error) {
	var all []BondBaseInfo
	page := 1

	for {
		params := url.Values{
			"sortColumns": {"PUBLIC_START_DATE"},
			"sortTypes":   {"-1"},
			"pageSize":    {"500"},
			"pageNumber":  {strconv.Itoa(page)},
			"reportName":  {"RPT_BOND_CB_LIST"},
			"columns":     {"ALL"},
			"source":      {"WEB"},
			"client":      {"WEB"},
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
			all = append(all, BondBaseInfo{
				BondCode:      item.Get("SECURITY_CODE").String(),
				BondName:      item.Get("SECURITY_NAME_ABBR").String(),
				StockCode:     item.Get("CONVERT_STOCK_CODE").String(),
				StockName:     item.Get("SECURITY_SHORT_NAME").String(),
				Rating:        item.Get("RATING").String(),
				SubscribeDate: item.Get("PUBLIC_START_DATE").String(),
				IssueScale:    item.Get("ACTUAL_ISSUE_SCALE").Float(),
				WinRate:       item.Get("ONLINE_GENERAL_LWR").Float(),
				ListingDate:   item.Get("LISTING_DATE").String(),
				ExpireDate:    item.Get("EXPIRE_DATE").String(),
				Term:          item.Get("BOND_EXPIRE").Float(),
				InterestDesc:  item.Get("INTEREST_RATE_EXPLAIN").String(),
			})
			return true
		})

		page++
	}
	return all, nil
}

func GetRealtimeQuotes() ([]BondQuote, error) {
	fs := common.FSDict["bond"]
	rq, err := common.GetRealtimeQuotesByFS(fs)
	if err != nil {
		return nil, err
	}
	quotes := make([]BondQuote, len(rq))
	for i, q := range rq {
		quotes[i] = BondQuote{RealtimeQuote: q}
	}
	return quotes, nil
}

func GetQuoteHistory(bondCodes []string, opts ...common.KlineOptions) (map[string]*BondKline, error) {
	var kopts common.KlineOptions
	if len(opts) > 0 {
		kopts = opts[0]
	} else {
		kopts = common.DefaultKlineOptions()
	}

	results := make(map[string]*BondKline)
	var mu sync.Mutex

	eg, _ := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, common.MaxConnections)

	for _, code := range bondCodes {
		c := code
		eg.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()

			secid, err := quote.GetQuoteID(c)
			if err != nil {
				return nil
			}
			r, err := common.GetQuoteHistorySingle(secid, kopts)
			if err != nil {
				return nil
			}
			mu.Lock()
			results[c] = &BondKline{Code: r.Code, Name: r.Name, Bars: r.Bars}
			mu.Unlock()
			return nil
		})
	}
	eg.Wait()
	return results, nil
}

func GetHistoryBill(bondCode string) ([]common.BillRecord, string, string, error) {
	secid, err := quote.GetQuoteID(bondCode)
	if err != nil {
		return nil, "", "", err
	}
	info, records, err := common.GetHistoryBill(secid)
	if err != nil {
		return nil, "", "", err
	}
	return records, info.Code, info.Name, nil
}

func GetTodayBill(bondCode string) ([]common.TodayBillRecord, string, string, error) {
	secid, err := quote.GetQuoteID(bondCode)
	if err != nil {
		return nil, "", "", err
	}
	info, records, err := common.GetTodayBill(secid)
	if err != nil {
		return nil, "", "", err
	}
	return records, info.Code, info.Name, nil
}

func GetDealDetail(bondCode string, maxCount int, quoteIDMode bool) ([]common.DealRecord, error) {
	var secid string
	if quoteIDMode {
		secid = bondCode
	} else {
		var err error
		secid, err = quote.GetQuoteID(bondCode)
		if err != nil {
			return nil, err
		}
	}
	return common.GetDealDetail(secid, maxCount)
}

func init() {
	_ = util.ParseFloat
	_ = strings.Join
	_ = context.Background
}
