package main

type Stacker interface {
	Push(v int)
	Pop() int
}

type stack struct {
	slice []int
}

func (s *stack) Push(v int) {
	s.slice = append(s.slice, v)
}

func (s *stack) Pop() int {
	if len(s.slice) == 0 || s.slice == nil {
		panic("стек пуст")
	}
	res := s.slice[len(s.slice)-1]
	s.slice = s.slice[:len(s.slice)-1]
	return res
}

func New() *stack {
	return &stack{slice: make([]int, 0)}
}
