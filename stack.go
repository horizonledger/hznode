package main

type Stack struct {
	items []interface{} // Using interface{} means the stack can hold elements of any type
}

// Push adds an item to the stack
func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

// Pop removes an item from the stack
func (s *Stack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}

	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1] // This reslices the slice to remove the last element

	return item
}

func (s *Stack) Len() interface{} {
	return len(s.items)
}
