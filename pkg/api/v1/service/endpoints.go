package v1.service

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"errors"
)

type TokenRequest struct{
	Claims map[string]interface{}
}

type TokenResponse struct{
	Response TokenResponseInterface
}

type TokenResponseInterface{
	isResponse()
}

type TokenServiceEndpoints struct{
	Endpoint endpoint.Endpoint
}

func MakeTokenServiceEndpoint(svc TokenServiceInterface)endpoint.Endpoint{
	return func (ctx context.Context, request interface{})(interface{}, error){
		req:=request.(TokenRequest)

		tkns,err:=svc.Generate(ctx, req.Claims)

		if err != nil{
			return nil, err
		}
		return TokenResponse{Response:tkns},nil
	}
}

func (te TokenServiceEndpoints)Generate(ctx context.Context, claims map[string]interface{}) (*AccessTokens, error) {
	req:=TokenRequest{Claims:claims}

	resp, err:=te.Endpoint(ctx,req)
	if err !=nil{
		return nil,err
	}

	tokenRespone,err1:=resp.(TokenResponse)
	if err1 !=nil{
		return nil,err1
	}
	
}