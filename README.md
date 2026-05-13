# efinance-go

<p align="center">
  <img src="https://img.shields.io/badge/version-v0.2.0-00ADD8?style=flat" alt="Version">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License">
  <img src="https://img.shields.io/badge/based_on-efinance_0.5.8-blue.svg" alt="Based On">
</p>

基于东方财富 API 的 Golang 金融数据库，提供 **股票、基金、债券、期货** 的行情数据获取能力。

> 本项目是 [efinance](https://github.com/Micro-sheep/efinance) 的 Go 语言重构版本，完整迁移了原 Python 库的全部功能。

**郑重声明：本项目仅供学习交流使用，不得用于商业用途。**

---

## 特性

- **多品种覆盖**：沪深A股、港股、美股、期货、可转债、公募基金
- **丰富的数据类型**：K线（日/周/月/分钟）、实时行情、资金流向、股东信息、龙虎榜、指数成分股等
- **行情ID自动解析**：直接输入股票代码或名称，自动解析为东方财富行情ID
- **双层缓存**：内存 + 本地文件，TTL 72小时自动过期
- **高并发**：errgroup + semaphore 限流，支持批量数据并发抓取
- **自动重试**：请求失败自动重试3次，间隔1秒
- **连接池复用**：全局 HTTP 客户端，连接池上限50
- **轻量依赖**：仅4个外部依赖

---

## 安装

```bash
go get github.com/T1anjiu/efinance-go
```

在代码中导入：

```go
import (
    "github.com/T1anjiu/efinance-go/stock"
    "github.com/T1anjiu/efinance-go/fund"
    "github.com/T1anjiu/efinance-go/bond"
    "github.com/T1anjiu/efinance-go/futures"
    "github.com/T1anjiu/efinance-go/quote"
    "github.com/T1anjiu/efinance-go/internal/cache"
)
```

> 缓存目录默认为 `~/.efinance/`，可通过环境变量 `EFINANCE_CACHE_DIR` 自定义。

---

## 快速开始

### 初始化缓存

```go
package main

import (
    "log"
    "github.com/T1anjiu/efinance-go/internal/cache"
    "github.com/T1anjiu/efinance-go/stock"
)

func init() {
    if err := cache.Load(); err != nil {
        log.Printf("加载缓存失败: %v", err)
    }
}

func main() {
    // ...
}
```

### 股票K线数据

```go
// 获取单只或多只股票日K线
results, _ := stock.GetQuoteHistory([]string{"600519", "000858"})
for code, kline := range results {
    fmt.Printf("%s(%s) 共 %d 条K线\n", code, kline.Name, len(kline.Bars))
    if len(kline.Bars) > 0 {
        last := kline.Bars[len(kline.Bars)-1]
        fmt.Printf("  最新: 日期=%s 开=%.2f 收=%.2f 高=%.2f 低=%.2f\n",
            last.Date, last.Open, last.Close, last.High, last.Low)
    }
}
```

### 实时行情

```go
// 全市场实时行情，按涨跌幅排序
quotes, _ := stock.GetRealtimeQuotes("沪深A股")
for _, q := range quotes[:5] {
    fmt.Printf("%s(%s) 涨跌幅:%.2f%% 最新价:%.2f\n",
        q.Name, q.Code, q.ChangeRate, q.Price)
}

// 指定市场
quotes, _ = stock.GetRealtimeQuotes("创业板")
quotes, _ = stock.GetRealtimeQuotes("港股")
quotes, _ = stock.GetRealtimeQuotes("ETF")
```

### 股票基本信息

```go
infos, _ := stock.GetBaseInfo([]string{"600519", "300750"})
for _, info := range infos {
    fmt.Printf("%s(%s) PE:%.2f PB:%.2f ROE:%.2f 行业:%s\n",
        info.Name, info.Code, info.PERatio, info.PBRatio, info.ROE, info.Industry)
}
```

### 资金流向

```go
// 历史资金流向
records, code, name, _ := stock.GetHistoryBill("600519")
fmt.Printf("%s(%s) 共 %d 条资金流向记录\n", name, code, len(records))

// 当日分钟级资金流向
todayRecords, code, name, _ := stock.GetTodayBill("600519")
fmt.Printf("%s(%s) 当日共 %d 条分钟级资金流向\n", name, code, len(todayRecords))
```

### 龙虎榜

```go
records, _ := stock.GetDailyBillboard("2024-01-01", "2024-01-31")
for _, r := range records {
    fmt.Printf("%s(%s) 上榜日期:%s 净买额:%.0f 涨跌幅:%.2f%%\n",
        r.StockName, r.StockCode, r.TradeDate, r.NetBuyAmt, r.ChangeRate)
}
```

### 行情快照（含五档买卖盘）

```go
snap, _ := stock.GetQuoteSnapshot("600519")
fmt.Printf("买一: %.2f (×%.0f)  卖一: %.2f (×%.0f)\n",
    snap.Buy1, snap.Buy1Count, snap.Sell1, snap.Sell1Count)
```

### 可转债

```go
// 可转债实时行情
quotes, _ := bond.GetRealtimeQuotes()
fmt.Printf("共 %d 只可转债\n", len(quotes))

// 可转债基本信息
infos, _ := bond.GetBaseInfo([]string{"123111", "113050"})
for _, info := range infos {
    fmt.Printf("%s(%s) 正股:%s(%s) 评级:%s\n",
        info.BondName, info.BondCode, info.StockName, info.StockCode, info.Rating)
}

// 全部可转债列表
all, _ := bond.GetAllBaseInfo()
fmt.Printf("共 %d 只可转债\n", len(all))
```

### 期货

```go
// 全部期货，行情ID可直接用于K线查询
infos, _ := futures.GetFuturesBaseInfo()
for _, info := range infos[:5] {
    fmt.Printf("%s(%s) 行情ID:%s 市场:%s\n",
        info.FuturesName, info.FuturesCode, info.QuoteID, info.MarketType)
}

// 期货K线，直接使用行情ID
results, _ := futures.GetQuoteHistory([]string{"115.ZCM", "114.jm"})
for id, kline := range results {
    fmt.Printf("%s(%s) K线:%d条\n", id, kline.Name, len(kline.Bars))
}

// 期货实时行情
quotes, _ := futures.GetRealtimeQuotes()
fmt.Printf("共 %d 个期货合约\n", len(quotes))
```

### 基金

```go
// 历史净值
navs, _ := fund.GetQuoteHistory("161725", 10)
for _, nav := range navs {
    fmt.Printf("%s 单位净值:%.4f 累计净值:%.4f 涨跌幅:%.2f%%\n",
        nav.Date, nav.UnitNAV, nav.AccNAV, nav.ChangeRate)
}

// 实时估算涨跌幅
rates, _ := fund.GetRealtimeIncreaseRate([]string{"161725", "005827"})
for _, r := range rates {
    fmt.Printf("%s(%s) 估算涨跌幅:%.2f%%\n",
        r.FundName, r.FundCode, r.EstChangeRate)
}

// 基金持仓
positions, _ := fund.GetInvestPosition("161725", nil)
for _, p := range positions[:5] {
    fmt.Printf("  %s(%s) 占比:%.2f%%\n", p.StockName, p.StockCode, p.HoldRatio)
}

// 基金阶段涨跌幅
changes, _ := fund.GetPeriodChange("161725")
for _, c := range changes {
    fmt.Printf("  %s: %.2f%% (同类平均:%.2f%% 排名:%d/%d)\n",
        c.Period, c.ReturnRate, c.AvgReturn, c.Rank, c.TotalCount)
}

// 行业分布
dist, _ := fund.GetIndustryDistribution("161725", nil)
for _, d := range dist {
    fmt.Printf("  %s: %.2f%%\n", d.Industry, d.HoldRatio)
}

// 全部基金代码
codes, _ := fund.GetFundCodes("")   // 全部类型
codes, _ = fund.GetFundCodes("gp")  // 股票型
codes, _ = fund.GetFundCodes("etf") // ETF
codes, _ = fund.GetFundCodes("hh")  // 混合型
```

### 行情ID搜索

```go
// 代码/名称 → 行情ID
id, _ := quote.GetQuoteID("600519")       // → "1.600519"
id, _ = quote.GetQuoteID("贵州茅台")        // → "1.600519"

// 搜索多个结果
results, _ := quote.SearchMulti("茅台", 5)
for _, q := range results {
    fmt.Printf("%s(%s) 行情ID:%s 类型:%s\n",
        q.Name, q.Code, q.QuoteID, q.SecurityTypeName)
}
```

---

## API 概览

### 股票 API (`stock`)

| 函数 | 说明 |
|---|---|
| `GetQuoteHistory` | 获取K线数据（日/周/月/分钟，支持前复权/后复权） |
| `GetRealtimeQuotes` | 获取全市场或指定市场实时行情 |
| `GetLatestQuote` | 获取指定股票的实时行情 |
| `GetBaseInfo` | 获取基本信息（PE/PB/ROE/毛利率等） |
| `GetHistoryBill` | 获取历史资金流向 |
| `GetTodayBill` | 获取当日分钟级资金流向 |
| `GetTop10Holders` | 获取前十大流通股东 |
| `GetAllReportDates` | 获取全部季报日期 |
| `GetAllCompanyPerformance` | 获取全市场季度业绩 |
| `GetLatestHolderNumber` | 获取股东人数变化 |
| `GetDailyBillboard` | 获取龙虎榜详情 |
| `GetIndexMembers` | 获取指数成分股及权重 |
| `GetLatestIPOInfo` | 获取IPO审核状态 |
| `GetQuoteSnapshot` | 获取行情快照（含五档买卖盘） |
| `GetDealDetail` | 获取成交明细（逐笔） |
| `GetBelongBoard` | 获取所属板块 |

### 债券 API (`bond`)

| 函数 | 说明 |
|---|---|
| `GetBaseInfo` | 获取可转债基本信息 |
| `GetAllBaseInfo` | 获取全部可转债列表 |
| `GetRealtimeQuotes` | 获取可转债实时行情 |
| `GetQuoteHistory` | 获取可转债K线数据 |
| `GetHistoryBill` | 获取历史资金流向 |
| `GetTodayBill` | 获取当日资金流向 |
| `GetDealDetail` | 获取成交明细 |

### 期货 API (`futures`)

| 函数 | 说明 |
|---|---|
| `GetFuturesBaseInfo` | 获取全部期货基本信息 |
| `GetRealtimeQuotes` | 获取期货实时行情 |
| `GetQuoteHistory` | 获取期货K线（直接使用行情ID） |
| `GetDealDetail` | 获取成交明细 |

### 基金 API (`fund`)

| 函数 | 说明 |
|---|---|
| `GetQuoteHistory` | 获取历史净值 |
| `GetQuoteHistoryMulti` | 批量获取净值 |
| `GetRealtimeIncreaseRate` | 获取实时估算涨跌幅 |
| `GetFundCodes` | 获取全部基金代码（可按类型筛选） |
| `GetFundManager` | 获取基金经理信息 |
| `GetInvestPosition` | 获取持仓股票及占比 |
| `GetPeriodChange` | 获取阶段涨跌幅 |
| `GetPublicDates` | 获取历史公开持仓日期 |
| `GetTypesPercentage` | 获取资产类型占比 |
| `GetBaseInfo` | 获取基金基本信息 |
| `GetIndustryDistribution` | 获取行业分布 |
| `GetPDFReports` | 下载PDF报告文件 |

### 行情ID (`quote`)

| 函数 | 说明 |
|---|---|
| `GetQuoteID` | 代码/名称 → 行情ID |
| `Search` | 搜索单个证券 |
| `SearchMulti` | 搜索多个候选证券 |

---

## 架构

```
efinance-go/
├── internal/
│   ├── client/      # 全局HTTP客户端（连接池50、重试3次、超时30s）
│   ├── cache/       # 双层缓存（sync.Map + JSON文件，TTL 72h）
│   └── util/        # 解析工具（数值转换、K线拆分）
├── common/
│   ├── config.go    # 全局常量（FSDict、字段映射、市场编号）
│   └── getter.go    # 通用数据获取层
├── quote/           # 行情ID解析（搜索API → secid）
├── stock/           # 股票模块
├── bond/            # 债券模块
├── futures/         # 期货模块
├── fund/            # 基金模块（独立请求头，伪装iPhone）
└── cmd/example/     # 示例程序
```

**数据流示意：**

```
用户输入 "600519"
    ↓
quote.GetQuoteID("600519")
    ├─ 查内存缓存 sync.Map
    ├─ 查本地文件 ~/.efinance/search-cache.json
    └─ 调用 searchapi.eastmoney.com → "1.600519"
    ↓
common.GetQuoteHistorySingle("1.600519", opts)
    ↓
HTTP GET push2his.eastmoney.com
    ↓
gjson 提取 klines → strings.Split → KlineBar struct
    ↓
返回 []KlineBar
```

---

## 市场行情类型

| 可用值 | 说明 |
|---|---|
| `stock` / `沪深A股` | 沪深京A股 |
| `沪A` | 沪市A股 |
| `深A` | 深市A股 |
| `北A` | 北证A股 |
| `创业板` | 创业板 |
| `科创板` | 科创板 |
| `美股` | 美股 |
| `港股` | 港股 |
| `英股` | 英股 |
| `新股` | 沪深新股 |
| `ETF` | ETF基金 |
| `LOF` | LOF基金 |
| `中概股` | 中国概念股 |
| `可转债` | 可转债 |
| `期货` | 期货 |
| `行业板块` | 行业板块 |
| `概念板块` | 概念板块 |
| `沪深系列指数` | 沪深系列指数 |

---

## 技术栈对应

| Python (efinance) | Go (efinance-go) |
|---|---|
| `requests` | `net/http` |
| `pandas DataFrame` | 结构化 `struct slice` |
| `jsonpath` | `github.com/tidwall/gjson` |
| `multitasking` | `golang.org/x/sync/errgroup` |
| `retry` | `github.com/avast/retry-go/v4` |
| `beautifulsoup4` | `golang.org/x/net/html` |

---

## K线周期说明

| klt 值 | 含义 |
|---|---|
| `1` | 1分钟 |
| `5` | 5分钟 |
| `15` | 15分钟 |
| `30` | 30分钟 |
| `60` | 60分钟 |
| `101` | 日K |
| `102` | 周K |
| `103` | 月K |

## 复权方式

| fqt 值 | 含义 |
|---|---|
| `0` | 不复权 |
| `1` | 前复权 |
| `2` | 后复权 |

---

## 相关项目

- [efinance](https://github.com/Micro-sheep/efinance) — 原始 Python 版本，作者 [micro sheep](https://github.com/Micro-sheep)
- [efinance-go](https://github.com/T1anjiu/efinance-go) — 本项目 Go 语言重构版

---

## License

MIT © 2021 micro sheep, 2026 efinance-go contributors