package interceptors

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WithRecoveryUnaryInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		defer func() {
			if err2 := recover(); err2 != nil {
				logger.Error().Interface("panic", err).Msg("recovered panic")
				err = status.Errorf(codes.Internal, "panic occurred: %v", err2)
			}
		}()

		res, err = handler(ctx, req)
		return
	}
}

func WithRecoveryStreamInterceptor(logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if err2 := recover(); err2 != nil {
				logger.Error().Interface("panic", err).Msg("recovered panic")
				err = status.Errorf(codes.Internal, "panic occurred: %v", err2)
			}
		}()

		err = handler(srv, ss)
		return
	}
}
