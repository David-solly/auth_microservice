package service

import (
	"context"

	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
)

//Encode and Decode Token Request
func EncodeGRPCTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	n := r.(TokenRequest)

	return &pb.TokenRequest{Claims: n.Claims}, nil
}

func DecodeGRPCTokenRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.TokenRequest)
	return TokenRequest{Claims: req.Claims}, nil
}

// Encode and Decode Token Response
func EncodeGRPCTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(TokenResponse)
	return &pb.TokenResponse{
		Error: &pb.ServiceError{
			Error: resp.Error.Error, Code: int32(resp.Error.Code)},
		Tokens: &pb.TokenPair{
			AuthToken:    resp.Response.AccessToken,
			RefreshToken: resp.Response.RefreshToken}}, nil
}

func DecodeGRPCTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.TokenResponse)
	return TokenResponse{
		Response: AccessTokens{
			AccessToken:  resp.Tokens.AuthToken,
			RefreshToken: resp.Tokens.RefreshToken,
		},
	}, nil
}
