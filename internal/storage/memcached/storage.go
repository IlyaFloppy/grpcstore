package memcached

import (
	"context"

	"github.com/pkg/errors"

	"github.com/IlyaFloppy/grpcstore/internal/config"
	"github.com/IlyaFloppy/grpcstore/sdk/memcached"
)

type Storage struct {
	cfg     config.MemcachedStorageConfig
	client  IMemcachedClient
	readyCh chan struct{}
}

func New(cfg config.MemcachedStorageConfig) *Storage {
	return &Storage{
		cfg:     cfg,
		client:  nil, // will be initialized in Run.
		readyCh: make(chan struct{}),
	}
}

func (s *Storage) Name() string {
	return "memcached-storage"
}

func (s *Storage) Run(ctx context.Context) error {
	var err error
	if s.cfg.UsePool {
		s.client, err = memcached.NewPoolWithAddress(s.cfg.Address, s.cfg.PoolSize)
	} else {
		s.client, err = memcached.NewConnWithAddress(s.cfg.Address)
	}
	if err != nil {
		return errors.Wrap(err, "failed to create memcached client")
	}

	close(s.readyCh)

	<-ctx.Done()
	return nil
}

func (s *Storage) ReadyCh() <-chan struct{} {
	return s.readyCh
}
