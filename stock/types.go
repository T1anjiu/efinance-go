package stock

type MarketType int

const (
	MarketTypeA  MarketType = iota
	MarketTypeHK
	MarketTypeUS
	MarketTypeUK
)

func (m MarketType) String() string {
	switch m {
	case MarketTypeA:
		return "AStock"
	case MarketTypeHK:
		return "HK"
	case MarketTypeUS:
		return "UsStock"
	case MarketTypeUK:
		return "LSE"
	default:
		return ""
	}
}
