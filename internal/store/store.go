package store

import "sync"

// Flag is a single feature-flag definition. It is the only thing the store
// ever holds: no user IDs from evaluate requests are persisted anywhere.
type Flag struct {
	Key            string `json:"key"`
	Description    string `json:"description"`
	Enabled        bool   `json:"enabled"`
	RolloutPercent int    `json:"rollout_percent"`
}

// Store is a thread-safe in-memory flag store.
type Store struct {
	mu    sync.Mutex
	flags map[string]Flag
}

// New returns an empty Store.
func New() *Store {
	return &Store{flags: make(map[string]Flag)}
}

// Create upserts the flag keyed by f.Key and returns the stored flag.
func (s *Store) Create(f Flag) Flag {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flags[f.Key] = f
	return f
}

// Get returns the flag for key and whether it exists.
func (s *Store) Get(key string) (Flag, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, ok := s.flags[key]
	return f, ok
}

// List returns every stored flag. It is never nil.
func (s *Store) List() []Flag {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Flag, 0, len(s.flags))
	for _, f := range s.flags {
		out = append(out, f)
	}
	return out
}

// Update overwrites the flag at key and returns it, or (zero, false) if the
// key does not exist. The key is preserved from the argument, not from f.
func (s *Store) Update(key string, f Flag) (Flag, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.flags[key]; !ok {
		return Flag{}, false
	}
	f.Key = key
	s.flags[key] = f
	return f, true
}

// Delete removes the flag at key and reports whether it existed.
func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.flags[key]; !ok {
		return false
	}
	delete(s.flags, key)
	return true
}
