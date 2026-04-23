package futures

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/T1anjiu/efinance-go/efinance/common"
	"github.com/T1anjiu/efinance-go/efinance/errors"
)

// FuturesInfo 期货信息
type FuturesInfo struct {
	Code     string `json:"code"`     // 期货代码
	Name     string `json:"name"`     // 期货名称
	QuoteID  string `json:"quote_id"` // 行情ID
	MktType  string `json:"mkt_type"` // 市场类型
}

// GetBaseInfo 获取四个交易所全部期货基本信息
func GetBaseInfo(ctx context.Context) ([]FuturesInfo, error) {
	// 使用股票模块的实时行情接口获取期货信息
	queryParams := map[string]string{
		"pn":    "1",
		"pz":    "1000",
		"po":    "1",
		"np":    "1",
		"fltt":  "2",
		"invt":  "2",
		"fid":   "f12",
		"fs":    common.FSMarketDict["futures"],
		"fields": "f2,f3,f4,f5,f6,f7,f8,f12,f14,f15,f16,f17",
	}

	raw, err := common.DefaultClient().GetJSON(ctx, common.EastMoneyQuoteURL, queryParams, common.HTTPHeaders)
	if err != nil {
		return nil, err
	}

	return parseFuturesInfo(raw)
}

// parseFuturesInfo 解析期货基本信息
func parseFuturesInfo(raw *json.RawMessage) ([]FuturesInfo, error) {
	var resp struct {
		Data struct {
			Diff []struct {
				F12 string `json:"f12"` // 代码
				F14 string `json:"f14"` // 名称
				F37 string `json:"f37"` // 市场编号
			} `json:"diff"`
		} `json:"data"`
	}

	if err := json.Unmarshal(*raw, &resp); err != nil {
		return nil, errors.ErrParse
	}

	infos := make([]FuturesInfo, 0, len(resp.Data.Diff))
	for _, d := range resp.Data.Diff {
		mktType := common.MarketNumberDict[d.F37]
		quoteID := d.F37 + "." + d.F12

		infos = append(infos, FuturesInfo{
			Code:    d.F12,
			Name:    d.F14,
			QuoteID: quoteID,
			MktType: mktType,
		})
	}

	return infos, nil
}

// GetKline 获取期货K线数据
func GetKline(ctx context.Context, quoteID string, beg, end string, klt common.KlineType, fqt common.AdjsType) ([]common.KlineItem, error) {
	if beg == "" {
		beg = common.DefaultBegDate
	}
	if end == "" {
		end = common.DefaultEndDate
	}

	queryParams := map[string]string{
		"fields1": "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13",
		"fields2": "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65,f66,f67,f68",
		"beg":     beg,
		"end":     end,
		"rtntype": "6",
		"secid":   quoteID,
		"klt":     strconv.Itoa(int(klt)),
		"fqt":     strconv.Itoa(int(fqt)),
	}

	raw, err := common.DefaultClient().GetJSON(ctx, common.EastMoneyKlineURL, queryParams, common.HTTPHeaders)
	if err != nil {
		return nil, err
	}

	return parseFuturesKline(raw, quoteID)
}

// parseFuturesKline 解析期货K线数据
func parseFuturesKline(raw *json.RawMessage, quoteID string) ([]common.KlineItem, error) {
	var resp struct {
		Data struct {
			Name   string   `json:"name"`
			Klines []string `json:"klines"`
		} `json:"data"`
	}

	if err := json.Unmarshal(*raw, &resp); err != nil {
		return nil, errors.ErrParse
	}

	if len(resp.Data.Klines) == 0 {
		return nil, errors.ErrNoData
	}

	code := strings.Split(quoteID, ".")[1]
	items := make([]common.KlineItem, 0, len(resp.Data.Klines))

	for _, kline := range resp.Data.Klines {
		fields := strings.Split(kline, ",")
		if len(fields) < 13 {
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
			Volume:       volume,
			Amount:       amount,
			Amplitude:    amplitude,
			ChangePCT:    changePct,
			ChangeAmt:    changeAmt,
			TurnoverRate: turnoverRate,
			Open:         open,
			Close:        close,
			High:         high,
			Low:          low,
		})
	}

	return items, nil
}
