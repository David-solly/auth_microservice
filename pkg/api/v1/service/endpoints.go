package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type TokenRequest struct {
	Claims map[string]string
}

type TokenResponse struct {
	Response AccessTokens
	Error    ServiceError
}

type TokenServiceEndpoints struct {
	Endpoint endpoint.Endpoint
}

func MakeTokenServiceEndpoint(svc TokenServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TokenRequest)

		tkns, err := svc.Generate(ctx, req.Claims)

		if err != nil {
			return nil, err
		}
		return TokenResponse{Response: *tkns}, nil
	}
}

func (te TokenServiceEndpoints) Generate(ctx context.Context, claims map[string]string) (*AccessTokens, error) {
	req := TokenRequest{Claims: claims}

	resp, err := te.Endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	tokenRespone := resp.(TokenResponse)

	if tokenRespone.Response.AccessToken == "" {
		return nil, errors.New(fmt.Sprintf("Response was errors [%v]", tokenRespone.Error))
	}

	return &tokenRespone.Response, nil

}
