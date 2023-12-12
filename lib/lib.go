package lib

func Must[T any](val T, err any) T {
	if err != nil {
		panic(err)
	}

	return val
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
