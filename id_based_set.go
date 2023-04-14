package gogpt

type idBasedItem interface {
	getId() string
}

type idBasedSet[T idBasedItem] struct {
	Content []T
}

// newIdBasedSet creates a new idBasedSet instance with a given capacity
func newIdBasedSet[T idBasedItem](capacity int) *idBasedSet[T] {
	return &idBasedSet[T]{
		Content: make([]T, 0, capacity),
	}
}

// add adds the given element to the current idBasedSet instance
func (s *idBasedSet[T]) add(itemToAdd T) bool {
	if s.contains(itemToAdd) {
		return false
	}
	s.Content = append(s.Content, itemToAdd)
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
	for _, item := range s.Content {
		if item.getId() == itemToVerify.getId() {
			return true
		}
	}
	return false
}

// find finds the element which has the given id in the current idBasedSet
func (s *idBasedSet[T]) find(id string) *T {
	for _, item := range s.Content {
		if item.getId() == id {
			return &item
		}
	}
	return nil
}

// size returns the length of the idBasedSet instance
func (s *idBasedSet[T]) size() int {
	return len(s.Content)
}
