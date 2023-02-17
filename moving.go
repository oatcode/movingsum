// Package movingsum provides utilities to compute moving sum
package movingsum

import (
	"container/list"
	"sync"
	"time"
)

// MovingSum calculates moving sum with a fixed queue size
type MovingSum struct {
	n     int
	queue *list.List
	sum   int
	sync.RWMutex
}

// used in MovingSumByTime
type timedEntry struct {
	t     time.Time
	value int
}

// MovingSumByTime calculates moving sum for a duration
type MovingSumByTime struct {
	duration time.Duration
	queue    *list.List
	sum      int
	sync.Mutex
}

// for UT purpose
var timeNow = time.Now
var timeSince = time.Since

// NewMovingSum creates moving sum with a max of n entries
func NewMovingSum(n int) *MovingSum {
	return &MovingSum{
		n:     n,
		queue: list.New(),
	}
}

// Add a new value to the moving sum
func (m *MovingSum) Add(value int) {
	m.Lock()
	defer m.Unlock()
	if m.queue.Len() >= m.n {
		e := m.queue.Back()
		m.sum -= e.Value.(int)
		m.queue.Remove(e)
	}
	m.queue.PushFront(value)
	m.sum += value
}

// Get returns the moving sum and the count of entries added
func (m *MovingSum) Get() (sum int, count int) {
	m.RLock()
	defer m.RUnlock()
	return m.sum, m.queue.Len()
}

// NewMovingSumByTime creates moving sum based on duration
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
		e := elem.Value.(timedEntry)
		if timeSince(e.t) < m.duration {
			// Nothing before this is expired
			break
		}
		m.sum -= e.value
		delete := elem
		elem = delete.Prev()
		m.queue.Remove(delete)
	}
}

// Add a new value to the moving sum
func (m *MovingSumByTime) Add(value int) {
	m.Lock()
	defer m.Unlock()
	m.expire()
	m.queue.PushFront(timedEntry{t: timeNow(), value: value})
	m.sum += value
}

// Get returns the moving sum and the count of entries added
func (m *MovingSumByTime) Get() (sum int, count int) {
	m.Lock()
	defer m.Unlock()
	m.expire()
	return m.sum, m.queue.Len()
}
