package service

var startingRange int64 = 1_000_000

func FetchNextIdRange() int64 {
	nextRange := startingRange
	startingRange += 100_000
	return nextRange
}
