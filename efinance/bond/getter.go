package bond

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	"github.com/efinance/efinance/efinance/common"
	"github.com/efinance/efinance/efinance/errors"
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
		return infos, errs[0]
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
				SECURITY_CODE     string `json:"SECURITY_CODE"`
				SECURITY_NAME_ABBR string `json:"SECURITY_NAME_ABBR"`
				正股代码           string `json:"正股代码"`
				正股名称           string `json:"正股名称"`
				债券评级           string `json:"债券评级"`
				申购日期           string `json:"申购日期"`
				发行规模           string `json:"发行规模"`
				上市日期           string `json:"上市日期"`
				到期日期           string `json:"到期日期"`
				期限               string `json:"期限"`
				利率说明           string `json:"利率说明"`
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
	publishScale, _ := strconv.ParseFloat(d.发行规模, 64)
	term, _ := strconv.Atoi(d.期限)

	return &BondInfo{
		Code:         d.SECURITY_CODE,
		Name:         d.SECURITY_NAME_ABBR,
		StockCode:    d.正股代码,
		StockName:    d.正股名称,
		Rating:       d.债券评级,
		PublishDate:  d.申购日期,
		PublishScale: publishScale,
		ListedDate:   d.上市日期,
		ExpireDate:   d.到期日期,
		Term:         term,
		RateDesc:     d.利率说明,
	}, nil
}

// parseBondBaseInfoList 解析债券列表
func parseBondBaseInfoList(raw *json.RawMessage) ([]BondInfo, error) {
	var resp struct {
		Result struct {
			Data []struct {
				SECURITY_CODE     string `json:"SECURITY_CODE"`
				SECURITY_NAME_ABBR string `json:"SECURITY_NAME_ABBR"`
				正股代码           string `json:"正股代码"`
				正股名称           string `json:"正股名称"`
				债券评级           string `json:"债券评级"`
				申购日期           string `json:"申购日期"`
				发行规模           string `json:"发行规模"`
				上市日期           string `json:"上市日期"`
				到期日期           string `json:"到期日期"`
				期限               string `json:"期限"`
				利率说明           string `json:"利率说明"`
			} `json:"data"`
		} `json:"result"`
	}

	if err := json.Unmarshal(*raw, &resp); err != nil {
		return nil, errors.ErrParse
	}

	infos := make([]BondInfo, 0, len(resp.Result.Data))
	for _, d := range resp.Result.Data {
		publishScale, _ := strconv.ParseFloat(d.发行规模, 64)
		term, _ := strconv.Atoi(d.期限)

		infos = append(infos, BondInfo{
			Code:         d.SECURITY_CODE,
			Name:         d.SECURITY_NAME_ABBR,
			StockCode:    d.正股代码,
			StockName:    d.正股名称,
			Rating:       d.债券评级,
			PublishDate:  d.申购日期,
			PublishScale: publishScale,
			ListedDate:   d.上市日期,
			ExpireDate:   d.到期日期,
			Term:         term,
			RateDesc:     d.利率说明,
		})
	}

	return infos, nil
}
