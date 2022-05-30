package memcached

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	mocknet "github.com/IlyaFloppy/grpcstore/sdk/memcached/mocks"
)

//go:generate mockgen --build_flags=--mod=mod -destination=mocks/conn.go net Conn
func TestConnHappyPaths(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctrl := gomock.NewController(t)
	nc := mocknet.NewMockConn(ctrl)

	key := "key"
	val := []byte("12345")
	c := NewConn(nc)

	nc.EXPECT().Write([]byte("set key 0 0 5\r\n12345\r\n")).Return(22, nil).Times(1)
	nc.EXPECT().Read(gomock.Any()).DoAndReturn(func(dst []byte) (n int, err error) {
		r := "STORED\r\n"
		copy(dst, []byte(r))
		return len(r), nil
	}).Times(1)
	err := c.Set(key, val)
	require.NoError(t, err)

	nc.EXPECT().Write([]byte("get key\r\n")).Return(9, nil).Times(1)
	nc.EXPECT().Read(gomock.Any()).DoAndReturn(func(dst []byte) (n int, err error) {
		r := "VALUE key 0 5\r\n12345\r\nEND\r\n"
		copy(dst, []byte(r))
		return len(r), nil
	})
	v, err := c.Get("key")
	require.NoError(t, err)
	require.Equal(t, val, v)

	nc.EXPECT().Write([]byte("delete key\r\n")).Return(12, nil).Times(1)
	nc.EXPECT().Read(gomock.Any()).DoAndReturn(func(dst []byte) (n int, err error) {
		r := "DELETED\r\n"
		copy(dst, []byte(r))
		return len(r), nil
	})
	err = c.Delete("key")
	require.NoError(t, err)

	nc.EXPECT().Write([]byte("delete key\r\n")).Return(12, nil).Times(1)
	nc.EXPECT().Read(gomock.Any()).DoAndReturn(func(dst []byte) (n int, err error) {
		r := "NOT_FOUND\r\n"
		copy(dst, []byte(r))
		return len(r), nil
	})
	err = c.Delete("key")
	require.ErrorIs(t, err, ErrNotFound)

	nc.EXPECT().Close().Return(nil)
	err = c.Close()
	require.NoError(t, err)
}

func TestValueHeaderRE(t *testing.T) {
	require.True(t, valueHeaderRE.MatchString("VALUE key 0 0\r\n"))
	require.True(t, valueHeaderRE.MatchString("VALUE key 123 123\r\n"))
	require.True(t, valueHeaderRE.MatchString("VALUE key 324 2343 3423\r\n"))
	require.True(t, valueHeaderRE.MatchString("VALUE key 0 0 0\r\n"))

	require.False(t, valueHeaderRE.MatchString("VALUE key 0"))
	require.False(t, valueHeaderRE.MatchString("VALUE key 0 0"))
	require.False(t, valueHeaderRE.MatchString("VALUE key 0 0\n"))
	require.False(t, valueHeaderRE.MatchString("VALUE key 0 0 0"))
	require.False(t, valueHeaderRE.MatchString("VALUE key 0 0 0\n"))
	require.False(t, valueHeaderRE.MatchString("sdfsdf"))
}
