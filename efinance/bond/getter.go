package bond

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/T1anjiu/efinance-go/efinance/common"
	"github.com/T1anjiu/efinance-go/efinance/errors"
)

// BondInfo 债券信息
type BondInfo struct {
	Code          string  `json:"code"`          // 债券代码
	Name          string  `json:"name"`          // 债券名称
	StockCode     string  `json:"stock_code"`    // 正股代码
	StockName     string  `json:"stock_name"`    // 正股名称
	Rating        string  `json:"rating"`        // 债券评级
	PublishDate   string  `json:"publish_date"`  // 申购日期
	PublishScale  float64 `json:"publish_scale"` // 发行规模(亿)
	ListedDate    string  `json:"listed_date"`   // 上市日期
	ExpireDate    string  `json:"expire_date"`   // 到期日期
	Term          int     `json:"term"`          // 期限(年)
	RateDesc      string  `json:"rate_desc"`     // 利率说明
}

// GetBaseInfo 获取单只债券基本信息
func GetBaseInfo(ctx context.Context, bondCode string) (*BondInfo, error) {
	queryParams := map[string]string{
		"reportName": "RPT_BOND_CB_LIST",
		"columns":    "ALL",
		"source":     "WEB",
		"client":     "WEB",
		"filter":     `(SECURITY_CODE="` + bondCode + `")`,
	}

	raw, err := common.DefaultClient().GetJSON(ctx, common.EastMoneyDataCenterURL, queryParams, common.HTTPHeaders)
	if err != nil {
		return nil, err
	}

	return parseBondBaseInfo(raw)
}

// GetBaseInfoMulti 获取多只债券基本信息
func GetBaseInfoMulti(ctx context.Context, bondCodes []string) ([]BondInfo, error) {
	type result struct {
		info *BondInfo
		err  error
	}

	results := make(chan result, len(bondCodes))

	var wg sync.WaitGroup
	sem := make(chan struct{}, common.MaxConnections)

	for _, code := range bondCodes {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			info, err := GetBaseInfo(ctx, code)
			results <- result{info: info, err: err}
		}(code)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var infos []BondInfo
	var errs []error

	for r := range results {
		if r.err != nil {
			errs = append(errs, r.err)
		} else if r.info != nil {
			infos = append(infos, *r.info)
		}
	}

	if len(errs) > 0 {
		return infos, fmt.Errorf("批量请求 %d 个失败: %v", len(errs), errs)
	}

	return infos, nil
}

// GetAllBaseInfo 获取全部债券基本信息
func GetAllBaseInfo(ctx context.Context) ([]BondInfo, error) {
	queryParams := map[string]string{
		"reportName": "RPT_BOND_CB_LIST",
		"columns":    "ALL",
		"source":     "WEB",
		"client":     "WEB",
		"pageSize":   "5000",
		"pageNumber": "1",
		"sortColumns": "LISTDATE",
		"sortTypes":   "-1",
	}

	raw, err := common.DefaultClient().GetJSON(ctx, common.EastMoneyDataCenterURL, queryParams, common.HTTPHeaders)
	if err != nil {
		return nil, err
	}

	return parseBondBaseInfoList(raw)
}

// parseBondBaseInfo 解析单只债券信息
func parseBondBaseInfo(raw *json.RawMessage) (*BondInfo, error) {
	var resp struct {
		Result struct {
			Data []struct {
				SECURITY_CODE      string `json:"SECURITY_CODE"`
				SECURITY_NAME_ABBR string `json:"SECURITY_NAME_ABBR"`
				CONVERT_STOCK_CODE string `json:"CONVERT_STOCK_CODE"`
				SECURITY_SHORT_NAME string `json:"SECURITY_SHORT_NAME"`
				RATING             string `json:"RATING"`
				PUBLIC_START_DATE  string `json:"PUBLIC_START_DATE"`
				ACTUAL_ISSUE_SCALE string `json:"ACTUAL_ISSUE_SCALE"`
				LISTING_DATE       string `json:"LISTING_DATE"`
				EXPIRE_DATE        string `json:"EXPIRE_DATE"`
				BOND_EXPIRE        string `json:"BOND_EXPIRE"`
				INTEREST_RATE_EXPLAIN string `json:"INTEREST_RATE_EXPLAIN"`
			} `json:"data"`
		} `json:"result"`
	}

	if err := json.Unmarshal(*raw, &resp); err != nil {
		return nil, errors.ErrParse
	}

	if len(resp.Result.Data) == 0 {
		return nil, errors.ErrNoData
	}

	d := resp.Result.Data[0]
	publishScale, _ := strconv.ParseFloat(d.ACTUAL_ISSUE_SCALE, 64)
	term, _ := strconv.Atoi(d.BOND_EXPIRE)

	return &BondInfo{
		Code:         d.SECURITY_CODE,
		Name:         d.SECURITY_NAME_ABBR,
		StockCode:    d.CONVERT_STOCK_CODE,
		StockName:    d.SECURITY_SHORT_NAME,
		Rating:       d.RATING,
		PublishDate:  d.PUBLIC_START_DATE,
		PublishScale: publishScale,
		ListedDate:   d.LISTING_DATE,
		ExpireDate:   d.EXPIRE_DATE,
		Term:         term,
		RateDesc:     d.INTEREST_RATE_EXPLAIN,
	}, nil
}

// parseBondBaseInfoList 解析债券列表
func parseBondBaseInfoList(raw *json.RawMessage) ([]BondInfo, error) {
	var resp struct {
		Result struct {
			Data []struct {
				SECURITY_CODE      string `json:"SECURITY_CODE"`
				SECURITY_NAME_ABBR string `json:"SECURITY_NAME_ABBR"`
				CONVERT_STOCK_CODE string `json:"CONVERT_STOCK_CODE"`
				SECURITY_SHORT_NAME string `json:"SECURITY_SHORT_NAME"`
				RATING             string `json:"RATING"`
				PUBLIC_START_DATE  string `json:"PUBLIC_START_DATE"`
				ACTUAL_ISSUE_SCALE string `json:"ACTUAL_ISSUE_SCALE"`
				LISTING_DATE       string `json:"LISTING_DATE"`
				EXPIRE_DATE        string `json:"EXPIRE_DATE"`
				BOND_EXPIRE        string `json:"BOND_EXPIRE"`
				INTEREST_RATE_EXPLAIN string `json:"INTEREST_RATE_EXPLAIN"`
			} `json:"data"`
		} `json:"result"`
	}

	if err := json.Unmarshal(*raw, &resp); err != nil {
		return nil, errors.ErrParse
	}

	infos := make([]BondInfo, 0, len(resp.Result.Data))
	for _, d := range resp.Result.Data {
		publishScale, _ := strconv.ParseFloat(d.ACTUAL_ISSUE_SCALE, 64)
		term, _ := strconv.Atoi(d.BOND_EXPIRE)

		infos = append(infos, BondInfo{
			Code:         d.SECURITY_CODE,
			Name:         d.SECURITY_NAME_ABBR,
			StockCode:    d.CONVERT_STOCK_CODE,
			StockName:    d.SECURITY_SHORT_NAME,
			Rating:       d.RATING,
			PublishDate:  d.PUBLIC_START_DATE,
			PublishScale: publishScale,
			ListedDate:   d.LISTING_DATE,
			ExpireDate:   d.EXPIRE_DATE,
			Term:         term,
			RateDesc:     d.INTEREST_RATE_EXPLAIN,
		})
	}

	return infos, nil
}