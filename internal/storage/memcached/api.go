package memcached

import "github.com/pkg/errors"

func (s *Storage) Get(key string) ([]byte, error) {
	res, err := s.client.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get key")
	}

	return res, nil
}

func (s *Storage) Set(key string, value []byte) error {
	err := s.client.Set(key, value)
	if err != nil {
		return errors.Wrap(err, "failed to set key")
	}

	return nil
}

func (s *Storage) Delete(key string) error {
	err := s.client.Delete(key)
	if err != nil {
		return errors.Wrap(err, "failed to delete key")
	}

	return nil
}
