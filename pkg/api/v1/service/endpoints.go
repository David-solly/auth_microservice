package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	hc "github.com/David-solly/auth_microservice/pkg/api/v1/hc"
	"github.com/David-solly/auth_microservice/pkg/api/v1/models"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
)

type TokenServiceEndpoints struct {
	GenerateEndpoint endpoint.Endpoint
	VerifyEndpoint   endpoint.Endpoint
}

type TokenRequest struct {
	Claims map[string]string `json:"claims,omitempty"`
}

type TokenRequest2 struct {
	Claims []string `json:"claims,omitempty"`
}

type TokenResponse struct {
	Response models.AccessTokens
	Error    models.ServiceError
}

type TokenVerifyRequest struct {
	Token   string `json:"token"`
	Service string `json:"service,omitempty"`
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

func MakeTokenServiceGenerateEndpoint(svc TokenServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TokenRequest)

		tkns, err := svc.Generate(ctx, req.Claims)

		if err != nil {
			return nil, err
		}
		return TokenResponse{Response: *tkns}, nil
	}
}

func MakeTokenServiceVerifyEndpoint(svc TokenServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TokenVerifyRequest)

		tkns, err := svc.VerifyToken(ctx, req)

		if err != nil {
			errorMessage := err.(*models.ResponseObject)

			return errorMessage, nil
		}
		return tkns, nil
	}
}

func (te TokenServiceEndpoints) Generate(ctx context.Context, claims map[string]string) (*models.AccessTokens, error) {
	req := TokenRequest{Claims: claims}

	resp, err := te.GenerateEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	tokenRespone := resp.(TokenResponse)

	if tokenRespone.Response.AccessToken == "" {
		return nil, errors.New(fmt.Sprintf("Response was errors [%v]", tokenRespone.Error))
	}

	return &tokenRespone.Response, nil

}

func (te TokenServiceEndpoints) VerifyToken(ctx context.Context, tokenToverify TokenVerifyRequest) (*models.TokenVerifyResponse, interface{}) {

	//fmt.Printf("token to verify %v\n", tokenToverify)
	resp, err := te.VerifyEndpoint(ctx, tokenToverify)
	if err != nil {
		//fmt.Printf("Error in tendpoint %v", err)
		return nil, err
	}
	//fmt.Printf("after in tendpoint %v\n%v", resp, err)

	if tokenRespone, k := resp.(*models.TokenVerifyResponse); k {
		if tokenRespone.Status == models.TokenStatus_INVALID {
			return nil, models.ResponseObject{Error: fmt.Sprintf("Token is invalid for id [%d]", tokenRespone.UserID), Code: http.StatusBadRequest}
			// return nil, errors.New(fmt.Sprintf("Response was errors [%v]", tokenRespone.Error))
		}
		//fmt.Printf("about to relay back @@@---\n%v\ntype:%v", tokenRespone, tokenRespone)

		return tokenRespone, nil
	}

	//fmt.Printf("\nJumped after in tendpoint %v\n%v\n", resp, err)

	if serviceRspone, k := resp.(*models.ResponseObject); k {
		return &models.TokenVerifyResponse{Error: models.ServiceError{Error: serviceRspone.Error, Code: serviceRspone.Code}}, nil
	}

	//fmt.Printf("\nJumped Again after in tendpoint %v\n%v\n", resp, err)
	return nil, models.ResponseObject{Error: fmt.Sprintf("Internal Error processsing verification"), Code: http.StatusBadRequest}

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
