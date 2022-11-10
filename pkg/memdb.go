package pkg

import (
	"sync"

	"github.com/google/uuid"
)

// In-Memory Fast Key-Value Storage/DB
type MemoryDB struct {
	table map[string]*Process
	mu    sync.Mutex
}

// Inits the db with an empty map of string->Process pointer
func (m *MemoryDB) Init() {
	m.table = make(map[string]*Process)
}

// Get a process pointer given its id
func (m *MemoryDB) Get(id string) *Process {
	m.mu.Lock()
	res := m.table[id]
	m.mu.Unlock()
	return res
}

// Store a pointer of a process and return its id
func (m *MemoryDB) Set(process *Process) string {
	id := uuid.Must(uuid.NewRandom()).String()
	m.mu.Lock()
	m.table[id] = process
	m.mu.Unlock()
	return id
}

// Update a process progress, given the process id
func (m *MemoryDB) Update(id string, info DownloadInfo) {
	m.mu.Lock()
	if m.table[id] != nil {
		m.table[id].Info = info
	}
	m.mu.Unlock()
}

func (m *MemoryDB) UpdateProgress(id string, progress DownloadProgress) {
	m.mu.Lock()
	if m.table[id] != nil {
		m.table[id].Progress = progress
	}
	m.mu.Unlock()
}

// Removes a process progress, given the process id
func (m *MemoryDB) Delete(id string) {
	m.mu.Lock()
	delete(m.table, id)
	m.mu.Unlock()
}

// Returns a slice of all currently stored processes id
func (m *MemoryDB) Keys() []string {
	m.mu.Lock()
	keys := make([]string, len(m.table))
	i := 0
	for k := range m.table {
		keys[i] = k
		i++
	}
	m.mu.Unlock()
	return keys
}

// Returns a slice of all currently stored processes progess
func (m *MemoryDB) All() []ProcessResponse {
	running := make([]ProcessResponse, len(m.table))
	i := 0
	for k, v := range m.table {
		if v != nil {
			running[i] = ProcessResponse{
				Id:       k,
				Info:     v.Info,
				Progress: v.Progress,
			}
			i++
		}
	}
	return running
}
