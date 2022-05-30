package memcached

import (
	"errors"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestPoolWithRealServer(t *testing.T) {
	defer goleak.VerifyNone(t)

	pool, err := NewPoolWithAddress("localhost:11211", 3)
	if errors.Is(err, syscall.ECONNREFUSED) {
		t.Skip() // skip when there is no memcached server.
	}
	require.NoError(t, err)

	key := "key"
	val := []byte("12345")

	err = pool.Set(key, val)
	require.NoError(t, err)

	v, err := pool.Get("key")
	require.NoError(t, err)
	require.Equal(t, val, v)

	err = pool.Delete("key")
	require.NoError(t, err)

	err = pool.Delete("key")
	require.ErrorIs(t, err, ErrNotFound)

	err = pool.Close()
	require.NoError(t, err)
}
