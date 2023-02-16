package movingsum

import (
	"container/list"
	"sync"
	"time"
)

type MovingSum struct {
	n     int
	queue *list.List
	sum   int
	sync.RWMutex
}

type timedEntry struct {
	t     time.Time
	value int
}

type MovingSumByTime struct {
	duration time.Duration
	queue    *list.List
	sum      int
	sync.Mutex
}

var timeNow = time.Now
var timeSince = time.Since

// This is moving sum with max entries
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

// This uses a queue to store entries within the duration
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
