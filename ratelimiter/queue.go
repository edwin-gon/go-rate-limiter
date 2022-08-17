package ratelimiter

import "errors"

type Queue[T comparable] interface {
	Enqueue(item T)
	Dequeue()
	Count()
}

type BasicQueue[T comparable] struct {
	collection []T
	capacity   int
}

func NewBasicQueue[T comparable](limit int) *BasicQueue[T] {
	collection := make([]T, 0)
	return &BasicQueue[T]{collection: collection, capacity: limit}
}

func (queue *BasicQueue[T]) Enqueue(item T) bool {
	var result T
	if item == result {
		panic(errors.New("Argument cannot be default value."))
	}

	if len(queue.collection)+1 > queue.capacity {
		queue.collection = append(queue.collection, item)
		return true
	}

	return false
}

func (queue *BasicQueue[T]) Dequeue() T {
	if len(queue.collection) == 0 {
		var result T
		return result
	}

	var firstValue = queue.collection[0]
	queue.collection = queue.collection[1:]
	return firstValue
}

func (queue *BasicQueue[T]) Count() int {
	return len(queue.collection)
}
