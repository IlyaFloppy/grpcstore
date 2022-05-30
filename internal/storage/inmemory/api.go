package inmemory

import "github.com/IlyaFloppy/grpcstore/internal/storage"

func (s *Storage) Get(key string) ([]byte, error) {
	if val, ok := s.hm[key]; ok {
		return val, nil
	}

	return nil, storage.ErrNotFound
}

func (s *Storage) Set(key string, value []byte) error {
	s.hm[key] = value

	return nil
}

func (s *Storage) Delete(key string) error {
	delete(s.hm, key)

	return nil
}
