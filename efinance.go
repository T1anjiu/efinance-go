// Package efinance 金融数据获取库
//
// 东方财富 EastMoney API 的 Go 语言实现
//
// 支持股票、基金、债券、期货等金融数据的获取
package efinance

import (
	"github.com/efinance/efinance/efinance/bond"
	"github.com/efinance/efinance/efinance/fund"
	"github.com/efinance/efinance/efinance/futures"
	"github.com/efinance/efinance/efinance/stock"
)

// Stock 股票模块
var Stock = stock.Module{
	GetKline:           stock.GetKline,
	GetKlineMulti:      stock.GetKlineMulti,
	GetRealtimeQuotes:  stock.GetRealtimeQuotes,
	GetLatestQuote:     stock.GetLatestQuote,
	Search:             stock.Search,
	GetQuoteID:         stock.GetQuoteID,
}

// Fund 基金模块
var Fund = fund.Module{
	GetQuoteHistory:      fund.GetQuoteHistory,
	GetQuoteHistoryMulti: fund.GetQuoteHistoryMulti,
	GetRealtimeEstimate:  fund.GetRealtimeEstimate,
}

// Bond 债券模块
var Bond = bond.Module{
	GetBaseInfo:      bond.GetBaseInfo,
	GetBaseInfoMulti: bond.GetBaseInfoMulti,
	GetAllBaseInfo:   bond.GetAllBaseInfo,
}

// Futures 期货模块
var Futures = futures.Module{
	GetBaseInfo: futures.GetBaseInfo,
	GetKline:    futures.GetKline,
}
