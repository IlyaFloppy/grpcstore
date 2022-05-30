package interceptors

import (
	"context"
	"time"

	"github.com/google/uuid"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func WithLoggingUnaryInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = logger.WithContext(ctx)
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.
				Str("request_id", uuid.New().String()).
				Str("ip", getIPAddr(ctx)).
				Str("method", info.FullMethod)
		})

		if logger.GetLevel() <= zerolog.DebugLevel {
			logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Interface("request", req)
			})
		}

		start := time.Now()
		res, err := handler(ctx, req)

		s, _ := status.FromError(err)
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.
				Dur("duration", time.Since(start)).
				Uint32("code", uint32(s.Code()))
		})

		details := s.Details()
		if len(details) > 0 {
			logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Interface("details", details)
			})
		}

		//nolint:exhaustive
		switch s.Code() {
		case codes.OK:
			logger.Info().Msg("unary handler finished successfully")
		case codes.Canceled:
			logger.Info().Msg("unary handler was canceled")
		case codes.DeadlineExceeded:
			logger.Info().Msg("unary handler was canceled due to deadline exceeding")
		default:
			logger.Err(err).Msg("unary handler error")
		}

		return res, err
	}
}

func WithLoggingStreamInterceptor(logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := logger.WithContext(ss.Context())
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.
				Str("request_id", uuid.New().String()).
				Str("ip", getIPAddr(ctx)).
				Str("method", info.FullMethod)
		})

		newStream := grpcmiddleware.WrapServerStream(ss)
		newStream.WrappedContext = ctx

		logger.Debug().Msg("handling stream request")

		start := time.Now()
		err := handler(srv, newStream)

		s, _ := status.FromError(err)
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.
				Dur("duration", time.Since(start)).
				Uint32("code", uint32(s.Code()))
		})

		details := s.Details()
		if len(details) > 0 {
			logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Interface("details", details)
			})
		}

		//nolint:exhaustive
		switch s.Code() {
		case codes.OK:
			logger.Info().Msg("stream handler finished successfully")
		case codes.Canceled:
			logger.Info().Msg("stream handler was canceled")
		case codes.DeadlineExceeded:
			logger.Info().Msg("stream handler was canceled due to deadline exceeding")
		default:
			logger.Err(err).Msg("stream handler error")
		}

		return err
	}
}

func getIPAddr(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if ok {
		return p.Addr.String()
	}

	return "<unknown ip>"
}
