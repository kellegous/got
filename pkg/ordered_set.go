package pkg

type OrderedSet[T comparable] struct {
	idx  map[T]int
	vals []T
}

func (s *OrderedSet[T]) Add(t T) int {
	if ix, ok := s.idx[t]; ok {
		return ix
	}

	ix := len(s.vals)
	s.idx[t] = ix
	s.vals = append(s.vals, t)
	return ix
}

func (s *OrderedSet[T]) Values() []T {
	return s.vals
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		idx: map[T]int{},
	}
}
