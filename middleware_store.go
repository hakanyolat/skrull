package skrull

type MiddlewareStore struct {
	data map[string]interface{}
}

// Set ..
func (s *MiddlewareStore) Set(key string, value interface{}) {
	s.data[key] = value
}

// Get ..
func (s *MiddlewareStore) Get(key string) interface{} {
	value, _ := s.data[key]
	return value
}

// Has ..
func (s *MiddlewareStore) Has(key string) bool {
	_, ok := s.data[key]
	return ok
}