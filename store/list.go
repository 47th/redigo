package store

type Node struct {
	val  string
	prev *Node
	next *Node
}

type LinkedList struct {
	head *Node
	tail *Node
	size int
}

func NewLinkedList(value string) *LinkedList {
	newNode := &Node{
		val: value,
	}

	newLinkedList := &LinkedList{
		head: newNode,
		tail: newNode,
		size: 1,
	}

	return newLinkedList
}

//think about the implementation of structs using sync.Pool

func (ll *LinkedList) rpush(value string) {
	node := &Node{
		val: value,
	}
	if ll.size == 0 {
		ll.head = node
		ll.tail = node
	} else {
		ll.tail.next = node
		node.prev = ll.tail
		ll.tail = node
	}
	ll.size++
}

func (ll *LinkedList) lpush(value string) {
	node := &Node{
		val: value,
	}
	if ll.size == 0 {
		ll.head = node
		ll.tail = node
	} else {
		ll.head.prev = node
		node.next = ll.head
		ll.head = node
	}
	ll.size++
}

func (ll *LinkedList) rpop() string {
	ret := ll.tail.val
	ll.tail = ll.tail.prev
	if ll.tail != nil {
		ll.tail.next = nil
	} else {
		ll.head = nil
	}
	ll.size--
	return ret
}

func (ll *LinkedList) lpop() string {
	ret := ll.head.val
	ll.head = ll.head.next
	if ll.head != nil {
		ll.head.prev = nil
	} else {
		ll.tail = nil
	}
	ll.size--
	return ret
}

func (ll *LinkedList) lrange(key string, start int, stop int) []string {
	// case 1 start +ve stop +ve -> if start < size -> iterate till min(stop, size)
	// case 2 start -ve stop +ve -> find start, start < 0, start becomes 0, if start
	// case 3 start +ve stop -ve -> normalize the stop,
	// case 4 start -ve stop -ve ->
	// combine these to -> if start is -ve, noarmalize it, then change it to 0 if it is still negative, then check if its less than size, then iterate.
	if start < 0 {
		start = ll.size + start
		if start < 0 {
			start = 0
		}
	}

	if stop < 0 {
		stop = ll.size + stop
	}

	if stop >= ll.size {
		stop = ll.size - 1
	}

	if start >= ll.size || stop < 0 {
		return []string{}
	}

	retArr := make([]string, stop-start+1)
	temp := ll.head
	for i, ind := 0, 0; i < ll.size; i++ {
		if i > stop {
			break
		}
		if i >= start {
			retArr[ind] = temp.val
			ind++
		}
		temp = temp.next
	}

	return retArr
}

func (ll *LinkedList) llen() int {
	return ll.size
}

// exposed functions

func (s *Store) LPUSH(key string, values []string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj := s.db[key]
	for _, value := range values {
		if valObj == nil {
			valObj = &Value{
				kind: List,
				list: NewLinkedList(value),
			}
			s.db[key] = valObj
		} else if valObj.kind == List {
			valObj.list.lpush(value)
		} else {
			return 0, false
		}
	}

	return valObj.list.size, true
}

func (s *Store) RPUSH(key string, values []string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj := s.db[key]
	for _, value := range values {
		if valObj == nil {
			valObj = &Value{
				kind: List,
				list: NewLinkedList(value),
			}
			s.db[key] = valObj
		} else if valObj.kind == List {
			valObj.list.rpush(value)
		} else {
			return 0, false
		}
	}

	return valObj.list.size, true
}

func (s *Store) LPOP(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj := s.db[key]
	if valObj == nil {
		return "NIL", false
	}
	if valObj.kind != List {
		return "WRONGTYPE", false
	}
	if valObj.list.size == 0 {
		return "NIL", false
	}

	ret := valObj.list.lpop()
	return ret, true
}

func (s *Store) RPOP(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj := s.db[key]
	if valObj == nil {
		return "NIL", false
	}
	if valObj.kind != List {
		return "WRONGTYPE", false
	}
	if valObj.list.size == 0 {
		return "NIL", false
	}

	ret := valObj.list.rpop()
	return ret, true
}

func (s *Store) LRANGE(key string, start int, stop int) ([]string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj := s.db[key]
	if valObj == nil {
		return []string{}, true
	}
	if valObj.kind != List {
		return []string{}, false
	}
	if valObj.list.size == 0 {
		return []string{}, true
	}

	values := valObj.list.lrange(key, start, stop)
	return values, true
}

func (s *Store) LLEN(key string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj := s.db[key]

	if valObj == nil {
		return 0, true
	}
	if valObj.kind != List {
		return 0, false
	}

	return valObj.list.llen(), true
}
