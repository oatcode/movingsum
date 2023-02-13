package movingsum

import (
	"container/list"
	"time"
)

type entry struct {
	t     time.Time
	value int
}
type MovingSumByTime struct {
	duration time.Duration
	queue    *list.List
	total    int
}

var timeNow = time.Now
var timeSince = time.Since

// This uses a ring buffer to store time slots
func NewMovingSumByTime(duration time.Duration) *MovingSumByTime {
	return &MovingSumByTime{
		queue:    list.New(),
		duration: duration,
	}
}

func (m *MovingSumByTime) expire() {
	// Check backwards on the queue
	elem := m.queue.Back()
	for elem != nil {
		e := elem.Value.(entry)
		if timeSince(e.t) < m.duration {
			// Nothing before this is expired
			break
		}
		m.total -= e.value
		delete := elem
		elem = delete.Prev()
		m.queue.Remove(delete)
	}
}

// Add a new value to the moving sum
func (m *MovingSumByTime) Add(value int) {
	m.expire()
	m.queue.PushFront(entry{t: timeNow(), value: value})
	m.total += value
}

// Get returns the moving sum and the count of entries added
func (m *MovingSumByTime) Get() (sum int, count int) {
	m.expire()
	return m.total, m.queue.Len()
}
