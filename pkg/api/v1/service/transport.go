package service

import (
	"context"

	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
	hc "github.com/David-solly/auth_microservice/pkg/api/v1/hc"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	generate grpctransport.Handler
}

type grpcServerHealth struct {
	health grpctransport.Handler
	watch  grpctransport.Handler
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

// create new grpc server
func NewGRPCServerHealth(ctx context.Context, endpoint EndpointsConsul) hc.HealthServer {
	return &grpcServerHealth{
		health: grpctransport.NewServer(
			endpoint.ConsulHealthCheckEndpoint,
			DecodeGRPCHealthServiceRequest,
			EncodeGRPCHealthServiceResponse,
		),
	}
}

func (s *grpcServerHealth) Check(ctx context.Context, r *hc.HealthCheckRequest) (*hc.HealthCheckResponse, error) {
	_, resp, err := s.health.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*hc.HealthCheckResponse), nil
}

func (s *grpcServerHealth) Watch(req *hc.HealthCheckRequest, srv hc.Health_WatchServer) error {
	_, _, err := s.watch.ServeGRPC(srv.Context(), req)

	return err
}
