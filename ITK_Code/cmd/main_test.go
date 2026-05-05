package main

import (
	"testing"
)

func TestStack_PushPop(t *testing.T) {
	s := New()

	s.Push(1)
	s.Push(2)
	s.Push(3)

	tests := []struct {
		expected int
	}{
		{3},
		{2},
		{1},
	}

	for _, tc := range tests {
		got := s.Pop()
		if got != tc.expected {
			t.Errorf("Pop() = %d; ожидалось %d", got, tc.expected)
		}
	}
}

func TestStack_PopEmpty(t *testing.T) {
	s := New()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Ожидалась паника при попытке извлечь элемент из пустого стека")
		}
	}()

	s.Pop()
}
