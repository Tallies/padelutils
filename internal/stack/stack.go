package stack

import (
	"errors"
)

type Stack[T any] struct {
	items []T
}

func (stack *Stack[T]) Push(item T) {
	stack.items = append(stack.items, item)
}

func (stack *Stack[T]) Pop() (T, error) {
	if len(stack.items) == 0 {
		var zero T
		return zero, errors.New("The stack is empty.")
	}

	item := stack.items[len(stack.items)-1]
	stack.items = stack.items[:len(stack.items)-1]
	return item, nil
}

func (stack *Stack[T]) Peek() (T, bool) {
	if len(stack.items) == 0 {
		var zero T
		return zero, false
	}

	return stack.items[len(stack.items)-1], true
}

func (stack *Stack[T]) IsEmpty() bool {
	return len(stack.items) == 0
}
