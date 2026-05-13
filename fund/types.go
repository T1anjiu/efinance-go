package fund

type FundNAV struct {
	FundCode string  `json:"fund_code"`
	Date     string  `json:"date"`
	UnitNAV  float64 `json:"unit_nav"`
	AccNAV   float64 `json:"acc_nav"`
	ChangeRate float64 `json:"change_rate"`
}

type FundRealtime struct {
	FundCode    string  `json:"fund_code"`
	FundName    string  `json:"fund_name"`
	LatestNAV   float64 `json:"latest_nav"`
	NAVDate     string  `json:"nav_date"`
	EstTime     string  `json:"est_time"`
	EstChangeRate float64 `json:"est_change_rate"`
}

type FundManager struct {
	FundCode      string `json:"fund_code"`
	StartDate     string `json:"start_date"`
	Company       string `json:"company"`
	ManagerName   string `json:"manager_name"`
	FundType      string `json:"fund_type"`
	FundScale     string `json:"fund_scale"`
	CurrentDate   string `json:"current_date"`
}

type Position struct {
	FundCode  string  `json:"fund_code"`
	StockCode string  `json:"stock_code"`
	StockName string  `json:"stock_name"`
	HoldRatio float64 `json:"hold_ratio"`
	Change    float64 `json:"change"`
	Date      string  `json:"date"`
}

type PeriodChange struct {
	FundCode   string  `json:"fund_code"`
	ReturnRate float64 `json:"return_rate"`
	AvgReturn  float64 `json:"avg_return"`
	Rank       int     `json:"rank"`
	TotalCount int     `json:"total_count"`
	Period     string  `json:"period"`
}

type AssetAllocation struct {
	FundCode    string  `json:"fund_code"`
	StockRatio  float64 `json:"stock_ratio"`
	BondRatio   float64 `json:"bond_ratio"`
	CashRatio   float64 `json:"cash_ratio"`
	TotalScale  float64 `json:"total_scale"`
	OtherRatio  float64 `json:"other_ratio"`
	Date        string  `json:"date"`
}

type FundBaseInfo struct {
	FundCode   string `json:"fund_code"`
	ShortName  string `json:"short_name"`
	EstabDate  string `json:"estab_date"`
	ChangeRate float64 `json:"change_rate"`
	NAV        float64 `json:"nav"`
	Company    string `json:"company"`
	NAVDate    string `json:"nav_date"`
	Comments   string `json:"comments"`
}

type IndustryDist struct {
	FundCode   string  `json:"fund_code"`
	Industry   string  `json:"industry"`
	HoldRatio  float64 `json:"hold_ratio"`
	Date       string  `json:"date"`
	MarketCap  float64 `json:"market_cap"`
}
