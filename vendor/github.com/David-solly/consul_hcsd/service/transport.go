package healthservice

import (
	"context"

	hc "github.com/David-solly/consul_hcsd/hc"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServerHealth struct {
	health grpctransport.Handler
	watch  grpctransport.Handler
}

// create new grpc server
func NewGRPCServerHealth(ctx context.Context, endpoint EndpointsConsul) hc.HealthServer {
	return &grpcServerHealth{
		health: grpctransport.NewServer(
			endpoint.ConsulHealthCheckEndpoint,
			DecodeGRPCHealthServiceRequest,
			EncodeGRPCHealthServiceResponse,
		),
	}
}

func (s *grpcServerHealth) Check(ctx context.Context, r *hc.HealthCheckRequest) (*hc.HealthCheckResponse, error) {
	_, resp, err := s.health.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*hc.HealthCheckResponse), nil
}

func (s *grpcServerHealth) Watch(req *hc.HealthCheckRequest, srv hc.Health_WatchServer) error {
	_, _, err := s.watch.ServeGRPC(srv.Context(), req)

	return err
}
