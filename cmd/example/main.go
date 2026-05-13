package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/T1anjiu/efinance-go/internal/cache"
	"github.com/T1anjiu/efinance-go/stock"
	"github.com/T1anjiu/efinance-go/bond"
	"github.com/T1anjiu/efinance-go/futures"
	"github.com/T1anjiu/efinance-go/fund"
	"github.com/T1anjiu/efinance-go/quote"
)

func init() {
	if err := cache.Load(); err != nil {
		log.Printf("Warning: failed to load cache: %v", err)
	}
}

func main() {
	fmt.Println("=== efinance-go 示例程序 ===")
	fmt.Println()

	exampleQuoteSearch()
	exampleStockKline()
	exampleStockRealtime()
	exampleStockBaseInfo()
	exampleBondRealtime()
	exampleFuturesRealtime()
	exampleFundNAV()
}

func exampleQuoteSearch() {
	fmt.Println("--- 行情ID搜索 ---")
	q, err := quote.GetQuoteID("600519")
	if err != nil {
		log.Printf("搜索失败: %v", err)
		return
	}
	fmt.Printf("贵州茅台 行情ID: %s\n\n", q)
}

func exampleStockKline() {
	fmt.Println("--- 股票K线数据 ---")
	results, err := stock.GetQuoteHistory([]string{"600519"})
	if err != nil {
		log.Printf("获取K线失败: %v", err)
		return
	}
	for code, kline := range results {
		fmt.Printf("代码: %s, 名称: %s, K线条数: %d\n", code, kline.Name, len(kline.Bars))
		if len(kline.Bars) > 0 {
			last := kline.Bars[len(kline.Bars)-1]
			data, _ := json.Marshal(last)
			fmt.Printf("最新K线: %s\n", string(data))
		}
	}
	fmt.Println()
}

func exampleStockRealtime() {
	fmt.Println("--- 股票实时行情 ---")
	quotes, err := stock.GetRealtimeQuotes()
	if err != nil {
		log.Printf("获取实时行情失败: %v", err)
		return
	}
	fmt.Printf("共获取 %d 条实时行情\n", len(quotes))
	if len(quotes) > 3 {
		for _, q := range quotes[:3] {
			fmt.Printf("  %s(%s) 涨跌幅:%.2f%% 最新价:%.2f\n", q.Name, q.Code, q.ChangeRate, q.Price)
		}
	}
	fmt.Println()
}

func exampleStockBaseInfo() {
	fmt.Println("--- 股票基本信息 ---")
	infos, err := stock.GetBaseInfo([]string{"600519"})
	if err != nil {
		log.Printf("获取基本信息失败: %v", err)
		return
	}
	for _, info := range infos {
		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Printf("%s\n", string(data))
	}
	fmt.Println()
}

func exampleBondRealtime() {
	fmt.Println("--- 债券实时行情 ---")
	quotes, err := bond.GetRealtimeQuotes()
	if err != nil {
		log.Printf("获取债券行情失败: %v", err)
		return
	}
	fmt.Printf("共获取 %d 条债券实时行情\n", len(quotes))
	fmt.Println()
}

func exampleFuturesRealtime() {
	fmt.Println("--- 期货实时行情 ---")
	infos, err := futures.GetFuturesBaseInfo()
	if err != nil {
		log.Printf("获取期货行情失败: %v", err)
		return
	}
	fmt.Printf("共获取 %d 条期货基本信息\n", len(infos))
	if len(infos) > 3 {
		for _, info := range infos[:3] {
			fmt.Printf("  %s(%s) 行情ID:%s 市场类型:%s\n", info.FuturesName, info.FuturesCode, info.QuoteID, info.MarketType)
		}
	}
	fmt.Println()
}

func exampleFundNAV() {
	fmt.Println("--- 基金历史净值 ---")
	navs, err := fund.GetQuoteHistory("161725", 5)
	if err != nil {
		log.Printf("获取基金净值失败: %v", err)
		return
	}
	fmt.Printf("共获取 %d 条净值数据\n", len(navs))
	for _, nav := range navs {
		fmt.Printf("  %s 单位净值:%.4f 累计净值:%.4f 涨跌幅:%.2f\n", nav.Date, nav.UnitNAV, nav.AccNAV, nav.ChangeRate)
	}
	fmt.Println()
}
