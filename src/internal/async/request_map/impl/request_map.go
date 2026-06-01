package request_map

import "sync"

type RequestAttemptsMap struct {
	mu   sync.Mutex
	data map[int]int
}

func NewRequestAttemptsMap() RequestAttemptsMap {
	return RequestAttemptsMap{data: make(map[int]int)}
}

func (m *RequestAttemptsMap) Put(id, attempts int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[id] = attempts
}

func (m *RequestAttemptsMap) Get(id int) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[id]
	return v, ok
}

func (m *RequestAttemptsMap) Delete(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
}
