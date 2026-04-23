# efinance-go

> 基于 [efinance](https://github.com/Micro-sheep/efinance) 的 Go 语言移植版本
>
> **数据源**：使用腾讯行情API（实时行情、K线数据）和东方财富API（股票搜索）

一个纯 Go 实现的股票行情数据获取库，专注于从公开金融API获取实时行情和K线数据。

**注意**：本项目仅供学习与研究使用，数据来源于第三方API，请遵守其服务条款。

## 功能特性

- 📈 **实时行情** - 获取股票当前价格、涨跌幅、成交量等
- 📊 **K线数据** - 支持日/周/月线及分钟K线
- 🔍 **股票搜索** - 通过关键词搜索股票
- 🗂️ **多市场支持** - 沪深A股、期货、基金、债券
- 🔄 **复权支持** - 前复权、后复权、不复权
- ⚡ **并发获取** - 支持批量股票并行请求
- 🎯 **配置管理** - 支持环境变量配置
- 💾 **智能缓存** - 内置缓存机制，支持过期清理
- 🛡️ **错误处理** - 完善的错误处理和边界检查

## 项目结构

```
efinance-go/
├── efinance/
│   ├── stock/          # 股票模块
│   │   ├── kline.go    # K线数据获取
│   │   ├── quote.go    # 实时行情获取
│   │   ├── search.go   # 股票搜索
│   │   └── module.go   # 模块定义
│   ├── futures/        # 期货模块
│   │   ├── getter.go   # 期货数据获取
│   │   └── module.go   # 模块定义
│   ├── fund/           # 基金模块
│   │   ├── getter.go   # 基金数据获取
│   │   └── module.go   # 模块定义
│   ├── bond/           # 债券模块
│   │   ├── getter.go   # 债券数据获取
│   │   └── module.go   # 模块定义
│   ├── common/         # 公共组件
│   │   ├── client.go   # HTTP客户端
│   │   ├── config.go   # 配置管理
│   │   └── types.go    # 类型定义
│   └── errors/         # 错误定义
│       └── errors.go   # 错误类型
├── efinance_test.go     # 原始测试文件
├── efinance_fix_test.go # 修复后的增强测试
├── debug_test.go        # API调试测试
├── network_test.go      # 网络连接测试
├── FIXES.md             # 修复报告
└── go.mod               # Go模块文件
```

## 快速开始

### 安装

```bash
go get github.com/T1anjiu/efinance-go/efinance@v0.1.1
```

### 获取实时行情

```go
package main

import (
    "context"
    "fmt"
    "github.com/T1anjiu/efinance-go/efinance/stock"
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
贵州茅台 (600519): ¥1415.53 (0.43%)
招商银行 (600036): ¥39.70 (-0.10%)
平安银行 (000001): ¥10.99 (0.09%)
```

### 获取K线数据

```go
package main

import (
    "context"
    "fmt"
    "github.com/T1anjiu/efinance-go/efinance/common"
    "github.com/T1anjiu/efinance-go/efinance/stock"
)

func main() {
    ctx := context.Background()
    
    params := stock.GetKlineParams{
        Code:       "600519",
        Beg:        "2024-01-01",
        End:        "2024-01-31",
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
  2024-01-02  开:1608.69 高:1611.88 低:1571.79 收:1578.70 量:32156
  2024-01-03  开:1574.80 高:1588.91 低:1570.02 收:1587.69 量:20229
  2024-01-04  开:1586.69 高:1586.69 低:1556.62 收:1562.69 量:21551
  ...
```

### 搜索股票

```go
package main

import (
    "context"
    "fmt"
    "github.com/T1anjiu/efinance-go/efinance/stock"
)

func main() {
    ctx := context.Background()
    
    results, err := stock.Search(ctx, "茅台", "")
    if err != nil {
        panic(err)
    }
    
    for _, r := range results {
        fmt.Printf("代码: %s 名称: %s 市场: %s\n", r.Code, r.Name, r.MktType)
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
| `MarketType` | `MarketType` | 市场类型筛选 | 空 |
| `SuppressErr` | bool | 遇到错误是否静默 | `false` |

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

#### `GetRealtimeQuotes(ctx, params) ([]QuoteItem, error)`

获取实时行情（支持市场类型筛选）。

```go
params := stock.QuoteParams{
    Markets: []string{"沪深A股", "上证A股"},
}
quotes, _ := stock.GetRealtimeQuotes(ctx, params)
```

### 股票搜索

#### `Search(ctx, keyword, marketType) ([]SearchResult, error)`

通过关键词搜索股票。

```go
results, _ := stock.Search(ctx, "茅台", "")
```

#### `GetQuoteID(ctx, code, marketType) (string, error)`

获取股票的行情ID。

```go
quoteID, _ := stock.GetQuoteID(ctx, "600519", "")
```

## 配置管理

### 环境变量配置

```bash
# API配置
export EASTMONEY_SEARCH_TOKEN="your_token"

# HTTP配置
export REQUEST_TIMEOUT="30s"
export MAX_CONNECTIONS="20"
export MAX_RETRIES="3"

# 缓存配置
export QUOTEID_CACHE_TTL="24h"
export SEARCH_CACHE_TTL="1h"
```

### 运行时配置

```go
import "github.com/T1anjiu/efinance-go/efinance/common"

// 修改HTTP客户端配置
common.SetRequestTimeout(30 * time.Second)
common.SetMaxConnections(20)
common.SetMaxRetries(3)

// 修改缓存配置
common.SetQuoteIDCacheTTL(24 * time.Hour)
common.SetSearchCacheTTL(1 * time.Hour)

// 清理缓存
common.DefaultQuoteIDCache().Cleanup()
common.DefaultSearchCache().Clear()
```

## 数据源

| 数据类型 | API | 域名 | 说明 |
|---------|-----|------|------|
| 实时行情 | 腾讯行情 | `qt.gtimg.cn` | 获取股票实时价格、涨跌幅等 |
| K线数据 | 腾讯/ifzq | `web.ifzq.gtimg.cn` | 获取历史K线数据 |
| 股票搜索 | 东方财富 | `searchapi.eastmoney.com` | 搜索股票代码和名称 |

## 注意事项

1. **股票代码格式**：代码可使用6位纯数字（会自动匹配市场前缀）或带前缀（如 `sh600519`）
2. **市场前缀**：上海 `sh`、深圳 `sz`、北京 `bj`
3. **网络要求**：需要能访问腾讯/东方财富域名
4. **请求频率**：请合理控制请求频率，避免对API造成压力
5. **数据准确性**：数据来源于第三方API，可能存在延迟或错误，请以官方数据为准
6. **边界检查**：代码已实现数组边界检查，不会出现越界panic
7. **缓存机制**：内置缓存机制，可减少API请求次数

## 测试

### 运行所有测试

```bash
cd efinance-go
go test -v
```

### 运行特定测试

```bash
# 功能测试
go test -v -run TestGetKline
go test -v -run TestGetLatestQuote
go test -v -run TestSearch

# 增强测试
go test -v -run TestGetKlineWithRetry
go test -v -run TestGetLatestQuoteWithBoundsCheck
go test -v -run TestConcurrentRequests
go test -v -run TestErrorHandling

# 调试测试
go test -v -run TestDebugKlineURL
go test -v -run TestMyCodeURL
go test -v -run TestNetwork
```

### 测试覆盖

- ✅ K线数据获取
- ✅ 实时行情获取
- ✅ 股票搜索
- ✅ 并发请求
- ✅ 错误处理
- ✅ 边界检查
- ✅ 网络连接

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 致谢

- 原项目 [efinance](https://github.com/Micro-sheep/efinance) - Python版本
- 腾讯财经 - 提供行情数据API
- 东方财富 - 提供搜索API
