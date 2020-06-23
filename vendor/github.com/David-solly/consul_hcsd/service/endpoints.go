package healthservice

import (
	"context"

	hc "github.com/David-solly/consul_hcsd/hc"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
)

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
