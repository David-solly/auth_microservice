package service

import (
	"context"

	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	generate grpctransport.Handler
	verify   grpctransport.Handler
	affect   grpctransport.Handler
	renew    grpctransport.Handler
}

// implement Server Interface in generate.pb.go
func (s *grpcServer) Generate(ctx context.Context, r *pb.TokenRequest) (*pb.TokenResponse, error) {
	_, resp, err := s.generate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TokenResponse), nil
}

func (s *grpcServer) VerifyToken(ctx context.Context, r *pb.TokenVerifyRequest) (*pb.TokenVerifyResponse, error) {

	_, resp, err := s.verify.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.TokenVerifyResponse), nil
}

func (s *grpcServer) AffectToken(ctx context.Context, r *pb.TokenAffectRequest) (*pb.TokenAffectResponse, error) {

	_, resp, err := s.affect.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.TokenAffectResponse), nil
}

func (s *grpcServer) RenewTokens(ctx context.Context, r *pb.TokenRenewRequest) (*pb.TokenResponse, error) {
	_, resp, err := s.renew.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TokenResponse), nil
}

// create new grpc server
func NewGRPCServer(ctx context.Context, endpoint TokenServiceEndpoints) pb.TokenServiceServer {
	return &grpcServer{
		generate: grpctransport.NewServer(
			endpoint.GenerateEndpoint,
			DecodeGRPCTokenRequest,
			EncodeGRPCTokenResponse,
		),
		verify: grpctransport.NewServer(
			endpoint.VerifyEndpoint,
			DecodeGRPCTokenVerifyRequest,
			EncodeGRPCTokenVerifyResponse,
		),
		affect: grpctransport.NewServer(
			endpoint.AffectEndpoint,
			DecodeGRPCTokenAffectRequest,
			EncodeGRPCTokenAffectResponse,
		),
		renew: grpctransport.NewServer(
			endpoint.RenewEndpoint,
			DecodeGRPCTokenRenewRequest,
			EncodeGRPCTokenResponse,
		),
	}
}
