package store

import "time"

type Val struct {
	Value  string
	Expiry time.Time
}

type Store struct {
	db map[string]Val
}

func NewStore() *Store {
	return &Store{
		db: make(map[string]Val),
	}
}

func (s *Store) Set(key string, value string, expiry time.Time) {
	s.db[key] = Val{
		Value:  value,
		Expiry: expiry,
	}
}

func (s *Store) Get(key string) (string, bool) {
	expiry := s.db[key].Expiry
	if !expiry.IsZero() && time.Now().After(expiry) {
		delete(s.db, key)
		return "", false
	}
	return s.db[key].Value, true
}
