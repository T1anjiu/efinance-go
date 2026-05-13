package bond

import "github.com/T1anjiu/efinance-go/common"

type BondBaseInfo struct {
	BondCode      string  `json:"bond_code"`
	BondName      string  `json:"bond_name"`
	StockCode     string  `json:"stock_code"`
	StockName     string  `json:"stock_name"`
	Rating        string  `json:"rating"`
	SubscribeDate string  `json:"subscribe_date"`
	IssueScale    float64 `json:"issue_scale"`
	WinRate       float64 `json:"win_rate"`
	ListingDate   string  `json:"listing_date"`
	ExpireDate    string  `json:"expire_date"`
	Term          float64 `json:"term"`
	InterestDesc  string  `json:"interest_desc"`
}

type BondQuote struct {
	common.RealtimeQuote
}

type BondKline struct {
	Code string           `json:"code"`
	Name string           `json:"name"`
	Bars []common.KlineBar `json:"bars"`
}

type BondBillRecord = common.BillRecord
type BondTodayBillRecord = common.TodayBillRecord
