package ratelimiter

import "errors"

type Queue[T comparable] interface {
	Enqueue(item T)
	Dequeue()
}

type BasicQueue[T comparable] struct {
	collection []T
	capacity   int
}

func (queue *BasicQueue[T]) NewBasicQueue(limit int) *BasicQueue[T] {
	collection := make([]T, 0)
	return &BasicQueue[T]{collection: collection, capacity: limit}
}

func (queue *BasicQueue[T]) Enqueue(item T) {
	var result T
	if item == result {
		panic(errors.New("Argument cannot be default value."))
	}

	if len(queue.collection)+1 > queue.capacity {
		queue.collection = append(queue.collection, item)
	}
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
