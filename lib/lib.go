package lib

func Must[T any](val T, err any) T {
	if err != nil {
		panic(err)
	}

	return val
}

func Clamp(n, min, max int) int {
	if n > max {
		return max
	}

	if n < min {
		return min
	}

	return n
}

type Queue[T any] []T

func NewQueue[T any]() Queue[T] {
	return make([]T, 0)
}

func (self *Queue[T]) Append(val T) {
	*self = append(*self, val)
}

func (self *Queue[T]) Pop() T {
	q := *self
	l := len(q)
	top := q[0]
	*self = q[1:l]
	return top
}

type Heap[T any] struct {
	elements  []T
	extractor func(T) int
}

func NewHeap[T any](extractor func(T) int) Heap[T] {
	return Heap[T]{
		elements:  make([]T, 0),
		extractor: extractor,
	}
}

func (self Heap[T]) Len() int {
	return len(self.elements)
}

func (self *Heap[T]) Push(val T) {
	h := *&self.elements
	h = append(h, val)
	newValIndex := len(h) - 1
	parentIndex := (newValIndex - 1) / 2
	for parentIndex >= 0 &&
		self.extractor(h[newValIndex]) < self.extractor(h[parentIndex]) {
		h[newValIndex] = h[parentIndex]
		h[parentIndex] = val
		newValIndex = parentIndex
		parentIndex = (parentIndex - 1) / 2
	}
	*&self.elements = h
}

func (self *Heap[T]) Pop() T {
	h := *&self.elements
	l := len(h)
	top := h[0]
	h[0] = h[l-1]
	sinkIndex := 0
	for sinkIndex < l-1 {
		left := sinkIndex*2 + 1
		right := sinkIndex*2 + 2

		if right < l-1 &&
			self.extractor(h[right]) < self.extractor(h[sinkIndex]) &&
			self.extractor(h[right]) < self.extractor(h[left]) {
			h[sinkIndex] = h[right]
			h[right] = h[l-1]
			sinkIndex = right
			continue
		}

		if left < l-1 &&
			self.extractor(h[left]) < self.extractor(h[sinkIndex]) {
			h[sinkIndex] = h[left]
			h[left] = h[l-1]
			sinkIndex = left
			continue
		}

		break
	}

	*&self.elements = h[0 : l-1]
	return top
}

type Stack[T any] []T

func NewStack[T any]() Stack[T] {
	return make([]T, 0)
}

func (self *Stack[T]) Append(val T) {
	*self = append(*self, val)
}

func (self *Stack[T]) Pop() T {
	s := *self
	l := len(s)
	last := s[l-1]
	*self = s[0 : l-1]
	return last
}

type Set[T comparable] map[T]bool

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (self *Set[T]) Add(val T) bool {
	s := *self
	exists := s[val]
	s[val] = true
	*self = s
	return !exists
}

func (self *Set[T]) Contains(val T) bool {
	s := *self
	return s[val]
}

func (self *Set[T]) Remove(val T) bool {
	s := *self
	exists := s[val]
	s[val] = false
	*self = s
	return exists
}

func PopSlice[T any](arr *[]T) T {
	a := *arr
	l := len(a)
	top := a[l-1]
	*arr = a[:l-1]
	return top
}

func Fill[T any](arr *[]T, val T) {
	a := *arr
	for i := 0; i < len(a); i++ {
		a[i] = val
	}
}
