package gogpt

type idBasedSetItem interface {
	getId() string
}

type idBasedSet[T idBasedSetItem] struct {
	content []T
}

// newIdBasedSet creates a new idBasedSet instance with a given capacity
func newIdBasedSet[T idBasedSetItem](capacity int) *idBasedSet[T] {
	return &idBasedSet[T]{
		content: make([]T, 0, capacity),
	}
}

// add adds the given element to the current idBasedSet instance
func (s *idBasedSet[T]) add(itemToAdd T) bool {
	if s.contains(itemToAdd) {
		return false
	}
	s.content = append(s.content, itemToAdd)
	return true
}

// addAll adds the given list of elements in tho the current idBasedSet instance
func (s *idBasedSet[T]) addAll(itemsToAdd []T) {
	for _, itemToAdd := range itemsToAdd {
		s.add(itemToAdd)
	}
}

// contains check if the given item is in the current idBasedSet instance or not
func (s *idBasedSet[T]) contains(itemToVerify T) bool {
	for _, item := range s.content {
		if item.getId() == itemToVerify.getId() {
			return true
		}
	}
	return false
}

// size returns the length of the idBasedSet instance
func (s *idBasedSet[T]) size() int {
	return len(s.content)
}
