package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type TokenRequest struct {
	Claims map[string]interface{}
}

type TokenResponse struct {
	Response TokenResponseInterface
}

type TokenResponseInterface interface {
	isResponse()
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
		return TokenResponse{Response: tkns}, nil
	}
}

func (te *TokenServiceEndpoints) Generate(ctx context.Context, claims map[string]interface{}) (*AccessTokens, error) {
	req := TokenRequest{Claims: claims}

	resp, err := te.Endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	tokenRespone := resp.(TokenResponse)
	response := tokenRespone.Response.(AccessTokens)
	if response.AccessToken == "" {
		return nil, errors.New(fmt.Sprintf("Response was errors [%v]", tokenRespone.Response))
	}

	return &response, nil

}
