package store

import (
	"strconv"
	"sync"
	"time"
)

type Type string

const (
	String  = "STR"
	Integer = "INT"
	Array   = "ARR"
	List    = "LST"
)

type Value struct {
	kind   Type
	str    string
	list   *LinkedList
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
		kind:   String,
		str:    value,
		expiry: expiry,
	}
}

func (s *Store) get(key string) (*Value, bool) {
	value, ok := s.db[key]
	if !ok || value.kind != String {
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
	return valObj.str, ok
}

func (s *Store) INCR(key string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj, ok := s.get(key)
	if !ok {
		s.set(key, "1", time.Time{})
		return 1, true
	}

	newVal, ok := incrStr(valObj.str)
	if !ok {
		return 0, ok
	}

	valObj.str = strconv.Itoa(newVal)
	return newVal, ok
}

func (s *Store) DECR(key string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj, ok := s.get(key)
	if !ok {
		s.set(key, "-1", time.Time{})
		return -1, true
	}

	newVal, ok := decrStr(valObj.str)
	if !ok {
		return 0, ok
	}

	valObj.str = strconv.Itoa(newVal)
	return newVal, ok
}

// Problem : redis increment will retrieve the value using the get command,
// increment the value and then store the new value using the set command
// but the set requires you to set an expiry but we dont have an idea about the expiry
// so we need to only set the value and not the expiry if possible

func incrStr(str string) (int, bool) {
	if str == "" {
		return 1, true
	}
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, false
	}

	// newVal := strconv.Itoa(n + 1)
	return n + 1, true
}

func decrStr(str string) (int, bool) {
	if str == "" {
		return -1, true
	}
	n, err := strconv.Atoi(str)
	if err != nil {
		return 0, false
	}

	// newVal := strconv.Itoa(n + 1)
	return n - 1, true
}
