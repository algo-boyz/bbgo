package volatility

import (
	v2 "github.com/c9s/bbgo/pkg/indicator/v2"
	"github.com/c9s/bbgo/pkg/types"
)

// OrderSide represents the side of an order: Buy (long) or Sell (short).
type OrderSide int

const (
	// Buy (long)
	Buy OrderSide = iota + 1

	// Sell (short)
	Sell
)

type ChandelierExitStream struct {
	// embedded struct
	*types.Float64Series

	atr *ATRStream
	min *v2.MinValueStream
	max *v2.MaxValueStream
}

// Chandelier Exit. It sets a trailing stop-loss based on the Average True Value (ATR).
//
// Chandelier Exit Long = 22-Period SMA High - ATR(22) * 3
// Chandelier Exit Short = 22-Period SMA Low + ATR(22) * 3
//
// Returns chandelierExitLong, chandelierExitShort
func ChandelierExit(source v2.KLineSubscription, os OrderSide, window int) *ChandelierExitStream {
	s := &ChandelierExitStream{
		atr: ATR2(source, window),
		min: v2.MinValue(v2.LowPrices(source), window),
		max: v2.MaxValue(v2.HighPrices(source), window)}

	source.AddSubscriber(func(v types.KLine) {
		high := s.max.Last(0)
		low := s.min.Last(0)
		atr := s.atr.Last(0)
		var chandelierExit float64
		if os == Buy {
			chandelierExit = high - atr*3
		} else {
			chandelierExit = low + atr*3
		}
		s.PushAndEmit(chandelierExit)
	})

	return s
}

func ChandelierExitDefault(source v2.KLineSubscription, os OrderSide) *ChandelierExitStream {
	return ChandelierExit(source, os, 22)
}

func (s *ChandelierExitStream) Truncate() {
	s.Slice = s.Slice.Truncate(5000)
}
