package deque

import (
	"fmt"
	"slices"
	"testing"
)

func TestNewDeque(t *testing.T) {
	testCases := []struct{ capacity int }{
		{-1},
		{0},
		{3},
	}
	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("Capacity/%d", tc.capacity),
			func(t *testing.T) {
				q := NewDeque[int](tc.capacity)

				if q.GetCapacity() != tc.capacity {
					t.Errorf("Expected capacity to be %d, got %d", tc.capacity, q.GetCapacity())
				}
				if q.GetLen() != 0 {
					t.Errorf("Expected len to be 0, got %d", q.GetLen())
				}
				if q.head != nil {
					t.Errorf("Expected head to be nil, got %v", q.head)
				}
				if q.tail != nil {
					t.Errorf("Expected tail to be nil, got %v", q.tail)
				}
			},
		)
	}
}

func TestNewUnlimitedDeque(t *testing.T) {
	q := NewUnlimitedDeque[int]()

	if q.GetCapacity() != -1 {
		t.Errorf("Expected capacity to be %d, got %d", -1, q.GetCapacity())
	}
	if q.GetLen() != 0 {
		t.Errorf("Expected len to be 0, got %d", q.GetLen())
	}
	if q.head != nil {
		t.Errorf("Expected head to be nil, got %v", q.head)
	}
	if q.tail != nil {
		t.Errorf("Expected tail to be nil, got %v", q.tail)
	}
}

func TestNewDequeFromSec(t *testing.T) {
	testCases := []struct {
		slice    []int
		capacity int
	}{
		{[]int{0, 1, 2, 3, 4, 5}, 10},
		{[]int{}, 10},
		{[]int{77, -2, 9}, 10},
		{[]int{7, 7, 7, 7, 7, 7, 7}, 3},
		{[]int{}, 0},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, -1},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("InputSlice/%v", tc.slice),
			func(t *testing.T) {
				q := NewDequeFromSeq(slices.Values(tc.slice), tc.capacity)

				if q.GetCapacity() != tc.capacity {
					t.Errorf("Expected capacity to be %d, got %d", tc.capacity, q.GetCapacity())
				}

				var expectedLen int
				if q.IsUnlimited() {
					expectedLen = len(tc.slice)
				} else {
					expectedLen = min(len(tc.slice), q.GetCapacity())
				}
				if q.GetLen() != expectedLen {
					t.Errorf("Expected len to be %d, got %d", expectedLen, q.GetLen())
				}

				for i, v := range q.All() {
					if tc.slice[i] != v {
						t.Errorf("Value at index %d do not match, expected %d, got %d", i, tc.slice[i], v)
					}
				}
			},
		)
	}
}

func TestGetCapacity(t *testing.T) {
	testCases := []struct{ capacity int }{
		{-1},
		{0},
		{1},
		{2},
		{5},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("Capacity/%d", tc.capacity),
			func(t *testing.T) {
				q := NewDeque[int](tc.capacity)
				if q.GetCapacity() != tc.capacity {
					t.Errorf("Expected capacity %d, but GetCapacity() returned %d", tc.capacity, q.GetCapacity())
				}
			},
		)
	}
}

func TestDefaultLen(t *testing.T) {
	testCases := []struct{ capacity, expectedLen int }{
		{-1, 0},
		{0, 0},
		{1, 0},
		{2, 0},
		{5, 0},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("Capacity/%d/ExpectedLen/%d", tc.capacity, tc.expectedLen),
			func(t *testing.T) {
				q := NewDeque[int](tc.capacity)
				if q.GetLen() != 0 {
					t.Errorf("Expected initial len 0, but GetLen() returned %d", q.GetLen())
				}
			},
		)
	}
}

func TestIsFull(t *testing.T) {
	testCases := []struct {
		slice    []int
		capacity int
	}{
		{[]int{}, 1},
		{[]int{0}, 1},
		{[]int{0}, -1},
		{[]int{0, 1, 2, 3, 4}, 3},
		{[]int{0, 1, 2, 3, 4}, 10},
	}
	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("Slice/%v/Capacity/%d", tc.slice, tc.capacity),
			func(t *testing.T) {
				q := NewDequeFromSeq(slices.Values(tc.slice), tc.capacity)

				if tc.capacity < 0 && q.IsFull() {
					t.Errorf("Unlimited deque cannot be full")
				}

				if tc.capacity > -1 && len(tc.slice) >= tc.capacity && !q.IsFull() {
					t.Errorf(
						"Slice size (%d) is bigger than deque capacity (%d), expected deque to be full, but it is NOT",
						len(tc.slice),
						tc.capacity,
					)
				}

				if tc.capacity > -1 && len(tc.slice) < tc.capacity && q.IsFull() {
					t.Errorf(
						"Slice size (%d) is smaller than deque capacity (%d), expected deque to NOT be full, but it IS",
						len(tc.slice),
						tc.capacity,
					)
				}
			},
		)
	}
}

