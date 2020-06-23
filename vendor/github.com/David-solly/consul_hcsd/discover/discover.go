package discover

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/David-solly/consul_hcsd/discover/models"
	hc "github.com/David-solly/consul_hcsd/hc"
	token_consul "github.com/David-solly/consul_hcsd/registration"
	health_grpc "github.com/David-solly/consul_hcsd/service"
	"google.golang.org/grpc"
)

// ErrChanHC : A channel to pass stop signal into from caller
var ErrChanHC = make(chan error)

// ConfigureAndAdvertise : The Api to create a new service health check and
// register the service with consul sd
func ConfigureAndAdvertise(address *models.AddressConfig, service *models.ServiceConfig) {
	print("Starting configuration -- announcing discovery")

	var (
		gRPCAddr = flag.String("grpc-health", ""+address.AdvertiseAddr+":"+address.AdvertiseHealthPort,
			"gRPC listen address")
	)
	flag.Parse()

	// Register Service to Consul
	registar := token_consul.RegisterService(address.ConsulAddr,
		address.ConsulPort,
		address.AdvertiseAddr,
		address.AdvertisePort,
		address.AdvertiseHealthPort,
		service.ID, service.Name, service.Tags)

	ctx := context.Background()

	var svcH health_grpc.Health
	svcH = health_grpc.HealthService{}
	// svcH = health_grpc.LoggingMiddlewareHealth(logger)(svcH)

	check := health_grpc.MakeHealthServiceCheckEndpoint(svcH)
	watch := health_grpc.MakeHealthServiceWatchEndpoint(svcH)

	// check = health_grpc.NewTokenBucketLimitter(rlbucket)(check)
	endpointsHealth := health_grpc.EndpointsConsul{
		ConsulHealthCheckEndpoint: check,
		ConsulHealthWatchEndpoint: watch,
	}

	go func() {
		fmt.Println("starting healthcheck grpc server")
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			ErrChanHC <- err
			return
		}
		handler := health_grpc.NewGRPCServerHealth(ctx, endpointsHealth)
		gRPCServer := grpc.NewServer()
		hc.RegisterHealthServer(gRPCServer, handler)
		registar.Register()
		fmt.Printf("Service info %v", gRPCServer.GetServiceInfo())
		ErrChanHC <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		ErrChanHC <- fmt.Errorf("%s", <-c)
	}()

	error := <-ErrChanHC
	// deregister service on exit
	registar.Deregister()

	log.Fatalln(error)

}

//go run . -consul.addr localhost -consul.port 8500 -advertise.addr localhost -advertise.port 8081 -health.port 8082
