package lib

import (
	"fmt"
	"math"
)

type TimeRange struct {
	Start float64
	End   float64
}

func (r TimeRange) String() string {
	return fmt.Sprintf("[%f-%f]", r.Start, r.End)
}

type TimeRangeState struct {
	TimeRange
	EndState bool
}

func (r TimeRangeState) String() string {
	return fmt.Sprintf("[%f-%f] -> %t", r.Start, r.End, r.EndState)
}

// mark the boundaries between sim.Similar states - when the state changes, e.g.
// input: 0t, 1t, 2f, 3f, 4f, 5t, 6t, 7t, 8f, 9f, 10t, 11t
// output: [1,2], [4,5], [7,8], [9,10]
func WalkBoundaryRanges(sims []SimilarFrame) []TimeRangeState {
	var ranges []TimeRangeState
	var last bool
	for i, sim := range sims {
		if i == 0 {
			last = sim.Similar
			continue
		}
		if last != sim.Similar {
			ranges = append(ranges, TimeRangeState{
				TimeRange: TimeRange{
					Start: sims[i-1].Time,
					End:   sim.Time,
				},
				EndState: sim.Similar,
			})
		}
		last = sim.Similar
	}
	return ranges
}

// Convert boundary time ranges to "keeper" ranges, i.e. the in-between time ranges that contain non-blank video.
// input: [6.1-6.2] -> false, [120.8-120.9] -> true, [350.4-350.5] -> false, [600.9-610.0] true
// output: [6.2-120.8], [350.5-600.9]
// input: [6.1-6.2] -> false
// output: [6.2-positive infinity]
func TimeRangesToKeepers(ranges []TimeRangeState) []TimeRange {
	var keepers []TimeRange
	for i, r := range ranges {
		if i == 0 {
			continue
		}
		if r.EndState {
			keepers = append(keepers, TimeRange{
				Start: ranges[i-1].End,
				End:   r.Start,
			})
		}
	}
	// open-ended range
	if !ranges[len(ranges)-1].EndState {
		keepers = append(keepers, TimeRange{
			Start: ranges[len(ranges)-1].End,
			End:   math.Inf(1),
		})
	}
	return keepers
}
