package server

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/IlyaFloppy/grpcstore/internal/config"
	servermocks "github.com/IlyaFloppy/grpcstore/internal/server/mocks"
	"github.com/IlyaFloppy/grpcstore/public-api/pb"
)

func TestGet(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctrl := gomock.NewController(t)

	storage := servermocks.NewMockIStorage(ctrl)
	server := New(zerolog.New(os.Stderr), config.ServerConfig{}, storage)

	t.Run("happy case", func(t *testing.T) {
		storage.EXPECT().Get("key").Return([]byte("12345"), nil)
		res, err := server.Get(context.Background(), &pb.GetRequest{
			Key: "key",
		})
		require.NoError(t, err)
		require.Equal(t, &pb.GetResult{
			Value: []byte("12345"),
		}, res)
	})

	t.Run("internal error", func(t *testing.T) {
		storage.EXPECT().Get("key").Return(nil, errors.New("failed on purpose"))
		res, err := server.Get(context.Background(), &pb.GetRequest{
			Key: "key",
		})
		require.Error(t, err)
		require.Nil(t, res)
		require.Equal(t, codes.Internal, status.Code(err))
	})
}

func TestSet(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctrl := gomock.NewController(t)

	storage := servermocks.NewMockIStorage(ctrl)
	server := New(zerolog.New(os.Stderr), config.ServerConfig{}, storage)

	t.Run("happy case", func(t *testing.T) {
		storage.EXPECT().Set("key", []byte("12345")).Return(nil)
		res, err := server.Set(context.Background(), &pb.SetRequest{
			Key:   "key",
			Value: []byte("12345"),
		})
		require.NoError(t, err)
		require.Equal(t, &pb.SetResult{}, res)
	})

	t.Run("internal error", func(t *testing.T) {
		storage.EXPECT().Set("key", []byte("12345")).Return(errors.New("failed on purpose"))
		res, err := server.Set(context.Background(), &pb.SetRequest{
			Key:   "key",
			Value: []byte("12345"),
		})
		require.Error(t, err)
		require.Nil(t, res)
		require.Equal(t, codes.Internal, status.Code(err))
	})
}

func TestDelete(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctrl := gomock.NewController(t)

	storage := servermocks.NewMockIStorage(ctrl)
	server := New(zerolog.New(os.Stderr), config.ServerConfig{}, storage)

	t.Run("happy case", func(t *testing.T) {
		storage.EXPECT().Delete("key").Return(nil)
		res, err := server.Delete(context.Background(), &pb.DeleteRequest{
			Key: "key",
		})
		require.NoError(t, err)
		require.Equal(t, &pb.DeleteResult{}, res)
	})

	t.Run("internal error", func(t *testing.T) {
		storage.EXPECT().Delete("key").Return(errors.New("failed on purpose"))
		res, err := server.Delete(context.Background(), &pb.DeleteRequest{
			Key: "key",
		})
		require.Error(t, err)
		require.Nil(t, res)
		require.Equal(t, codes.Internal, status.Code(err))
	})
}
