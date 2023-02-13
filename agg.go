package movingsum

import "time"

type slot struct {
	sum   int
	count int
}

type AggregatedMovingSumByTime struct {
	duration time.Duration
	slots    []slot
	start    time.Time
	current  int
	total    slot
}

// This uses a ring buffer to store time slots
func NewAggregatedMovingSumByTime(duration time.Duration, n int) *AggregatedMovingSumByTime {
	return &AggregatedMovingSumByTime{
		slots:    make([]slot, n),
		duration: duration,
		start:    time.Now(),
	}
}

func (m *AggregatedMovingSumByTime) getCurrent() *slot {
	// slot is delta / (duration / slotCount)
	// Do multiplication first => (delta * slotCount) / duration
	delta := timeSince(m.start)
	newPos := int((delta * time.Duration(len(m.slots))) / m.duration)
	if newPos > m.current {
		if m.current+len(m.slots) <= newPos {
			// newPos beyond slot count. Clear all
			for i := range m.slots {
				m.slots[i].count = 0
				m.slots[i].sum = 0
			}
			m.total.count = 0
			m.total.sum = 0
		} else {
			// push current to total
			s := &m.slots[m.current%len(m.slots)]
			m.total.count += s.count
			m.total.sum += s.sum
			// pop others before newPos
			for i := m.current + 1; i <= newPos; i++ {
				s = &m.slots[i%len(m.slots)]
				m.total.count -= s.count
				m.total.sum -= s.sum
				s.count = 0
				s.sum = 0
			}
		}
		m.current = newPos
	}
	return &m.slots[m.current%len(m.slots)]
}

// TODO mutex!
// Add a new value to the moving sum
func (m *AggregatedMovingSumByTime) Add(value int) {
	s := m.getCurrent()
	s.count++
	s.sum += value
}

// Get returns the moving sum and the count of entries added
func (m *AggregatedMovingSumByTime) Get() (sum int, count int) {
	s := m.getCurrent()
	// include current slot
	return m.total.sum + s.sum, m.total.count + s.count
}
