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
	ll.tail.next = node
	node.prev = ll.tail
	ll.tail = node
	ll.size++
}

func (ll *LinkedList) lpush(value string) {
	node := &Node{
		val: value,
	}
	ll.head.prev = node
	node.next = ll.head
	ll.head = node
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

// exposed functions

func (s *Store) LPUSH(key string, values []string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj := s.db[key]
	for _, value := range values {
		if valObj == nil || valObj.list.size == 0 {
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
		if valObj == nil || valObj.list.size == 0 {
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
	if valObj == nil || valObj.kind != List || valObj.list.size == 0 {
		return "", false
	}

	ret := valObj.list.lpop()
	return ret, true
}

func (s *Store) RPOP(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	valObj := s.db[key]
	if valObj == nil || valObj.kind != List || valObj.list.size == 0 {
		return "", false
	}

	ret := valObj.list.rpop()
	return ret, true
}
