package client

import (
	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// New : gRPc callable Endpoints and data handlers
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

	var tokenAffectEndpoint = grpctransport.NewClient(
		conn, "v1.TokenService", "AffectToken",
		token_grpc.EncodeGRPCTokenAffectRequest,
		token_grpc.DecodeGRPCTokenAffectResponse,
		pb.TokenAffectResponse{},
	).Endpoint()

	var tokenRenewEndpoint = grpctransport.NewClient(
		conn, "v1.TokenService", "RenewTokens",
		token_grpc.EncodeGRPCTokenRenewRequest,
		token_grpc.DecodeGRPCTokenResponse,
		pb.TokenResponse{},
	).Endpoint()

	return token_grpc.TokenServiceEndpoints{
		GenerateEndpoint: tokenEndpoint,
		VerifyEndpoint:   tokenVerifyEndpoint,
		AffectEndpoint:   tokenAffectEndpoint,
		RenewEndpoint:    tokenRenewEndpoint,
	}
}
