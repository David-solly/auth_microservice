package client

import (
	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func New(conn *grpc.ClientConn) token_grpc.TokenServiceInterface {
	var tokenEndpoint = grpctransport.NewClient(
		conn, "v1.TokenService", "Generate",
		token_grpc.EncodeGRPCTokenRequest,
		token_grpc.DecodeGRPCTokenResponse,
		pb.TokenResponse{},
	).Endpoint()

	var tokenVerifyEndpoint = grpctransport.NewClient(
		conn, "v1.TokenService", "VerifyToken",
		token_grpc.EncodeGRPCTokenVerifyRequest,
		token_grpc.DecodeGRPCTokenVerifyResponse,
		pb.TokenVerifyResponse{},
	).Endpoint()

	return token_grpc.TokenServiceEndpoints{
		GenerateEndpoint: tokenEndpoint,
		VerifyEndpoint:   tokenVerifyEndpoint,
	}
}
