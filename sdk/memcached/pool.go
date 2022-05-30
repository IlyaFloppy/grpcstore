package memcached

import (
	"net"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Pool struct {
	semaphore chan *Conn
	size      int
}

func NewPoolWithAddress(addr string, size int) (pool *Pool, err error) {
	semaphore := make(chan *Conn, size)

	defer func() {
		if err != nil {
			var eg errgroup.Group
			eg.Go(func() error {
				return err
			})
			for i := 0; i < len(semaphore); i++ {
				c := <-semaphore
				eg.Go(func() error {
					return c.Close()
				})
			}
			close(semaphore)

			pool = nil
			err = eg.Wait()
		}
	}()

	for i := 0; i < size; i++ {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to dial")
		}

		semaphore <- NewConn(c)
	}

	return &Pool{
		semaphore: semaphore,
		size:      size,
	}, nil
}

func (p *Pool) Close() error {
	var eg errgroup.Group
	for i := 0; i < p.size; i++ { // close will wait for all other operations to finish to close all underlying connections.
		c := <-p.semaphore

		eg.Go(func() error {
			return c.Close()
		})
	}
	close(p.semaphore)

	return eg.Wait()
}

func (p *Pool) Set(key string, value []byte) error {
	c := <-p.semaphore
	defer func() { p.semaphore <- c }()

	return c.Set(key, value)
}

func (p *Pool) Get(key string) ([]byte, error) {
	c := <-p.semaphore
	defer func() { p.semaphore <- c }()

	return c.Get(key)
}

func (p *Pool) Delete(key string) error {
	c := <-p.semaphore
	defer func() { p.semaphore <- c }()

	return c.Delete(key)
}
