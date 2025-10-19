package generator

import "fmt"

type RangeTracker struct {
	Start   uint64
	End     uint64
	Current uint64
}

var currentRange *RangeTracker

func GetNextId() uint64 {
	// 1. If current range is empty or consumed completely, get new range from DB.
	// 2. Save current range in local memory
	// 3. Get next Id from the range
	if currentRange == nil || currentRange.Current == currentRange.End {
		currentRange = &RangeTracker{
			Start:   1_000_000,
			End:     1_100_000,
			Current: 1_000_000,
		}
	}

	fmt.Printf("Current Range: %v\n", currentRange)

	nextId := currentRange.Current
	currentRange.Current += 1
	return nextId
}
