package store

type Store struct {
	store map[string]string	
}

func New() Store {
	s := Store{ store: make(map[string]string) }
	return s
}

func (s Store) Set(key string, value string) {
	s.store[key] = value
}

func (s Store) Get(key string) string {
	v := s.store[key]
	return v
}
