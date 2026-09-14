package datastore

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Item represents a simple resource in the REST API.
type Item struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// Store is a simple in-memory store for demonstration.
// In production, replace with a real database (Postgres, Redis, etc.).
type Store struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]Item),
	}
}

func (s *Store) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Item, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, item)
	}
	return result
}

func (s *Store) Get(id string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	return item, ok
}

func (s *Store) Create(title string) Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	item := Item{
		ID:        uuid.New().String(),
		Title:     title,
		Completed: false,
		CreatedAt: time.Now().UTC(),
	}
	s.items[item.ID] = item
	return item
}

func (s *Store) Update(id string, title string, completed bool) (Item, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok {
		return Item{}, false
	}
	item.Title = title
	item.Completed = completed
	s.items[id] = item
	return item, true
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false
	}
	delete(s.items, id)
	return true
}
