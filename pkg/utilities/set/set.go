package set

// Set is a custom set data structure
type Set[T comparable] map[T]bool

// New creates a new set
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// Add adds an item to the set
func (s Set[T]) Add(i T) {
	s[i] = true
}

// Remove removes an item from the set
func (s Set[T]) Remove(i T) {
	delete(s, i)
}

// Contains returns true if the set contains the given item
func (s Set[T]) Contains(i T) bool {
	_, ok := s[i]
	return ok
}

// Length returns the number of items in the set
func (s Set[T]) Length() int {
	return len(s)
}

// ToArray converts the set to an array
func (s Set[T]) ToArray() []T {
	result := make([]T, 0, s.Length()) 
	for k := range s {
		result = append(result, k)
	}
	return result
}