func TestIter(t *testing.T) {
	const TOTAL_ITEMS int = 4

	q := NewUnlimitedDeque[int]()
	for i := range TOTAL_ITEMS {
		q.Append(i)
	}

	counter := 0
	for i, v := range q.All() {
		if counter >= TOTAL_ITEMS {
			t.Errorf("Expected up to %d elements, got %d", TOTAL_ITEMS-1, counter)
		}
		if counter != v {
			t.Errorf("Expected next element be %d, got %d", counter, v)
		}
		if v != i {
			t.Errorf("Expected next index to be %d, got %d", v, i)
		}
		counter++

	}
}

func TestAppend(t *testing.T) {
	testValues := make([]int, 16)
	for v := range len(testValues) {
		testValues[v] = v
	}
	testCases := []struct{ capacity int }{
		{-1},
		{0},
		{1},
		{2},
		{5},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("Capacity/%d", tc.capacity),
			func(t *testing.T) {
				q := NewDeque[int](tc.capacity)

				for i, v := range testValues {
					q.Append(v)
					if q.GetCapacity() != tc.capacity {
						t.Errorf("Capacity has changed, should be %d, but it's %d", tc.capacity, q.GetCapacity())
					}

					var expectedLen int
					if tc.capacity < 0 {
						expectedLen = i + 1
					} else {
						expectedLen = min(i+1, tc.capacity)
					}
					if q.GetLen() != expectedLen {
						t.Errorf("Length is invalid, expected %d, got %d", expectedLen, q.GetLen())
					}

					if tc.capacity != 0 && q.tail.value != v {
						t.Errorf("Expected tail to be equal %d, got %d", v, q.tail.value)
					} else if tc.capacity == 0 && q.tail != nil {
						t.Errorf("Expected tail to be equal to nil, got %d", q.tail.value)
					}
				}
			},
		)
	}
}

func TestAppendLeft(t *testing.T) {
	testValues := make([]int, 16)
	for v := range len(testValues) {
		testValues[v] = v
	}
	testCases := []struct{ capacity int }{
		{-1},
		{0},
		{1},
		{2},
		{5},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("Capacity/%d", tc.capacity),
			func(t *testing.T) {
				q := NewDeque[int](tc.capacity)

				for i, v := range testValues {
					q.AppendLeft(v)
					if q.GetCapacity() != tc.capacity {
						t.Fatalf("Capacity has changed, should be %d, but it's %d", tc.capacity, q.GetCapacity())
					}

					var expectedLen int
					if tc.capacity < 0 {
						expectedLen = i + 1
					} else {
						expectedLen = min(i+1, tc.capacity)
					}
					if q.GetLen() != expectedLen {
						t.Fatalf("Length is invalid, expected %d, got %d", expectedLen, q.GetLen())
					}

					if tc.capacity != 0 && q.head.value != v {
						t.Fatalf("Expected tail to be equal %d, got %d", v, q.head.value)
					} else if tc.capacity == 0 && q.head != nil {
						t.Fatalf("Expected tail to be equal to nil, got %d", q.head.value)
					}
				}
			},
		)
	}
}

func TestTryPop(t *testing.T) {
	q := NewDeque[float32](8)
	_, err := q.TryPop()
	if err == nil {
		t.Errorf("Expected Pop() to error out")
	}
	q.AppendLeft(2.0)
	_, err = q.TryPop()
	if err != nil {
		t.Errorf("Expected Pop() to remove element, but got an error")
	}
}

func TestTryPopLeft(t *testing.T) {
	q := NewDeque[float32](8)
	_, err := q.TryPopLeft()
	if err == nil {
		t.Errorf("Expected PopLeft() to error out")
	}
	q.Append(8.8)
	_, err = q.TryPopLeft()
	if err != nil {
		t.Errorf("Expected PopLeft() to remove element, but got an error")
	}
}

