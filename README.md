# Description

An implementation of double ended queue via a linked list


# Auto Documentation

package deque // import "github.com/Booplicate/dequeue/pkg"


TYPES

```go
type Deque[T comparable] struct {
	// Has unexported fields.
}
```
    Double ended queue

```go
func NewDeque[T comparable](capacity int) *Deque[T]
```
    Creates a new deque with the given capacity. Capacity -1 creates a deque of
    unlimited size

```go
func NewDequeFromSeq[T comparable](seq iter.Seq[T], capacity int) *Deque[T]
```
    Creates a new deque from some iterator

```go
func NewUnlimitedDeque[T comparable]() *Deque[T]
```
    Creates a new deque with unlimited capacity

```go
func (self *Deque[T]) All() iter.Seq2[int, T]
```
    Returns iterator over deque values and their indices

```go
func (self *Deque[T]) Append(value T)
```
    Appends a new element to the right end of the deque. If capacity is not -1
    and the deque is full, an element is popped from the left end

```go
func (self *Deque[T]) AppendLeft(value T)
```
    Appends a new element to the left end of the deque. If capacity is not -1
    and the deque is full, an element is popped from the right end

```go
func (self *Deque[T]) Clear()
```
    Removes all elements from the deque

```go
func (self *Deque[T]) Copy() *Deque[T]
```
    Creates a shallow copy of the deque

```go
func (self *Deque[T]) Count(value T) int
```
    Returns the number of occurrences of the value given in the deque

```go
func (self *Deque[T]) GetCapacity() int
```
    Returns deque capacity

```go
func (self *Deque[T]) GetLen() int
```
    Returns current size of the deque

```go
func (self *Deque[T]) IsEmpty() bool
```
    Checks if the deque is empty

```go
func (self *Deque[T]) IsFull() bool
```
    Checks if the deque is full

```go
func (self *Deque[T]) IsUnlimited() bool
```
    Checks if the deque is of unlimited capacity

```go
func (self *Deque[T]) Peek(index int) T
```
    Returns an element at the given index

```go
func (self *Deque[T]) Rotate(n int)
```
    Rotates the deque by the given number of steps TODO: optimise: can be done
    without a loop

```go
func (self *Deque[T]) String() string
```

```go
func (self *Deque[T]) TryPeek(index int) (T, error)
```
    Returns an element at the given index or error if there's no element at such
    index

```go
func (self *Deque[T]) TryPop() (T, error)
```
    Removes an element from the right end and returns it. If the deque is empty,
    returns an error

```go
func (self *Deque[T]) TryPopLeft() (T, error)
```
    Removes an element from the left end and returns it. If the deque is empty,
    returns an error

```go
func (self *Deque[T]) Values() iter.Seq[T]
```
    Returns iterator over deque values

```go
type PeekError struct {
	Index int
}
```

```go
func (self *PeekError) Error() string
```

```go
type PopError struct{}
```

```go
func (self *PopError) Error() string
```

