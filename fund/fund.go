package fund

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/T1anjiu/efinance-go/internal/client"
	"github.com/T1anjiu/efinance-go/internal/util"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

const fundAPIBase = "https://fundmobapi.eastmoney.com/FundMNewApi"

func GetQuoteHistory(fundCode string, pz int) ([]FundNAV, error) {
	if pz <= 0 {
		pz = 40000
	}
	data := url.Values{
		"FCODE":         {fundCode},
		"IsShareNet":    {"true"},
		"MobileKey":     {"1"},
		"appType":       {"ttjj"},
		"appVersion":    {"6.2.8"},
		"cToken":        {"1"},
		"deviceid":      {"1"},
		"pageIndex":     {"1"},
		"pageSize":      {fmt.Sprintf("%d", pz)},
		"plat":          {"Iphone"},
		"product":       {"EFund"},
		"serverVersion": {"6.2.8"},
		"uToken":        {"1"},
		"userId":        {"1"},
		"version":       {"6.2.8"},
	}
	u := fmt.Sprintf("%s/FundMNHisNetList", fundAPIBase)
	body, err := client.PostForm(u, data, EastmoneyFundHeaders)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	datas := result.Get("Datas")
	if !datas.Exists() || !datas.IsArray() || len(datas.Array()) == 0 {
		return nil, nil
	}

	var navs []FundNAV
	datas.ForEach(func(_, item gjson.Result) bool {
		navs = append(navs, FundNAV{
			FundCode:   fundCode,
			Date:       item.Get("FSRQ").String(),
			UnitNAV:    util.ParseFloat(item.Get("DWJZ").String()),
			AccNAV:     util.ParseFloat(item.Get("LJJZ").String()),
			ChangeRate: util.ParseFloat(item.Get("JZZZL").String()),
		})
		return true
	})
	return navs, nil
}

func GetQuoteHistoryMulti(fundCodes []string, pz int) (map[string][]FundNAV, error) {
	results := make(map[string][]FundNAV)
	var mu sync.Mutex

	eg, _ := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, 20)

	for _, code := range fundCodes {
		c := code
		eg.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()

			navs, err := GetQuoteHistory(c, pz)
			if err != nil {
				return nil
			}
			mu.Lock()
			results[c] = navs
			mu.Unlock()
			return nil
		})
	}
	eg.Wait()
	return results, nil
}

func GetRealtimeIncreaseRate(fundCodes []string) ([]FundRealtime, error) {
	data := url.Values{
		"pageIndex":   {"1"},
		"pageSize":    {"300000"},
		"Sort":        {""},
		"Fcodes":      {strings.Join(fundCodes, ",")},
		"SortColumn":  {""},
		"IsShowSE":    {"false"},
		"P":           {"F"},
		"deviceid":    {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"plat":        {"Iphone"},
		"product":     {"EFund"},
		"version":     {"6.2.8"},
	}
	u := fmt.Sprintf("%s/FundMNFInfo", fundAPIBase)
	body, err := client.PostForm(u, data, EastmoneyFundHeaders)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	datas := result.Get("Datas")
	if !datas.Exists() || !datas.IsArray() {
		return nil, nil
	}

	var rates []FundRealtime
	datas.ForEach(func(_, item gjson.Result) bool {
		rates = append(rates, FundRealtime{
			FundCode:      item.Get("FCODE").String(),
			FundName:      item.Get("SHORTNAME").String(),
			LatestNAV:     util.ParseFloat(item.Get("ACCNAV").String()),
			NAVDate:       item.Get("PDATE").String(),
			EstTime:       item.Get("GZTIME").String(),
			EstChangeRate: util.ParseFloat(item.Get("GSZZL").String()),
		})
		return true
	})
	return rates, nil
}

func GetFundCodes(ft string) ([]string, error) {
	params := url.Values{
		"op":   {"dy"},
		"dt":   {"kf"},
		"rs":   {""},
		"gs":   {"0"},
		"sc":   {"qjzf"},
		"st":   {"desc"},
		"es":   {"0"},
		"qdii": {""},
		"pi":   {"1"},
		"pn":   {"50000"},
		"dx":   {"0"},
	}
	if ft != "" {
		params.Set("ft", ft)
	}

	u := "http://fund.eastmoney.com/data/rankhandler.aspx"
	body, err := client.Get(u, params, FundRankHeaders)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`"(\d{6}),(.*?),`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	var codes []string
	for _, m := range matches {
		if len(m) >= 2 {
			codes = append(codes, m[1])
		}
	}
	return codes, nil
}

func GetFundManager(fundCode string) (*FundManager, error) {
	u := fmt.Sprintf("http://fundf10.eastmoney.com/jjjl_%s.html", fundCode)
	body, err := client.Get(u, nil, nil)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	var labels []string
	var parseLabels func(*html.Node)
	parseLabels = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "label" {
			var text string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					text += strings.TrimSpace(c.Data)
				}
				if c.Type == html.ElementNode && c.Data == "span" {
					for sc := c.FirstChild; sc != nil; sc = sc.NextSibling {
						if sc.Type == html.TextNode {
							text += strings.TrimSpace(sc.Data)
						}
					}
				}
				if c.Type == html.ElementNode && c.Data == "a" {
					for sc := c.FirstChild; sc != nil; sc = sc.NextSibling {
						if sc.Type == html.TextNode {
							text += strings.TrimSpace(sc.Data)
						}
					}
				}
			}
			if text != "" {
				labels = append(labels, text)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseLabels(c)
		}
	}
	parseLabels(doc)

	fm := &FundManager{FundCode: fundCode}
	_ = labels

	return fm, nil
}

