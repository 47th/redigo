package store

import (
	"strconv"
	"sync"
	"time"
)

type Value struct {
	val    string
	expiry time.Time
}

type Store struct {
	mu sync.Mutex
	db map[string]*Value
}

func NewStore() *Store {
	return &Store{
		db: make(map[string]*Value),
	}
}

func (s *Store) set(key string, value string, expiry time.Time) {
	s.db[key] = &Value{
		val:    value,
		expiry: expiry,
	}
}

func (s *Store) get(key string) (*Value, bool) {
	value, ok := s.db[key]
	if !ok {
		return value, ok
	}

	if !(value.expiry.IsZero()) && time.Now().After(value.expiry) {
		delete(s.db, key)
		return s.get(key)
	}
	return value, ok
}

func (s *Store) SET(key string, value string, expiry time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.set(key, value, expiry)
}

func (s *Store) GET(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj, ok := s.get(key)
	if !ok {
		return "", ok
	}
	return valObj.val, ok
}

func (s *Store) INCR(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj, ok := s.get(key)
	if !ok {
		s.set(key, "1", time.Time{})
		return "1", true
	}

	newVal, ok := incrStr(valObj.val)
	if !ok {
		return "", ok
	}

	valObj.val = newVal
	return newVal, ok
}

// Problem : redis increment will retrieve the value using the get command,
// increment the value and then store the new value using the set command
// but the set requires you to set an expiry but we dont have an idea about the expiry
// so we need to only set the value and not the expiry if possible

// func isExpired(val Value) bool {
// 	if !(val.expiry.IsZero()) && time.Now().After(val.expiry) {
// 		return true
// 	}
// 	return false
// }

func incrStr(str string) (string, bool) {
	if str == "" {
		return "1", true
	}
	n, err := strconv.Atoi(str)
	if err != nil {
		return "", false
	}

	newVal := strconv.Itoa(n + 1)
	return newVal, true
}
