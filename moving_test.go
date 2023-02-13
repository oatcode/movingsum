package movingsum

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	ms := NewMovingSumByTime(time.Duration(10))

	msAddAtTime(ms, 10, 0)
	msAssertByTime(t, ms, 10, 1, 0)
	msAddAtTime(ms, 10, 1)
	msAssertByTime(t, ms, 20, 2, 1)
	msAddAtTime(ms, 10, 2)
	msAssertByTime(t, ms, 30, 3, 2)
	msAddAtTime(ms, 10, 3)
	msAssertByTime(t, ms, 40, 4, 3)

	// At 10ns. First entry is popped
	msAssertByTime(t, ms, 30, 3, 10)

	// Skip entire duration and add
	msAddAtTime(ms, 100, 20)
	msAssertByTime(t, ms, 100, 1, 20)

	// Skip entire duration
	msAssertByTime(t, ms, 0, 0, 30)

	// Add
	msAddAtTime(ms, 10, 31)
	msAddAtTime(ms, 10, 33)
	msAddAtTime(ms, 10, 35)
	msAssertByTime(t, ms, 30, 3, 39)

	// Skip many and add
	msAddAtTime(ms, 100, 1000)
	msAssertByTime(t, ms, 100, 1, 100)
}

func msAddAtTime(ms *MovingSumByTime, value int, ns int) {
	timeNow = func() time.Time {
		return time.Time{}.Add(time.Duration(ns))
	}
	timeSince = func(t time.Time) time.Duration {
		return time.Duration(ns - t.Nanosecond())
	}
	ms.Add(value)
}

func msAssertByTime(t *testing.T, ms *MovingSumByTime, sum int, count int, ns int) {
	timeSince = func(t time.Time) time.Duration {
		return time.Duration(ns - t.Nanosecond())
	}
	s, c := ms.Get()
	require.Equal(t, sum, s)
	require.Equal(t, count, c)
}
