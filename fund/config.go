package fund

import (
	"net/http"
)

var EastmoneyFundHeaders = http.Header{
	"User-Agent":  {"EMProjJijin/6.2.8 (iPhone; iOS 13.6; Scale/2.00)"},
	"GTOKEN":      {"98B423068C1F4DEF9842F82ADF08C5db"},
	"clientInfo":  {"ttjj-iPhone10,1-iOS-iOS13.6"},
	"Content-Type": {"application/x-www-form-urlencoded"},
	"Host":        {"fundmobapi.eastmoney.com"},
	"Referer":     {"https://mpservice.com/516939c37bdb4ba2b1138c50cf69a2e1/release/pages/FundHistoryNetWorth"},
}

var FundRankHeaders = http.Header{
	"Connection":    {"keep-alive"},
	"User-Agent":    {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36 Edg/87.0.664.75"},
	"Accept":        {"*/*"},
	"Referer":       {"http://fund.eastmoney.com/data/fundranking.html"},
	"Accept-Language": {"zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"},
}

var PDFHeaders = http.Header{
	"Connection":    {"keep-alive"},
	"User-Agent":    {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.128 Safari/537.36 Edg/89.0.774.77"},
	"Accept":        {"*/*"},
	"Referer":       {"http://fundf10.eastmoney.com/"},
	"Accept-Language": {"zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6"},
}
