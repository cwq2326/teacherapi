package set

// Set is a custom set data structure.
type Set[T comparable] map[T]bool

// Returns a new set.
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// Adds an item to the set.
func (s Set[T]) Add(i T) {
	s[i] = true
}

// Removes an item from the set.
func (s Set[T]) Remove(i T) {
	delete(s, i)
}

// Returns true if the set contains the given item.
func (s Set[T]) Contains(i T) bool {
	_, ok := s[i]
	return ok
}

// Returns the number of items in the set.
func (s Set[T]) Length() int {
	return len(s)
}

// Converts the set to an array and return it.
func (s Set[T]) ToArray() []T {
	result := make([]T, 0, s.Length()) 
	for k := range s {
		result = append(result, k)
	}
	return result
}
