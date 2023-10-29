package bollmaker

import "github.com/c9s/bbgo/pkg/indicator/v2/volatility"

type PriceTrend string

const (
	NeutralTrend PriceTrend = "neutral"
	UpTrend      PriceTrend = "upTrend"
	DownTrend    PriceTrend = "downTrend"
	UnknownTrend PriceTrend = "unknown"
)

func detectPriceTrend(inc *volatility.BollingerStream, price float64) PriceTrend {
	if inBetween(price, inc.DownBand.Last(0), inc.UpBand.Last(0)) {
		return NeutralTrend
	}

	if price < inc.DownBand.Last(0) {
		return DownTrend
	}

	if price > inc.UpBand.Last(0) {
		return UpTrend
	}

	return UnknownTrend
}