func GetInvestPosition(fundCode string, dates []string) ([]Position, error) {
	if dates == nil {
		dates = []string{""}
	}

	var allPositions []Position
	for _, date := range dates {
		params := url.Values{
			"FCODE":         {fundCode},
			"appType":       {"ttjj"},
			"deviceid":      {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
			"plat":          {"Iphone"},
			"product":       {"EFund"},
			"serverVersion": {"6.2.8"},
			"version":       {"6.2.8"},
		}
		if date != "" {
			params.Set("DATE", date)
		}

		u := fmt.Sprintf("%s/FundMNInverstPosition", fundAPIBase)
		body, err := client.Get(u, params, EastmoneyFundHeaders)
		if err != nil {
			continue
		}

		result := gjson.ParseBytes(body)
		stocks := result.Get("Datas.fundStocks")
		pubDate := result.Get("Expansion").String()

		if !stocks.Exists() || !stocks.IsArray() {
			continue
		}

		stocks.ForEach(func(_, item gjson.Result) bool {
			allPositions = append(allPositions, Position{
				FundCode:  fundCode,
				StockCode: item.Get("GPDM").String(),
				StockName: item.Get("GPJC").String(),
				HoldRatio: util.ParseFloat(item.Get("JZBL").String()),
				Change:    util.ParseFloat(item.Get("PCTNVCHG").String()),
				Date:      pubDate,
			})
			return true
		})
	}
	return allPositions, nil
}

func GetPeriodChange(fundCode string) ([]PeriodChange, error) {
	params := url.Values{
		"AppVersion":  {"6.3.8"},
		"FCODE":       {fundCode},
		"MobileKey":   {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"OSVersion":   {"14.3"},
		"deviceid":    {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"passportid":  {"3061335960830820"},
		"plat":        {"Iphone"},
		"product":     {"EFund"},
		"version":     {"6.3.6"},
	}

	u := fmt.Sprintf("%s/FundMNPeriodIncrease", fundAPIBase)
	body, err := client.Get(u, params, EastmoneyFundHeaders)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	datas := result.Get("Datas")
	if !datas.Exists() || !datas.IsArray() {
		return nil, nil
	}

	titles := []string{"近一周", "近一月", "近三月", "近六月", "近一年", "近两年", "近三年", "近五年", "今年以来", "成立以来"}

	var changes []PeriodChange
	i := 0
	datas.ForEach(func(_, item gjson.Result) bool {
		title := ""
		if i < len(titles) {
			title = titles[i]
		}
		changes = append(changes, PeriodChange{
			FundCode:   fundCode,
			ReturnRate: util.ParseFloat(item.Get("syl").String()),
			AvgReturn:  util.ParseFloat(item.Get("avg").String()),
			Rank:       int(item.Get("rank").Int()),
			TotalCount: int(item.Get("sc").Int()),
			Period:     title,
		})
		i++
		return true
	})
	return changes, nil
}

func GetPublicDates(fundCode string) ([]string, error) {
	params := url.Values{
		"FCODE":         {fundCode},
		"appVersion":    {"6.3.8"},
		"deviceid":      {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"plat":          {"Iphone"},
		"product":       {"EFund"},
		"serverVersion": {"6.3.6"},
		"version":       {"6.3.8"},
	}

	u := fmt.Sprintf("%s/FundMNIVInfoMultiple", fundAPIBase)
	body, err := client.Get(u, params, EastmoneyFundHeaders)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	datas := result.Get("Datas")
	if !datas.Exists() {
		return nil, nil
	}

	var dates []string
	datas.ForEach(func(_, v gjson.Result) bool {
		dates = append(dates, v.String())
		return true
	})
	return dates, nil
}

func GetTypesPercentage(fundCode string, dates []string) ([]AssetAllocation, error) {
	if dates == nil {
		dates = []string{""}
	}

	var all []AssetAllocation
	for _, date := range dates {
		params := url.Values{
			"FCODE":         {fundCode},
			"OSVersion":     {"14.3"},
			"appVersion":    {"6.3.8"},
			"deviceid":      {"3EA024C2-7F21-408B-95E4-383D38160FB3"},
			"plat":          {"Iphone"},
			"product":       {"EFund"},
			"serverVersion": {"6.3.6"},
			"version":       {"6.3.8"},
		}
		if date != "" {
			params.Set("DATE", date)
		}

		u := fmt.Sprintf("%s/FundMNAssetAllocationNew", fundAPIBase)
		body, err := client.Get(u, params, EastmoneyFundHeaders)
	if err != nil {
			continue
	}

		result := gjson.ParseBytes(body)
		datas := result.Get("Datas")
		if !datas.Exists() || !datas.IsArray() || len(datas.Array()) == 0 {
			continue
		}

		datas.ForEach(func(_, item gjson.Result) bool {
			all = append(all, AssetAllocation{
				FundCode:   fundCode,
				StockRatio: util.ParseFloat(item.Get("GP").String()),
				BondRatio:  util.ParseFloat(item.Get("ZQ").String()),
				CashRatio:  util.ParseFloat(item.Get("HB").String()),
				TotalScale: util.ParseFloat(item.Get("JZC").String()),
				OtherRatio: util.ParseFloat(item.Get("QT").String()),
			})
			return true
		})
	}
	return all, nil
}

func GetBaseInfo(fundCodes []string) ([]FundBaseInfo, error) {
	if len(fundCodes) == 1 {
		info, err := getBaseInfoSingle(fundCodes[0])
		if err != nil {
			return nil, err
		}
		if info == nil {
			return nil, nil
		}
		return []FundBaseInfo{*info}, nil
	}

	var results []FundBaseInfo
	var mu sync.Mutex

	eg, _ := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, 20)

	for _, code := range fundCodes {
		c := code
		eg.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()

			info, err := getBaseInfoSingle(c)
			if err != nil || info == nil {
				return nil
			}
			mu.Lock()
			results = append(results, *info)
			mu.Unlock()
			return nil
		})
	}
	eg.Wait()
	return results, nil
}

func getBaseInfoSingle(fundCode string) (*FundBaseInfo, error) {
	params := url.Values{
		"FCODE":    {fundCode},
		"deviceid": {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
		"plat":     {"Iphone"},
		"product":  {"EFund"},
		"version":  {"6.3.8"},
	}

	u := fmt.Sprintf("%s/FundMNNBasicInformation", fundAPIBase)
	body, err := client.Get(u, params, EastmoneyFundHeaders)
	if err != nil {
		return nil, err
	}

	result := gjson.ParseBytes(body)
	datas := result.Get("Datas")
	if !datas.Exists() {
		return nil, nil
	}

	item := datas.Array()[0]
	info := &FundBaseInfo{
		FundCode:   item.Get("FCODE").String(),
		ShortName:  item.Get("SHORTNAME").String(),
		EstabDate:  item.Get("ESTABDATE").String(),
		ChangeRate: util.ParseFloat(item.Get("RZDF").String()),
		NAV:        util.ParseFloat(item.Get("DWJZ").String()),
		Company:    item.Get("JJGS").String(),
		NAVDate:    item.Get("FSRQ").String(),
		Comments:   strings.TrimSpace(item.Get("COMMENTS").String()),
	}
	info.Comments = strings.ReplaceAll(info.Comments, "\n", " ")
	return info, nil
}

func GetIndustryDistribution(fundCode string, dates []string) ([]IndustryDist, error) {
	if dates == nil {
		dates = []string{""}
	}

	var all []IndustryDist
	for _, date := range dates {
		params := url.Values{
			"FCODE":         {fundCode},
			"OSVersion":     {"14.4"},
			"appVersion":    {"6.3.8"},
			"deviceid":      {"3EA024C2-7F22-408B-95E4-383D38160FB3"},
			"plat":          {"Iphone"},
			"product":       {"EFund"},
			"serverVersion": {"6.3.6"},
			"version":       {"6.3.8"},
		}
		if date != "" {
			params.Set("DATE", date)
		}

		u := fmt.Sprintf("%s/FundMNSectorAllocation", fundAPIBase)
		body, err := client.Get(u, params, EastmoneyFundHeaders)
		if err != nil {
			continue
		}

		result := gjson.ParseBytes(body)
		datas := result.Get("Datas")
		if !datas.Exists() || !datas.IsArray() {
			continue
		}

		datas.ForEach(func(_, item gjson.Result) bool {
			all = append(all, IndustryDist{
				FundCode:  fundCode,
				Industry:  item.Get("HYMC").String(),
				HoldRatio: util.ParseFloat(item.Get("ZJZBL").String()),
				Date:      item.Get("FSRQ").String(),
				MarketCap: util.ParseFloat(item.Get("SZ").String()),
			})
			return true
		})
	}
	return all, nil
}

func GetPDFReports(fundCode string, maxCount int, saveDir string) error {
	if maxCount <= 0 {
		maxCount = 12
	}
	if saveDir == "" {
		saveDir = "pdf"
	}

	params := url.Values{
		"fundcode":  {fundCode},
		"pageIndex": {"1"},
		"pageSize":  {"200000"},
		"type":      {"3"},
	}
	u := "http://api.fund.eastmoney.com/f10/JJGG"
	body, err := client.Get(u, params, PDFHeaders)
	if err != nil {
		return err
	}

	result := gjson.ParseBytes(body)
	data := result.Get("Data")
	if !data.Exists() || !data.IsArray() {
		return nil
	}

	items := data.Array()
	start := len(items) - maxCount
	if start < 0 {
		start = 0
	}

	dir := saveDir + "/" + fundCode
	os.MkdirAll(dir, 0755)

	for _, item := range items[start:] {
		id := item.Get("ID").String()
		title := item.Get("TITLE").String()
		downloadURL := fmt.Sprintf("http://pdf.dfcfw.com/pdf/H2_%s_1.pdf", id)

		resp, err := client.DefaultClient.Get(downloadURL)
		if err != nil {
			continue
		}
		f, err := os.Create(fmt.Sprintf("%s/%s.pdf", dir, title))
		if err != nil {
			resp.Body.Close()
			continue
		}
		io.Copy(f, resp.Body)
		f.Close()
		resp.Body.Close()
	}
	return nil
}

func init() {
	_ = context.Background
	_ = time.Now
}
