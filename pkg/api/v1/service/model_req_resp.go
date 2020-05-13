package service

import (
	"context"

	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
)

//Encode and Decode Token Request
func EncodeGRPCTokenRequest(_ context.Context, r interface{}) (interface{}, error) {
	n := r.(TokenRequest)

	return &pb.TokenRequest{}, nil
}

func DecodeGRPCTokenRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.TokenRequest)
	return TokenRequest{
		RequestType: req.RequestType,
		Max:         req.Max,
		Min:         req.Min,
	}, nil
}

// Encode and Decode Token Response
func EncodeGRPCTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(TokenResponse)
	return &pb.TokenResponse{
		Message: resp.Message,
		Err:     resp.Err,
	}, nil
}

func DecodeGRPCTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.TokenResponse)
	return TokenResponse{
		Message: resp.Message,
		Err:     resp.Err,
	}, nil
}
