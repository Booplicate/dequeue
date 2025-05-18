package deque

import (
	"fmt"
	"iter"
	"sync"
)

// Represents a node in deque
type node[T any] struct {
	value T
	next  *node[T]
	prev  *node[T]
}

func (self *node[T]) String() string {
	return fmt.Sprintf("Node{value:%v, next:%v}", self.value, self.next)
}

// Double ended queue
type Deque[T comparable] struct {
	head     *node[T]
	tail     *node[T]
	len      int
	capacity int
	mutex    sync.Mutex
}

func (self *Deque[T]) String() string {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return fmt.Sprintf("Deque{capacity:%v, head:%v, tail:%v}", self.capacity, self.head, self.tail)
}

// Creates a new deque with the given capacity.
// Capacity -1 creates a deque of unlimited size
func NewDeque[T comparable](capacity int) *Deque[T] {
	return &Deque[T]{nil, nil, 0, capacity, sync.Mutex{}}
}

// Creates a new deque with unlimited capacity
func NewUnlimitedDeque[T comparable]() *Deque[T] {
	return NewDeque[T](-1)
}

// Creates a new deque from some iterator
func NewDequeFromSeq[T comparable](seq iter.Seq[T], capacity int) *Deque[T] {
	q := NewDeque[T](capacity)
	for value := range seq {
		q.Append(value)
	}
	return q
}

// Returns current size of the deque
func (self *Deque[T]) GetLen() int {
	return self.len
}

// Checks if the deque is empty
func (self *Deque[T]) IsEmpty() bool {
	return self.GetLen() == 0
}

// Returns deque capacity
func (self *Deque[T]) GetCapacity() int {
	return self.capacity
}

// Checks if the deque is of unlimited capacity
func (self *Deque[T]) IsUnlimited() bool {
	return self.GetCapacity() < 0
}

// Checks if the deque is full
func (self *Deque[T]) IsFull() bool {
	return !self.IsUnlimited() && self.GetLen() >= self.GetCapacity()
}

// Checks if the deque is over its capacity
func (self *Deque[T]) isOverflowing() bool {
	return !self.IsUnlimited() && self.GetLen() > self.GetCapacity()
}

// Returns iterator over deque values
func (self *Deque[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		self.mutex.Lock()
		defer self.mutex.Unlock()

		for item := self.head; item != nil; item = item.next {
			if !yield(item.value) {
				return
			}
		}
	}
}

// Returns iterator over deque values and their indices
func (self *Deque[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for value := range self.Values() {
			if !yield(i, value) {
				return
			}
			i++
		}
	}
}

// Creates a shallow copy of the deque
func (self *Deque[T]) Copy() *Deque[T] {
	rv := NewDeque[T](self.capacity)
	rv.mutex.Lock()
	defer rv.mutex.Unlock()

	for value := range self.Values() {
		// Avoid mutex overhead
		rv.append(value)
	}

	return rv
}

// NOTE: assumes the mutex is acquired
func (self *Deque[T]) append(value T) {
	n := &node[T]{value, nil, nil}

	switch self.GetLen() {
	case 0:
		self.head = n
		self.tail = n
	case 1:
		n.prev = self.head
		self.head.next = n
		self.tail = n
	default:
		n.prev = self.tail
		self.tail.next = n
		self.tail = n
	}

	self.len++

	if self.isOverflowing() {
		self.tryPopLeft()
	}

}

// Appends a new element to the right end of the deque.
// If capacity is not -1 and the deque is full, an element is popped from the left end
func (self *Deque[T]) Append(value T) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.append(value)
}

// NOTE: assumes the mutex is acquired
func (self *Deque[T]) appendLeft(value T) {
	n := &node[T]{value, nil, nil}

	switch self.GetLen() {
	case 0:
		self.head = n
		self.tail = n
	default:
		n.next = self.head
		self.head.prev = n
		self.head = n
	}

	self.len++

	if self.isOverflowing() {
		self.tryPop()
	}
}

// Appends a new element to the left end of the deque.
// If capacity is not -1 and the deque is full, an element is popped from the right end
func (self *Deque[T]) AppendLeft(value T) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.appendLeft(value)
}

// NOTE: assumes the mutex is acquired
func (self *Deque[T]) tryPop() (T, error) {
	var value T

	switch self.GetLen() {
	case 0:
		return value, &PopError{}
	case 1:
		self.head = nil
		self.tail = nil
	default:
		oldTail := self.tail
		newTail := oldTail.prev
		newTail.next = nil
		oldTail.next = nil
		oldTail.prev = nil
		self.tail.prev = nil
		self.tail = newTail
	}
	self.len--

	return value, nil
}

// Removes an element from the right end and returns it.
// If the deque is empty, returns an error
func (self *Deque[T]) TryPop() (T, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return self.tryPop()
}

// NOTE: assumes the mutes is acquired
func (self *Deque[T]) tryPopLeft() (T, error) {
	var value T

	switch self.GetLen() {
	case 0:
		return value, &PopError{}
	case 1:
		self.head = nil
		self.tail = nil
	default:
		newHead := self.head.next
		newHead.prev = nil
		self.head = newHead
	}
	self.len--

	return value, nil
}

// Removes an element from the left end and returns it.
// If the deque is empty, returns an error
func (self *Deque[T]) TryPopLeft() (T, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return self.tryPopLeft()
}

// Removes all elements from the deque
func (self *Deque[T]) Clear() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.head = nil
	self.tail = nil
	self.len = 0
}

// Returns the number of occurrences of the value given in the deque
func (self *Deque[T]) Count(value T) int {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	i := 0
	n := self.head
	for n != nil {
		if n.value == value {
			i++
		}
		n = n.next
	}
	return i
}

// Returns an element at the given index or error if there's no element at such index
func (self *Deque[T]) TryPeek(index int) (T, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	var value T

	if index < 0 || index >= self.GetLen() {
		return value, &PeekError{index}
	}

	if index < self.GetLen()/2 {
		// Start from the head in the index is in the first half
		n := self.head
		for range index {
			n = n.next
		}
		value = n.value
	} else {
		// Otherwise start from the tail, this gives us O(1) for head/tail lookup
		n := self.tail
		for range self.GetLen() - 1 - index {
			n = n.prev
		}
		value = n.value
	}

	return value, nil
}

// Returns an element at the given index
func (self *Deque[T]) Peek(index int) T {
	value, err := self.TryPeek(index)
	if err != nil {
		panic(err)
	}
	return value
}

// Rotates the deque to the right, unsafe
func (self *Deque[T]) rotateRight() {
	tail := self.tail
	tail.next = self.head
	self.tail = tail.prev
	self.tail.next = nil
	tail.prev = nil
	self.head = tail
}

// Rotates the deque to the right, unsafe
func (self *Deque[T]) rotateLeft() {
	head := self.head
	self.head = head.next
	head.next = nil
	tail := self.tail
	tail.next = head
	self.tail = head
}

// Rotates the deque by the given number of steps
// TODO: optimise: can be done without a loop
func (self *Deque[T]) Rotate(n int) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.GetLen() < 2 {
		return
	}

	var doRotation func()
	if n >= 0 {
		doRotation = self.rotateRight
	} else {
		doRotation = self.rotateLeft
	}

	if n < 0 {
		n = -n
	}
	for range n {
		doRotation()
	}
}
