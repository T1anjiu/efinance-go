package futures

import (
	"context"
	"strings"
	"sync"

	"github.com/T1anjiu/efinance-go/common"
	"golang.org/x/sync/errgroup"
)

func GetFuturesBaseInfo() ([]FuturesBaseInfo, error) {
	quotes, err := GetRealtimeQuotes()
	if err != nil {
		return nil, err
	}

	var infos []FuturesBaseInfo
	for _, q := range quotes {
		infos = append(infos, FuturesBaseInfo{
			FuturesCode: q.RealtimeQuote.Code,
			FuturesName: q.RealtimeQuote.Name,
			QuoteID:     q.RealtimeQuote.QuoteID,
			MarketType:  q.RealtimeQuote.MarketType,
		})
	}
	return infos, nil
}

func GetRealtimeQuotes() ([]FuturesQuote, error) {
	fs := common.FSDict["futures"]
	rq, err := common.GetRealtimeQuotesByFS(fs)
	if err != nil {
		return nil, err
	}
	quotes := make([]FuturesQuote, len(rq))
	for i, q := range rq {
		quotes[i] = FuturesQuote{RealtimeQuote: q}
	}
	return quotes, nil
}

func GetQuoteHistory(quoteIDs []string, opts ...common.KlineOptions) (map[string]*FuturesKline, error) {
	var kopts common.KlineOptions
	if len(opts) > 0 {
		kopts = opts[0]
	} else {
		kopts = common.DefaultKlineOptions()
	}

	results := make(map[string]*FuturesKline)
	var mu sync.Mutex

	eg, _ := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, common.MaxConnections)

	for _, qid := range quoteIDs {
		id := qid
		eg.Go(func() error {
			sem <- struct{}{}
			defer func() { <-sem }()

			r, err := common.GetQuoteHistorySingle(id, kopts)
			if err != nil {
				return nil
			}
			mu.Lock()
			results[id] = &FuturesKline{Code: r.Code, Name: r.Name, Bars: r.Bars}
			mu.Unlock()
			return nil
		})
	}
	eg.Wait()
	return results, nil
}

func GetDealDetail(quoteID string, maxCount int) ([]common.DealRecord, error) {
	return common.GetDealDetail(quoteID, maxCount)
}

func init() {
	_ = context.Background
	_ = strings.Join
}
