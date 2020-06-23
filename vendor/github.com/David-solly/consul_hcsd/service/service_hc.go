package healthservice

import (
	"context"

	hc "github.com/David-solly/consul_hcsd/hc"
)

const (
	HealthCheckResponse_UNKNOWN     int = 0
	HealthCheckResponse_SERVING     int = 1
	HealthCheckResponse_NOT_SERVING int = 2
)

type Health interface {
	Check(ctx context.Context, service string) (int, error)
	Watch(req *hc.HealthCheckRequest, srv hc.Health_WatchServer) error
}

type HealthService struct {
}

func (HealthService) Check(_ context.Context, service string) (int, error) {
	//Service health checks here...
	switch service {
	case "grpc.health.v1.Health":
		return HealthCheckResponse_SERVING, nil

	default:
		return HealthCheckResponse_UNKNOWN, nil
	}

}

func (HealthService) Watch(req *hc.HealthCheckRequest, srv hc.Health_WatchServer) error {
	return nil
}
