package service

import (
	"context"

	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	generate grpctransport.Handler
}

// implement LoremServer Interface in generate.pb.go
func (s *grpcServer) Generate(ctx context.Context, r *pb.TokenRequest) (*pb.TokenResponse, error) {
	_, resp, err := s.generate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TokenResponse), nil
}

// create new grpc server
func NewGRPCServer(ctx context.Context, endpoint TokenServiceEndpoints) pb.TokenServiceServer {
	return &grpcServer{
		generate: grpctransport.NewServer(

			endpoint.Endpoint,
			DecodeGRPCTokenRequest,
			EncodeGRPCTokenResponse,
		),
	}
}
