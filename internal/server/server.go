package server

import (
	"context"
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/IlyaFloppy/grpcstore/internal/config"
	"github.com/IlyaFloppy/grpcstore/internal/server/interceptors"
	"github.com/IlyaFloppy/grpcstore/public-api/pb"
)

type Server struct {
	pb.UnimplementedGRPCStoreServiceServer

	logger     zerolog.Logger
	cfg        config.ServerConfig
	grpcServer *grpc.Server
	readyCh    chan struct{}
	storage    IStorage
}

func New(logger zerolog.Logger, cfg config.ServerConfig, storage IStorage) *Server {
	logger = logger.With().Str("component", (*Server)(nil).Name()).Logger()

	return &Server{
		logger: logger,
		cfg:    cfg,
		grpcServer: grpc.NewServer(
			grpc.KeepaliveParams(keepalive.ServerParameters{
				Time:    cfg.KeepAliveTime,
				Timeout: cfg.KeepAliveTimeout,
			}),
			grpc.WriteBufferSize(cfg.WriteBufferSize),
			grpc.ReadBufferSize(cfg.ReadBufferSize),
			grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
				interceptors.WithLoggingUnaryInterceptor(logger),
				interceptors.WithRecoveryUnaryInterceptor(logger),
			)),
			grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
				interceptors.WithLoggingStreamInterceptor(logger),
				interceptors.WithRecoveryStreamInterceptor(logger),
			)),
		),
		readyCh: make(chan struct{}),
		storage: storage,
	}
}

func (s *Server) Name() string {
	return "api-server"
}

func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		return err
	}

	pb.RegisterGRPCStoreServiceServer(s.grpcServer, s)

	go func() {
		<-ctx.Done()
		s.grpcServer.GracefulStop()
	}()

	s.logger.Info().Str("address", s.cfg.Address).Msg("server started listening")
	close(s.readyCh)

	return s.grpcServer.Serve(lis)
}

func (s *Server) ReadyCh() <-chan struct{} {
	return s.readyCh
}