func TestTryPeek(t *testing.T) {
	const SIZE int = 10

	q := NewDeque[int](SIZE)
	for i := range SIZE {
		q.Append(i)
		v, err := q.TryPeek(0)
		if err != nil {
			t.Errorf("Expected to get 0, got: err=%v", err)
		} else if v != 0 {
			t.Errorf("Expected PeekN() to return 0, got %d", v)
		}
	}

	// NOTE: element index is equal to its value, so we only need idx
	indecies := []int{1, q.GetLen() - 1, 5, 8}
	for _, idx := range indecies {
		t.Run(
			fmt.Sprintf("index/%d", idx),
			func(t *testing.T) {
				v, err := q.TryPeek(idx)
				if err != nil {
					t.Errorf("Expected to get %d, got: err=%v", idx, err)
				} else if v != idx {
					t.Errorf("Expected PeekN() to return %d, got %d", idx, v)
				}
			},
		)
	}
}

func TestTryPeekOutOutBounds(t *testing.T) {
	t.Run(
		"peek-into-empty-deque",
		func(t *testing.T) {
			tests := []int{-5, -1, 0, 1, 7}
			q := NewDeque[float64](10)

			for _, idx := range tests {
				_, err := q.TryPeek(idx)
				if err == nil {
					t.Errorf("Expected an error when peeking at %d", idx)
				}
			}
		},
	)
	t.Run(
		"peek-into-semi-filled-deque",
		func(t *testing.T) {
			tests := []struct {
				idx   int
				isErr bool
			}{
				{-5, true},
				{-1, true},
				{0, false},
				{1, false},
				{3, false},
				{7, true},
			}
			q := NewDeque[float64](10)
			q.AppendLeft(0.0)
			q.Append(0.1)
			q.Append(0.2)
			q.Append(0.3)

			for _, testCase := range tests {
				_, err := q.TryPeek(testCase.idx)
				if testCase.isErr != (err != nil) {
					t.Errorf("Expected err: %v, but err is %v", testCase.isErr, err)
				}
			}
		},
	)
}

func TestCount(t *testing.T) {
	const SIZE int = 10

	t.Run(
		"different-elements",
		func(t *testing.T) {
			q := NewDeque[int](SIZE)
			for i := range SIZE {
				q.Append(i)
				if q.Count(i) != 1 {
					t.Errorf("Expected Count() to return 1, got %d", q.Count(i))
				}
			}
		},
	)
	t.Run(
		"same-element",
		func(t *testing.T) {
			q := NewDeque[int](SIZE)
			for i := 1; i <= SIZE; i++ {
				q.Append(0)
				if q.Count(0) != i {
					t.Errorf("Expected Count() to return %d, got %d", i, q.Count(0))
				}
			}
		},
	)
}

func TestIsEmpty(t *testing.T) {
	q := NewUnlimitedDeque[string]()
	if !q.IsEmpty() {
		t.Errorf("Expected IsEmpty() to return true")
	}
	items := [3]string{"foo", "egg", "bar"}
	for _, v := range items {
		q.Append(v)
		if q.IsEmpty() {
			t.Errorf("Expected IsEmpty() to return false")
		}
	}
}

func TestClear(t *testing.T) {
	q := NewUnlimitedDeque[float32]()
	q.Append(1.0)
	q.AppendLeft(42.0)
	q.Append(3.14)
	q.AppendLeft(7.4)
	q.Clear()
	if !q.IsEmpty() {
		t.Errorf("Expected dequeue to become empty after Clear()")
	}
}

func TestCopy(t *testing.T) {
	const CAPACITY int = 5
	q := NewDeque[int](CAPACITY)
	for i := range CAPACITY {
		q.AppendLeft(i)
	}
	nq := q.Copy()
	if q.GetCapacity() != nq.GetCapacity() {
		t.Errorf("Capacity mismatch")
	}
	if q.GetLen() != nq.GetLen() {
		t.Errorf("Length mismatch")
	}
	for i := range CAPACITY - 1 {
		v1 := q.Peek(i)
		v2 := nq.Peek(i)
		if v1 != v2 {
			t.Errorf("Elements at the same index mismatch")
		}
	}
}

func TestRotate(t *testing.T) {
	const CAPACITY int = 3
	q := NewDeque[int](CAPACITY)
	for i := range CAPACITY {
		q.Append(i)
	}
	q.Rotate(1)
	if v := q.Peek(0); v != CAPACITY-1 {
		t.Errorf("Expected first element to be %d, got %d", CAPACITY-1, v)
	}
	q.Rotate(1)
	if v := q.Peek(0); v != CAPACITY-2 {
		t.Errorf("Expected first element to be %d, got %d", CAPACITY-2, v)
	}
	q.Rotate(-1)
	if v := q.Peek(0); v != CAPACITY-1 {
		t.Errorf("Expected first element to be %d, got %d", CAPACITY-1, v)
	}
}
