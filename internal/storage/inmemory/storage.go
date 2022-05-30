package inmemory

import (
	"context"
)

type Storage struct {
	readyCh chan struct{}
	hm      map[string][]byte
}

func New() *Storage {
	return &Storage{
		readyCh: make(chan struct{}),
		hm:      make(map[string][]byte),
	}
}

func (s *Storage) Name() string {
	return "inmemory-storage"
}

func (s *Storage) Run(ctx context.Context) error {
	close(s.readyCh)

	<-ctx.Done()
	return nil
}

func (s *Storage) ReadyCh() <-chan struct{} {
	return s.readyCh
}
