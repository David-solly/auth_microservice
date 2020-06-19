package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
	"github.com/David-solly/auth_microservice/pkg/api/v1/models"
)

// Encode and Decode Renew request
func DecodeGRPCTokenRenewRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.TokenRenewRequest)
	return TokenRenewRequest{RefreshToken: req.RefreshToken}, nil
}

//Encode and Decode Token Request
func EncodeGRPCTokenRenewRequest(_ context.Context, r interface{}) (interface{}, error) {
	n := r.(TokenRenewRequest)

	return &pb.TokenRenewRequest{RefreshToken: n.RefreshToken}, nil
}

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
	resp := r.(*TokenResponse)

	return &pb.TokenResponse{
		Error: &pb.ServiceError{
			Error: resp.Error.Error, Code: int32(resp.Error.Code)},
		Tokens: &pb.TokenPair{
			AuthToken:    resp.Response.AccessToken,
			RefreshToken: resp.Response.RefreshToken}}, nil
}

func DecodeGRPCTokenResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.TokenResponse)
	rs := TokenResponse{
		Error: models.ServiceError{Error: resp.Error.Error, Code: int(resp.Error.Code)},
		Response: models.AccessTokens{
			AccessToken:  resp.Tokens.AuthToken,
			RefreshToken: resp.Tokens.RefreshToken,
		},
	}
	return rs, nil

}

//Encode and Decode Token Affect Request
func EncodeGRPCTokenAffectRequest(_ context.Context, r interface{}) (interface{}, error) {
	n := r.(models.TokenAffectRequest)

	return &pb.TokenAffectRequest{Token: n.Token, DesiredState: pb.TokenState(n.DesiredState)}, nil
}

func DecodeGRPCTokenAffectRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.TokenAffectRequest)
	return models.TokenAffectRequest{Token: req.Token, DesiredState: models.TokenState(req.DesiredState)}, nil
}

// Encode and Decode Token Response
func EncodeGRPCTokenAffectResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*models.TokenAffectResponse)
	if resp.Error != nil {
		return &pb.TokenAffectResponse{
			Error: &pb.ServiceError{
				Error: resp.Error.Error, Code: int32(resp.Error.Code)},
			EffectApplied: resp.EffectApplied}, nil
	}
	return &pb.TokenAffectResponse{
		EffectApplied: resp.EffectApplied}, nil
}

func DecodeGRPCTokenAffectResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.TokenAffectResponse)
	if resp.Error == nil {
		return models.TokenAffectResponse{
			EffectApplied: resp.EffectApplied,
		}, nil
	}
	return models.TokenAffectResponse{
		EffectApplied: resp.EffectApplied,
		Error: &models.ServiceError{
			Error: resp.Error.Error, Code: int(resp.Error.Code)},
	}, nil
}

/// Verify Request response - encode/decode
func EncodeGRPCTokenVerifyRequest(_ context.Context, r interface{}) (interface{}, error) {
	n := r.(TokenVerifyRequest)

	return &pb.TokenVerifyRequest{Token: n.Token, Service: n.Service}, nil
}

func DecodeGRPCTokenVerifyRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.TokenVerifyRequest)
	return TokenVerifyRequest{Token: req.Token, Service: req.Service}, nil
}

// Encode and Decode Token Response
func EncodeGRPCTokenVerifyResponse(_ context.Context, r interface{}) (interface{}, error) {

	if resp, k := r.(*models.TokenVerifyResponse); k {
		return &pb.TokenVerifyResponse{Access: &pb.ServiceAccess{
			Status: pb.TokenStatus(resp.Access.Status), UserId: fmt.Sprintf("%d", resp.Access.UserID), Claims: resp.MapClaims()}}, nil
	}

	if resp, k := r.(*models.ResponseObject); k {

		return &pb.TokenVerifyResponse{
			Error: &pb.ServiceError{
				Error: resp.Error, Code: int32(resp.Code)},
		}, nil
	}

	return &pb.TokenVerifyResponse{
		Error: &pb.ServiceError{
			Error: "Unknown error", Code: int32(http.StatusInternalServerError)},
	}, nil
}

func DecodeGRPCTokenVerifyResponse(_ context.Context, r interface{}) (interface{}, error) {

	if resp, k := r.(*pb.TokenVerifyResponse); k {
		id, _ := strconv.ParseUint(resp.Access.UserId, 10, 64)

		if resp.Error != nil {
			return &models.ResponseObject{Error: resp.Error.Error, Code: int(resp.Error.Code)}, nil
		}

		c := MergeClaims(resp.Access.Claims)

		return &models.TokenVerifyResponse{Access: models.ServiceAccess{
			UserID: id,
			Status: models.TokenStatus(resp.Access.Status),
			Claims: &c,
		}}, nil

	}

	if resp, k := r.(*models.ResponseObject); k {
		return resp, nil
	}

	return &models.TokenVerifyResponse{
		Error: models.ServiceError{
			Error: "Unknown error 123res", Code: http.StatusInternalServerError},
	}, nil
}
