package service

import (
	"context"
	"errors"
	"fmt"

	hc "github.com/David-solly/auth_microservice/pkg/api/v1/hc"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
)

type TokenServiceEndpoints struct {
	Endpoint endpoint.Endpoint
}

type TokenRequest struct {
	Claims map[string]string `json:"claims,omitempty"`
}

type TokenRequest2 struct {
	Claims []string `json:"claims,omitempty"`
}

type TokenResponse struct {
	Response AccessTokens
	Error    ServiceError
}

//wrapper for the endpoints
type EndpointsConsul struct {
	ConsulHealthCheckEndpoint endpoint.Endpoint
	ConsulHealthWatchEndpoint endpoint.Endpoint
}

type HealthServiceRequest struct {
	Service string `json:"service,omitempty"`
}

type HealthServiceResponse struct {
	Status int `json:"status,omitempty"`
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

func MakeHealthServiceCheckEndpoint(svc Health) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(HealthServiceRequest)
		healthStatus, err := svc.Check(ctx, req.Service)
		if err != nil {
			return HealthServiceResponse{Status: HealthCheckResponse_UNKNOWN}, err
		}

		return HealthServiceResponse{Status: healthStatus}, nil
	}
}

func MakeHealthServiceWatchEndpoint(svc Health) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// req := request.(HealthServiceRequest)

		// sv := HealthWatchServer{}
		// err := svc.Watch(nil, sv)
		// if err != nil {
		// 	return HealthServiceResponse{Status: HealthCheckResponse_UNKNOWN}, err
		// }

		return HealthServiceResponse{Status: HealthCheckResponse_UNKNOWN}, nil
	}
}

type HealthWatchServer struct {
	grpc.ServerStream
}

func (x *HealthWatchServer) Send(m *hc.HealthCheckResponse) error {
	return x.ServerStream.SendMsg(m)
}
