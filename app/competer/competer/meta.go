package competer

import (
	"sync"
)

type CompeterMeta struct {
	count int

	mu sync.Mutex
}

func NewCompeterMeta() *CompeterMeta {
	return &CompeterMeta{
		count: 0,
	}
}

func (m *CompeterMeta) GetCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.count
}

func (m *CompeterMeta) UpdateCount(count int) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.count += count

	return m.count
}
