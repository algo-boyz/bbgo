package pattern

import (
	v2 "github.com/c9s/bbgo/pkg/indicator/v2"
	"github.com/c9s/bbgo/pkg/types"
)

type DarkCloudStream struct {
	*types.Float64Series

	window int
}

// Dark Cloud Cover is a candlestick pattern that shows a shift in momentum to the downside
// following a price rise.
// The pattern is composed of a bearish candle that opens above but then closes below the midpoint of
// the prior bullish candle.
// Both candles should be relatively large, showing strong participation by traders and investors.
// When the pattern occurs with small candles it is typically less significant.
func DarkCloud(source v2.KLineSubscription) *DarkCloudStream {
	s := &DarkCloudStream{
		Float64Series: types.NewFloat64Series(),
		window:        2,
	}

	source.AddSubscriber(func(kLine types.KLine) {
		var (
			i      = source.Length()
			output = Neutral
		)
		if i < s.window {
			s.PushAndEmit(output)
			return
		}
		var (
			two             = source.Last(1)
			one             = source.Last(0)
			twoMidpoint     = (two.Close + two.Open) / 2
			isFirstBullish  = two.Close > two.Open
			isSecondBearish = one.Close < one.Open
			isDarkCloud     = one.Open > two.High &&
				one.Close < twoMidpoint && one.Close > two.Open
		)

		if isFirstBullish && isSecondBearish && isDarkCloud {
			output = Bear
		}

		s.PushAndEmit(output)

	})

	return s
}

func (s *DarkCloudStream) Truncate() {
	s.Slice = s.Slice.Truncate(MaxNumOfPattern)
}
