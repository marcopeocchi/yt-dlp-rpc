package main

import (
	"sync"
)

// In-Memory Key-Value Database
type MemoryDB struct {
	table map[int]*Process
	mu    sync.Mutex
}

func (m *MemoryDB) Init() {
	m.table = make(map[int]*Process)
}

func (m *MemoryDB) Get(id int) *Process {
	m.mu.Lock()
	res := m.table[id]
	m.mu.Unlock()
	return res
}

func (m *MemoryDB) Set(id int, process *Process) {
	m.mu.Lock()
	m.table[id] = process
	m.mu.Unlock()
}

func (m *MemoryDB) Keys() []int {
	m.mu.Lock()
	keys := make([]int, len(m.table))
	i := 0
	for k := range m.table {
		keys[i] = k
		i++
	}
	m.mu.Unlock()
	return keys
}

func (m *MemoryDB) Update(id int, progress Progress) {
	m.mu.Lock()
	if m.table[id] != nil {
		m.table[id].progress = progress
	}
	m.mu.Unlock()
}

func (m *MemoryDB) Delete(id int) {
	m.mu.Lock()
	delete(m.table, id)
	m.mu.Unlock()
}

func (m *MemoryDB) All() []Progress {
	running := make([]Progress, len(m.table))
	i := 0
	for _, v := range m.table {
		if v != nil {
			running[i] = v.progress
			i++
		}
	}
	return running
}
