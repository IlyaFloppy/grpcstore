package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/IlyaFloppy/grpcstore/public-api/pb"
)

func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResult, error) {
	data, err := s.storage.Get(req.GetKey())
	if err != nil {
		return nil, status.Errorf(errCode(err), "failed to get key: %s", err.Error())
	}

	return &pb.GetResult{
		Value: data,
	}, nil
}

func (s *Server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResult, error) {
	err := s.storage.Set(req.GetKey(), req.GetValue())
	if err != nil {
		return nil, status.Errorf(errCode(err), "failed to set key: %s", err.Error())
	}

	return &pb.SetResult{}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResult, error) {
	err := s.storage.Delete(req.GetKey())
	if err != nil {
		return nil, status.Errorf(errCode(err), "failed to delete key: %s", err.Error())
	}

	return &pb.DeleteResult{}, nil
}

func errCode(err error) codes.Code {
	switch {
	case implements[interface{ NotFoundErrorMarker() }](err):
		return codes.NotFound
	case implements[interface{ UnknownErrorMarker() }](err):
		return codes.Unknown
	}

	return codes.Internal
}

func implements[T any](err error) bool {
	for err != nil {
		if _, ok := err.(T); ok {
			return true
		}

		err = errors.Unwrap(err)
	}

	return false
}
