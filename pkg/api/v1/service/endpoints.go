package service

import (
	"context"
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
	AffectEndpoint   endpoint.Endpoint
	RenewEndpoint    endpoint.Endpoint
}

type TokenRequest struct {
	Claims map[string]string `json:"claims,omitempty"`
}

type TokenRenewRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
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

func MakeTokenServicAffectEndpoint(svc TokenServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(models.TokenAffectRequest)

		effect, err := svc.AffectToken(ctx, req)

		if err != nil {
			return nil, err
		}
		return effect, nil
	}
}

func MakeTokenServiceGenerateEndpoint(svc TokenServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TokenRequest)

		tkns, err := svc.Generate(ctx, req.Claims)

		if err != nil {
			return nil, err
		}
		return &TokenResponse{Response: *tkns}, nil
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

func MakeTokenServiceRenewEndpoint(svc TokenServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TokenRenewRequest)

		tkns, err := svc.RenewTokens(ctx, req)

		if err != nil {
			return &TokenResponse{Error: models.ServiceError{Error: err.Error(), Code: http.StatusExpectationFailed}}, nil
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
		return nil, fmt.Errorf("Response was errors [%v]", tokenRespone.Error)
	}

	return &tokenRespone.Response, nil

}

func (te TokenServiceEndpoints) AffectToken(ctx context.Context, tokenAffectRequest models.TokenAffectRequest) (*models.TokenAffectResponse, error) {

	resp, err := te.AffectEndpoint(ctx, tokenAffectRequest)
	if err != nil {
		return nil, err
	}

	tokenRespone, k := resp.(models.TokenAffectResponse)

	if !k {
		return nil, fmt.Errorf("Error Transfering request")
	}

	return &tokenRespone, nil

}

func (te TokenServiceEndpoints) RenewTokens(ctx context.Context, refreshToken TokenRenewRequest) (*TokenResponse, error) {

	resp, err := te.RenewEndpoint(ctx, refreshToken)
	if err != nil {
		println(err.Error())
	}

	tokenRespone := resp.(TokenResponse)
	return &tokenRespone, nil

}

func (te TokenServiceEndpoints) VerifyToken(ctx context.Context, tokenToverify TokenVerifyRequest) (*models.TokenVerifyResponse, interface{}) {
	resp, err := te.VerifyEndpoint(ctx, tokenToverify)
	if err != nil {
		return nil, err
	}

	if tokenRespone, k := resp.(*models.TokenVerifyResponse); k {
		if tokenRespone.Access.Status == models.TokenStatus_INVALID {
			return nil, models.ResponseObject{Error: "Token is invalid or expired.", Code: http.StatusBadRequest}
		}

		if tokenRespone.Access.Status == models.TokenStatus_RESTRICTED {
			if svc, ok := tokenRespone.MapClaims()["service"]; ok {
				return nil, models.ResponseObject{Error: fmt.Sprintf("Token is Restricted for service : %s", svc), Code: http.StatusUnauthorized}
			}
			return nil, models.ResponseObject{Error: "Token is restricted for the service.", Code: http.StatusUnauthorized}

		}
		return tokenRespone, nil
	}

	if serviceRspone, k := resp.(*models.ResponseObject); k {
		return &models.TokenVerifyResponse{Error: models.ServiceError{Error: serviceRspone.Error, Code: serviceRspone.Code}}, nil
	}

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
		return HealthServiceResponse{Status: HealthCheckResponse_UNKNOWN}, nil
	}
}

type HealthWatchServer struct {
	grpc.ServerStream
}

func (x *HealthWatchServer) Send(m *hc.HealthCheckResponse) error {
	return x.ServerStream.SendMsg(m)
}
