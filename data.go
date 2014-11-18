package main

import (
	"fmt"
	"sync"
)

// Payload structure
type Payload struct {
	Id     string   `json:"id"`
	Params []string `json:"params"`
	Wid    string   `json:"wid"`
	Wiid   string   `json:"wiid"`
	Xid    string   `json:"xid"`
	Xuri   string   `json:"xuri"`
	Status string   `json:"status"` // "pending", "working", "blocked", "error:*", "finished"
}

func (p *Payload) String() string {
	return fmt.Sprintf("[%s] %s %s", p.Id, p.Wid, p.Wiid)
}

// Workflow structure
type Workflow struct {
	Id       string `json:"id"`
	Workflow []struct {
		Id     string   `json:"id"`
		Params []string `json:"params"`
		URL    string   `json:"url"`
	} `json:"workflow"`
}

// The Store interface defines methods to manipulate items.
type Store interface {
	Get(id string) *Payload
	GetAll() []*Payload
	Add(p *Payload) (string, error)
}

// Thread-safe in-memory map.
type payloadStore struct {
	sync.RWMutex
	m map[string]*Payload
}

// Store instance.
var store Store

func init() {
	store = &payloadStore{
		m: make(map[string]*Payload),
	}
	// Fill the store with some dummy data
	store.Add(&Payload{Wid: "wid-1", Wiid: "wiid-1", Xid: "xid-1", Xuri: "xuri-1"})
}

// GetAll returns all payloads from memory
func (store *payloadStore) GetAll() []*Payload {
	store.RLock()
	defer store.RUnlock()
	if len(store.m) == 0 {
		return nil
	}
	ar := make([]*Payload, len(store.m))
	i := 0
	for _, v := range store.m {
		ar[i] = v
		i++
	}
	return ar
}

// Get returns a single Payload identified by its id, or nil.
func (store *payloadStore) Get(id string) *Payload {
	store.RLock()
	defer store.RUnlock()
	return store.m[id]
}

// Add stores a new Payload and returns its newly generated id, or an error.
func (store *payloadStore) Add(p *Payload) (string, error) {
	store.Lock()
	defer store.Unlock()
	// Store it
	store.m[p.Id] = p
	return p.Id, nil
}

func (store *payloadStore) Update(p *Payload) error {
	store.Lock()
	defer store.Unlock()
	store.m[p.Id] = p
	return nil
}

func (store *payloadStore) Delete(id string) {
	store.Lock()
	defer store.Unlock()
	delete(store.m, id)
}
