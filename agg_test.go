package movingsum

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// With 10ns and 3 slots, based on the calculation:
//
//	slot 0 is 0,1,2,3ns
//	slot 1 is 4,5,6ns
//	slot 2 is 7,8,9ns
//	slot 3 is 10,11,12,13ns
//	...
func TestSlots(t *testing.T) {
	ms := NewAggregatedMovingSumByTime(time.Duration(10), 3)

	// Slot 0 for first 3 ns
	addAtTime(ms, 10, 0)
	assertByTime(t, ms, 10, 1, 0)
	addAtTime(ms, 10, 1)
	assertByTime(t, ms, 20, 2, 1)
	addAtTime(ms, 10, 2)
	assertByTime(t, ms, 30, 3, 2)
	addAtTime(ms, 10, 3)
	assertByTime(t, ms, 40, 4, 3)

	// Slot 1 at 4ns
	assertByTime(t, ms, 40, 4, 4)
	addAtTime(ms, 1, 4)
	assertByTime(t, ms, 41, 5, 4)

	// Slot 3 at 9ns
	assertByTime(t, ms, 41, 5, 9)

	// At 10ns. slot 0 is popped
	assertByTime(t, ms, 1, 1, 10)

	// Skip entire duration and add
	addAtTime(ms, 100, 20)
	assertByTime(t, ms, 100, 1, 20)

	// Skip entire duration
	assertByTime(t, ms, 0, 0, 30)

	// Add
	addAtTime(ms, 10, 31)
	addAtTime(ms, 10, 33)
	addAtTime(ms, 10, 35)
	assertByTime(t, ms, 30, 3, 39)

	// Skip many and add
	addAtTime(ms, 100, 1000)
	assertByTime(t, ms, 100, 1, 100)
}

func addAtTime(ms *AggregatedMovingSumByTime, value int, ns int) {
	timeSince = func(_ time.Time) time.Duration {
		return time.Duration(ns)
	}
	ms.Add(value)
}

// func assertByTime(t, ms *AggregatedMovingSumByTime, sum int, count int, ns int) bool {
// 	timeSince = func(_ time.Time) time.Duration {
// 		return time.Duration(ns)
// 	}
// 	s, c := ms.Get()
// 	return s == sum && c == count
// }

func assertByTime(t *testing.T, ms *AggregatedMovingSumByTime, sum int, count int, ns int) {
	timeSince = func(t time.Time) time.Duration {
		return time.Duration(ns)
	}
	s, c := ms.Get()
	require.Equal(t, sum, s)
	require.Equal(t, count, c)
}
