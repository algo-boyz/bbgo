package volume

import (
	"github.com/c9s/bbgo/pkg/fixedpoint"
	v2 "github.com/c9s/bbgo/pkg/indicator/v2"
	"github.com/c9s/bbgo/pkg/types"
)

// Accumulation/Distribution Indicator (A/D). Cumulative indicator
// that uses volume and price to assess whether a stock is
// being accumulated or distributed.
//
// MFM = ((Closing - Low) - (High - Closing)) / (High - Low)
// MFV = MFM * Period Volume
// AD = Previous AD + CMFV
type AccumulationDistributionStream struct {
	*types.Float64Series
}

func AccumulationDistribution(source v2.KLineSubscription) *AccumulationDistributionStream {
	s := &AccumulationDistributionStream{
		Float64Series: types.NewFloat64Series(),
	}

	source.AddSubscriber(func(v types.KLine) {
		var (
			i      = s.Slice.Length()
			output = fixedpoint.NewFromInt(0)
		)

		if i > 0 {
			output = fixedpoint.NewFromFloat(s.Slice.Last(0))
		}

		output += v.Volume * ((v.Close - v.Low) - (v.High-v.Close)/(v.High-v.Low))

		s.PushAndEmit(output.Float64())
	})

	return s
}
