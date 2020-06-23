package healthservice

import (
	"context"

	hc "github.com/David-solly/consul_hcsd/hc"
)

///Consul Health service checks

func EncodeGRPCHealthServiceRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(HealthServiceRequest)

	return &hc.HealthCheckRequest{Service: req.Service}, nil
}

func DecodeGRPCHealthServiceRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*hc.HealthCheckRequest)
	return HealthServiceRequest{
		Service: req.Service,
	}, nil
}

//response

func EncodeGRPCHealthServiceResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(HealthServiceResponse)
	return &hc.HealthCheckResponse{Status: hc.HealthCheckResponse_ServingStatus(resp.Status)}, nil
}

func DecodeGRPCHealthServiceResponse(ctx context.Context, r interface{}) (interface{}, error) {
	resp := r.(*hc.HealthCheckResponse)
	return HealthServiceResponse{
		Status: int(resp.Status),
	}, nil
}
