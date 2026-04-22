# efinance-go

> 基于 [efinance](https://github.com/Micro-sheep/efinance) 的 Go 语言移植版本
>
> 数据源变更：原项目使用东方财富行情，本项目改用腾讯行情API

一个纯 Go 实现的股票行情数据获取库，专注于从公开金融API获取实时行情和K线数据。

**注意**：本项目仅供学习与研究使用，数据来源于第三方API，请遵守其服务条款。

## 功能特性

- 📈 **实时行情** - 获取股票当前价格、涨跌幅、成交量等
- 📊 **K线数据** - 支持日/周/月线及分钟K线
- 🔍 **股票搜索** - 通过关键词搜索股票
- 🗂️ **多市场支持** - 沪深A股、期货、基金、债券
- 🔄 **复权支持** - 前复权、后复权、不复权
- ⚡ **并发获取** - 支持批量股票并行请求

## 项目结构

```
efinance-go/
├── efinance/
│   ├── stock/          # 股票模块
│   │   ├── kline.go    # K线数据
│   │   ├── quote.go    # 实时行情
│   │   └── search.go   # 股票搜索
│   ├── futures/        # 期货模块
│   ├── fund/           # 基金模块
│   ├── bond/           # 债券模块
│   ├── common/         # 公共组件 (HTTP客户端、类型定义)
│   └── errors/         # 错误定义
├── efinance_test.go     # 测试文件
└── go.mod
```

## 快速开始

### 安装

```bash
go get github.com/efinance/efinance/efinance/stock
```

### 获取实时行情

```go
package main

import (
    "context"
    "fmt"
    "github.com/efinance/efinance/efinance/stock"
)

func main() {
    ctx := context.Background()
    
    // 获取指定股票的实时行情
    quotes, err := stock.GetLatestQuote(ctx, []string{"600519", "600036", "000001"})
    if err != nil {
        panic(err)
    }
    
    for _, q := range quotes {
        fmt.Printf("%s (%s): ¥%.2f (%.2f%%)\n", 
            q.Name, q.Code, q.LatestPrice, q.ChangePCT)
    }
}
```

**输出示例**：
```
贵州茅台 (600519): ¥1409.50 (-0.18%)
招商银行 (600036): ¥39.74 (-0.33%)
平安银行 (000001): ¥10.98 (-0.90%)
```

### 获取K线数据

```go
package main

import (
    "context"
    "fmt"
    "github.com/efinance/efinance/efinance/common"
    "github.com/efinance/efinance/efinance/stock"
)

func main() {
    ctx := context.Background()
    
    params := stock.GetKlineParams{
        Code:       "600519",
        Beg:        "2024-01-01",
        End:        "2024-12-31",
        KlineType:  common.KlineDaily,
        AdjustType: common.AdjsFront,
    }
    
    result, err := stock.GetKline(ctx, params)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("股票: %s (%s), 共 %d 条数据\n", result.Name, result.Code, len(result.Items))
    
    // 显示最新5条
    items := result.Items
    if len(items) > 5 {
        items = items[len(items)-5:]
    }
    for _, item := range items {
        fmt.Printf("  %s  开:%.2f 高:%.2f 低:%.2f 收:%.2f 量:%v\n",
            item.Date, item.Open, item.High, item.Low, item.Close, item.Volume)
    }
}
```

**输出示例**：
```
股票: 贵州茅台 (600519), 共 22 条数据
  2024-11-29  开:1485.00 高:1500.00 低:1478.50 收:1495.20 量:3.20万
  2024-12-02  开:1490.00 高:1498.00 低:1482.00 收:1488.50 量:2.85万
  2024-12-03  开:1488.00 高:1495.00 低:1480.00 收:1492.30 量:3.10万
  ...
```

### 搜索股票

```go
package main

import (
    "context"
    "fmt"
    "github.com/efinance/efinance/efinance/stock"
)

func main() {
    ctx := context.Background()
    
    results, err := stock.SearchStock(ctx, "贵州茅台")
    if err != nil {
        panic(err)
    }
    
    for _, r := range results {
        fmt.Printf("代码: %s 名称: %s 市场: %s\n", r.Code, r.Name, r.Market)
    }
}
```

## API 参考

### K线数据

#### `GetKline(ctx, params) (*KlineResult, error)`

获取单只股票的K线数据。

**参数** `GetKlineParams`:
| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `Code` | string | 股票代码（6位数字）或名称 | 必填 |
| `Beg` | string | 开始日期 `YYYY-MM-DD` | `2024-01-01` |
| `End` | string | 结束日期 `YYYY-MM-DD` | 空（当前） |
| `KlineType` | `KlineType` | K线周期 | `KlineDaily` |
| `AdjustType` | `AdjsType` | 复权类型 | `AdjsFront` |

**K线周期** `KlineType`:
```go
KlineDaily   // 日K
KlineWeekly  // 周K
KlineMonthly // 月K
Kline1Min    // 1分钟
Kline5Min    // 5分钟
Kline15Min   // 15分钟
Kline30Min   // 30分钟
Kline60Min   // 60分钟
```

**复权类型** `AdjsType`:
```go
AdjsFront  // 前复权（默认）
AdjsBack   // 后复权
AdjsNone   // 不复权
```

**返回** `KlineResult`:
```go
type KlineResult struct {
    Code  string             // 股票代码
    Name  string             // 股票名称
    Items []common.KlineItem // K线数据
}
```

**K线数据项** `KlineItem`:
```go
type KlineItem struct {
    Code   string  // 股票代码
    Name   string  // 股票名称
    Date   string  // 日期 "YYYY-MM-DD"
    Open   float64 // 开盘价
    Close  float64 // 收盘价
    High   float64 // 最高价
    Low    float64 // 最低价
    Volume float64 // 成交量
    Amount float64 // 成交额
}
```

#### `GetKlineMulti(ctx, params, workers) (map[string]*KlineResult, error)`

批量获取多只股票K线数据。

```go
params := []stock.GetKlineParams{
    {Code: "600519", KlineType: common.KlineDaily},
    {Code: "600036", KlineType: common.KlineDaily},
}
results, _ := stock.GetKlineMulti(ctx, params, 5)
```

### 实时行情

#### `GetLatestQuote(ctx, codes) ([]QuoteItem, error)`

获取指定股票的实时行情。

**参数**:
- `codes`: 股票代码列表，支持6位纯数字或带市场前缀

```go
quotes, _ := stock.GetLatestQuote(ctx, []string{"sh600519", "sz000001"})
// 或
quotes, _ := stock.GetLatestQuote(ctx, []string{"600519", "000001"})
```

**返回** `QuoteItem`:
```go
type QuoteItem struct {
    Code           string  // 股票代码
    Name           string  // 股票名称
    LatestPrice    float64 // 最新价
    ChangePCT      float64 // 涨跌幅 (%)
    ChangeAmt      float64 // 涨跌额
    Open           float64 // 开盘价
    High           float64 // 最高价
    Low            float64 // 最低价
    Volume         int64   // 成交量
    Amount         float64 // 成交额
    TurnoverRate   float64 // 换手率
    DynamicPE      float64 // 动态市盈率
    // ...更多字段
}
```

### 股票搜索

#### `SearchStock(ctx, keyword) ([]SearchResult, error)`

通过关键词搜索股票。

```go
results, _ := stock.SearchStock(ctx, "茅台")
```

## 数据源

| 数据类型 | API | 域名 |
|---------|-----|------|
| 实时行情 | 腾讯行情 | `qt.gtimg.cn` |
| K线数据 | 腾讯/ifzq | `web.ifzq.gtimg.cn` |
| 股票搜索 | 东方财富 | `search-codename.eastmoney.com` |

## 注意事项

1. **股票代码格式**：代码可使用6位纯数字（会自动匹配市场前缀）或带前缀（如 `sh600519`）
2. **市场前缀**：上海 `sh`、深圳 `sz`、北京 `bj`
3. **网络要求**：需要能访问腾讯/东方财富域名
4. **请求频率**：请合理控制请求频率，避免对API造成压力

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test -v -run "TestGetKline"
go test -v -run "TestRealtimeQuote"
```

## 许可证

MIT License