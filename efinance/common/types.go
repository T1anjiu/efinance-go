package common

import "time"

// MarketType 市场类型
type MarketType string

const (
	AStock              MarketType = "AStock"           // A股
	AStockIndex         MarketType = "Index"            // A股指数
	BStock              MarketType = "BStock"           // B股
	Index               MarketType = "Index"            // 沪深京指数
	STARMarket          MarketType = "23"              // 科创板
	CSIFreeFloat         MarketType = "24"              // 中证系列指数
	NEEQ                MarketType = "NEEQ"             // 京A，新三板
	BK                  MarketType = "BK"               // 板块
	Hongkong            MarketType = "HK"                // 港股
	USStock             MarketType = "UsStock"         // 美股
	LondonStock         MarketType = "LSE"              // 英股
	LondonStockIOB      MarketType = "LSEIOB"           // 伦敦交易所国际挂盘册
	UniversalIndex      MarketType = "UniversalIndex"    // 国外指数
	SIXSwiss            MarketType = "SIX"              // SIX瑞士股市
)

// KlineType K线周期类型
type KlineType int

const (
	Kline1Min   KlineType = 1   // 1分钟
	Kline5Min   KlineType = 5   // 5分钟
	Kline15Min  KlineType = 15  // 15分钟
	Kline30Min  KlineType = 30  // 30分钟
	Kline60Min  KlineType = 60  // 60分钟
	KlineDaily  KlineType = 101 // 日K
	KlineWeekly KlineType = 102 // 周K
	KlineMonthly KlineType = 103 // 月K
)

// AdjsType 复权类型
type AdjsType int

const (
	AdjsNone  AdjsType = 0 // 不复权
	AdjsFront AdjsType = 1 // 前复权
	AdjsBack  AdjsType = 2 // 后复权
)

// KlineItem K线数据项
type KlineItem struct {
	Code        string    `json:"code"`        // 代码
	Name        string    `json:"name"`        // 名称
	Date        string    `json:"date"`        // 日期 (字符串格式)
	Open        float64   `json:"open"`        // 开盘
	Close       float64   `json:"close"`       // 收盘
	High        float64   `json:"high"`        // 最高
	Low         float64   `json:"low"`         // 最低
	Volume      float64   `json:"volume"`      // 成交量
	Amount      float64   `json:"amount"`      // 成交额
	Amplitude   float64   `json:"amplitude"`   // 振幅
	ChangePCT   float64   `json:"change_pct"`  // 涨跌幅
	ChangeAmt   float64   `json:"change_amt"`  // 涨跌额
	TurnoverRate float64  `json:"turnover_rate"` // 换手率
}

// QuoteItem 行情数据项
type QuoteItem struct {
	Code           string   `json:"code"`          // 股票代码
	Name           string   `json:"name"`          // 股票名称
	ChangePCT      float64  `json:"change_pct"`    // 涨跌幅
	LatestPrice    float64  `json:"latest_price"`  // 最新价
	High           float64  `json:"high"`          // 最高
	Low            float64  `json:"low"`           // 最低
	Open           float64  `json:"open"`          // 今开
	ChangeAmt      float64  `json:"change_amt"`    // 涨跌额
	TurnoverRate   float64  `json:"turnover_rate"` // 换手率
	VolumeRatio    float64  `json:"volume_ratio"`  // 量比
	DynamicPE      float64  `json:"dynamic_pe"`    // 动态市盈率
	Volume         int64    `json:"volume"`        // 成交量
	Amount         float64  `json:"amount"`        // 成交额
	YesterdayClose float64  `json:"yesterday_close"` // 昨日收盘
	TotalMarketCap float64 `json:"total_market_cap"` // 总市值
	FlowMarketCap  float64  `json:"flow_market_cap"`  // 流通市值
	QuoteID        string   `json:"quote_id"`     // 行情ID
	MarketType     string   `json:"market_type"`  // 市场类型
}

// StockBaseInfo 股票基本信息
type StockBaseInfo struct {
	Code           string  `json:"code"`          // 股票代码
	Name           string  `json:"name"`          // 股票名称
	PE             float64 `json:"pe"`            // 市盈率(动)
	PB             float64 `json:"pb"`            // 市净率
	Industry       string  `json:"industry"`       // 所处行业
	TotalMarketCap float64 `json:"total_market_cap"` // 总市值
	FlowMarketCap  float64 `json:"flow_market_cap"`  // 流通市值
	BlockCode      string  `json:"block_code"`   // 板块编号
	ROE            float64 `json:"roe"`           // ROE
	NetMargin      float64 `json:"net_margin"`    // 净利率
	NetProfit      float64 `json:"net_profit"`    // 净利润
	GrossMargin    float64 `json:"gross_margin"`  // 毛利率
}

// FundInfo 基金信息
type FundInfo struct {
	Code           string  `json:"code"`           // 基金代码
	Name           string  `json:"name"`           // 基金名称
	NAV            float64 `json:"nav"`            // 单位净值
	ACCNAV         float64 `json:"acc_nav"`        // 累计净值
	ChangePCT      float64 `json:"change_pct"`     // 涨跌幅
	EstNAV         float64 `json:"est_nav"`        // 估算净值
	EstTime        string  `json:"est_time"`       // 估算时间
	PublishDate    string  `json:"publish_date"`  // 最新净值公开日期
}

