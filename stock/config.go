package stock

var BaseInfoFields = map[string]string{
	"f57":  "股票代码",
	"f58":  "股票名称",
	"f162": "市盈率(动)",
	"f167": "市净率",
	"f127": "所处行业",
	"f116": "总市值",
	"f117": "流通市值",
	"f198": "板块编号",
	"f173": "ROE",
	"f187": "净利率",
	"f105": "净利润",
	"f186": "毛利率",
}

var BillboardFields = map[string]string{
	"SECURITY_CODE":     "股票代码",
	"SECURITY_NAME_ABBR": "股票名称",
	"TRADE_DATE":       "上榜日期",
	"EXPLAIN":          "解读",
	"CLOSE_PRICE":      "收盘价",
	"CHANGE_RATE":      "涨跌幅",
	"TURNOVERRATE":     "换手率",
	"BILLBOARD_NET_AMT":    "龙虎榜净买额",
	"BILLBOARD_BUY_AMT":    "龙虎榜买入额",
	"BILLBOARD_SELL_AMT":   "龙虎榜卖出额",
	"BILLBOARD_DEAL_AMT":   "龙虎榜成交额",
	"ACCUM_AMOUNT":     "市场总成交额",
	"DEAL_NET_RATIO":   "净买额占总成交比",
	"DEAL_AMOUNT_RATIO": "成交额占总成交比",
	"FREE_MARKET_CAP":  "流通市值",
	"EXPLANATION":      "上榜原因",
}

var CompanyPerformanceFields = map[string]string{
	"SECURITY_CODE":      "股票代码",
	"SECURITY_NAME_ABBR": "股票简称",
	"NOTICE_DATE":        "公告日期",
	"TOTAL_OPERATE_INCOME": "营业收入",
	"YSTZ":              "营业收入同比增长",
	"YSHZ":              "营业收入季度环比",
	"PARENT_NETPROFIT":  "净利润",
	"SJLTZ":             "净利润同比增长",
	"SJLHZ":             "净利润季度环比",
	"BASIC_EPS":         "每股收益",
	"BPS":               "每股净资产",
	"WEIGHTAVG_ROE":     "净资产收益率",
	"XSMLL":             "销售毛利率",
	"MGJYXJJE":          "每股经营现金流量",
}

var HolderNumberFields = map[string]string{
	"SECURITY_CODE":      "股票代码",
	"SECURITY_NAME_ABBR": "股票名称",
	"HOLDER_NUM":         "股东人数",
	"HOLDER_NUM_RATIO":   "股东人数增减",
	"HOLDER_NUM_CHANGE":  "较上期变化百分比",
	"END_DATE":           "股东户数统计截止日",
	"AVG_MARKET_CAP":     "户均持股市值",
	"AVG_HOLD_NUM":       "户均持股数量",
	"TOTAL_MARKET_CAP":   "总市值",
	"TOTAL_A_SHARES":     "总股本",
	"HOLD_NOTICE_DATE":   "公告日期",
}

var IPOFields = map[string]string{
	"ISSUER_NAME":   "发行人全称",
	"CHECK_STATUS":  "审核状态",
	"REG_ADDRESS":   "注册地",
	"CSRC_INDUSTRY": "证监会行业",
	"RECOMMEND_ORG": "保荐机构",
	"ACCOUNT_FIRM":  "会计师事务所",
	"UPDATE_DATE":   "更新日期",
	"ACCEPT_DATE":   "受理日期",
	"TOLIST_MARKET": "拟上市地点",
}
