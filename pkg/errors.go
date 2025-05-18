package deque

import "fmt"

type PopError struct{}

func (self *PopError) Error() string {
	return "deque: pop from empty queue"
}

type PeekError struct {
	Index int
}

func (self *PeekError) Error() string {
	return fmt.Sprintf("deque: index %d out of bounds", self.Index)
}
