package movingsum

import (
	"sync"
	"time"
)

type slot struct {
	sum   int
	count int
}

type AggregatedMovingSumByTime struct {
	duration   time.Duration
	slots      []slot
	currentPos int
	sum        int
	count      int
	sync.Mutex
}

// This uses a ring buffer to store time slots
func NewAggregatedMovingSumByTime(duration time.Duration, n int) *AggregatedMovingSumByTime {
	return &AggregatedMovingSumByTime{
		slots:    make([]slot, n),
		duration: duration,
	}
}

func (m *AggregatedMovingSumByTime) getCurrent() *slot {
	// slot is delta / (duration / slotCount)
	// Do multiplication first => (delta * slotCount) / duration
	delta := timeSince(time.Time{})
	newPos := int((delta * time.Duration(len(m.slots))) / m.duration)
	if newPos > m.currentPos {
		if m.currentPos+len(m.slots) <= newPos {
			// newPos beyond slot count. Clear all
			for i := range m.slots {
				m.slots[i].count = 0
				m.slots[i].sum = 0
			}
			m.count = 0
			m.sum = 0
		} else {
			// push current to total
			s := &m.slots[m.currentPos%len(m.slots)]
			m.count += s.count
			m.sum += s.sum
			// pop others before newPos
			for i := m.currentPos + 1; i <= newPos; i++ {
				s = &m.slots[i%len(m.slots)]
				m.count -= s.count
				m.sum -= s.sum
				s.count = 0
				s.sum = 0
			}
		}
		m.currentPos = newPos
	}
	return &m.slots[m.currentPos%len(m.slots)]
}

// Add a new value to the moving sum
func (m *AggregatedMovingSumByTime) Add(value int) {
	m.Lock()
	defer m.Unlock()
	s := m.getCurrent()
	s.count++
	s.sum += value
}

// TODO sum int64?
// Get returns the moving sum and the count of entries added
func (m *AggregatedMovingSumByTime) Get() (sum int, count int) {
	m.Lock()
	defer m.Unlock()
	s := m.getCurrent()
	// include current slot
	return m.sum + s.sum, m.count + s.count
}