// BondInfo 债券信息
type BondInfo struct {
	Code          string  `json:"code"`          // 债券代码
	Name          string  `json:"name"`          // 债券名称
	StockCode     string  `json:"stock_code"`    // 正股代码
	StockName     string  `json:"stock_name"`    // 正股名称
	Rating        string  `json:"rating"`        // 债券评级
	PublishDate   string  `json:"publish_date"`  // 申购日期
	PublishScale  float64 `json:"publish_scale"` // 发行规模(亿)
	ListedDate    string  `json:"listed_date"`  // 上市日期
	ExpireDate    string  `json:"expire_date"`  // 到期日期
	Term          int     `json:"term"`          // 期限(年)
}

// FlowItem 资金流向数据项
type FlowItem struct {
	Code         string    `json:"code"`          // 股票代码
	Name         string    `json:"name"`          // 股票名称
	Date         time.Time `json:"date"`         // 日期
	MainNetIn    float64   `json:"main_net_in"`   // 主力净流入
	SmallNetIn   float64   `json:"small_net_in"`  // 小单净流入
	MidNetIn     float64   `json:"mid_net_in"`    // 中单净流入
	LargeNetIn   float64   `json:"large_net_in"` // 大单净流入
	SuperNetIn   float64   `json:"super_net_in"` // 超大单净流入
	MainRatio    float64   `json:"main_ratio"`   // 主力净流入占比
	ClosePrice   float64   `json:"close_price"`  // 收盘价
	ChangePCT    float64   `json:"change_pct"`   // 涨跌幅
}

// BillboardItem 龙虎榜数据项
type BillboardItem struct {
	Code           string  `json:"code"`           // 股票代码
	Name           string  `json:"name"`           // 股票名称
	Date           string  `json:"date"`           // 上榜日期
	Interpretation string  `json:"interpretation"` // 解读
	ClosePrice     float64 `json:"close_price"`   // 收盘价
	ChangePCT      float64 `json:"change_pct"`    // 涨跌幅
	TurnoverRate   float64 `json:"turnover_rate"` // 换手率
	NetBuyAmt      float64 `json:"net_buy_amt"`   // 龙虎榜净买额
	BuyAmt         float64 `json:"buy_amt"`       // 龙虎榜买入额
	SellAmt        float64 `json:"sell_amt"`     // 龙虎榜卖出额
	TurnoverAmt    float64 `json:"turnover_amt"` // 龙虎榜成交额
	MarketTotAmt   float64 `json:"market_tot_amt"` // 市场总成交额
	NetBuyRatio    float64 `json:"net_buy_ratio"` // 净买额占总成交比
	FlowMarketCap  float64 `json:"flow_market_cap"` // 流通市值
	Reason         string  `json:"reason"`        // 上榜原因
}

// PerformanceItem 业绩表现数据项
type PerformanceItem struct {
	Code                string   `json:"code"`                // 股票代码
	Name                string   `json:"name"`                // 股票简称
	NoticeDate          string   `json:"notice_date"`         // 公告日期
	TotalRevenue        float64  `json:"total_revenue"`       // 营业收入
	RevenueGrowthYoY    float64  `json:"revenue_growth_yoy"`  // 营业收入同比增长
	RevenueGrowthQoQ    float64  `json:"revenue_growth_qoq"`  // 营业收入季度环比
	NetProfit           float64  `json:"net_profit"`          // 净利润
	NetProfitGrowthYoY  float64  `json:"net_profit_growth_yoy"` // 净利润同比增长
	NetProfitGrowthQoQ  float64  `json:"net_profit_growth_qoq"` // 净利润季度环比
	EPS                 float64  `json:"eps"`                // 每股收益
	BPS                 float64  `json:"bps"`                // 每股净资产
	ROE                 float64  `json:"roe"`                // 净资产收益率
	GrossMargin         float64  `json:"gross_margin"`        // 销售毛利率
	OpCashFlowPerShare  float64  `json:"op_cash_flow_per_share"` // 每股经营现金流量
}

// HolderInfo 股东信息
type HolderInfo struct {
	Code            string  `json:"code"`          // 股票代码
	UpdateDate      string  `json:"update_date"`   // 更新日期
	HolderCode      string  `json:"holder_code"`  // 股东代码
	HolderName      string  `json:"holder_name"`  // 股东名称
	HoldingShares   string  `json:"holding_shares"` // 持股数
	HoldingRatio    string  `json:"holding_ratio"` // 持股比例
	Change          string  `json:"change"`        // 增减
	ChangeRate      string  `json:"change_rate"`  // 变动率
}

// HolderNumInfo 股东人数信息
type HolderNumInfo struct {
	Code                string   `json:"code"`                // 股票代码
	Name                string   `json:"name"`                // 股票名称
	HolderNum           int      `json:"holder_num"`          // 股东人数
	HolderNumChange     float64  `json:"holder_num_change"`  // 股东人数增减
	ChangePCT           float64  `json:"change_pct"`         // 较上期变化百分比
	StatEndDate         string   `json:"stat_end_date"`      // 股东户数统计截止日
	AvgMarketCap        float64  `json:"avg_market_cap"`     // 户均持股市值
	AvgHoldingNum       float64  `json:"avg_holding_num"`     // 户均持股数量
	TotalMarketCap      float64  `json:"total_market_cap"`   // 总市值
	TotalShares         int64    `json:"total_shares"`       // 总股本
	NoticeDate          string   `json:"notice_date"`        // 公告日期
}
