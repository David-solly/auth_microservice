package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strings"

	"os"
	"os/signal"

	"syscall"

	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	"github.com/David-solly/consul_hcsd/discover"
	"github.com/David-solly/consul_hcsd/discover/models"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
)

func main() {
	// #### For debugging and local testing
	// if err := godotenv.Load("../../.env"); err != nil {
	// 	log.Fatal("Server could not load environmental variables")
	// }

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	errChan := make(chan error)

	if _, err := token_grpc.RedisInit(); err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	//fmt.Print"Starting grpc server...")

	// parse variable from input command
	var (
		consulAddr          = flag.String("consul.addr", "", "consul address")
		consulPort          = flag.String("consul.port", "", "consul port")
		advertiseAddr       = flag.String("advertise.addr", "", "advertise address")
		advertisePort       = flag.String("advertise.port", "", "advertise port")
		advertiseHealthPort = flag.String("health.port", "", "health port")
	)
	flag.Parse()

	if envAdd := os.Getenv("POD_IP"); envAdd != "" {
		repl := strings.ReplaceAll(envAdd, ".", "-")
		if advertiseAddr != nil {
			newAdd := fmt.Sprintf("%s.%s", repl, *advertiseAddr)
			advertiseAddr = &newAdd
		} else {
			advertiseAddr = &envAdd
		}

	}

	var (
		gRPCAddr = flag.String("grpc", ""+*advertiseAddr+":"+*advertisePort,
			"gRPC listen address")
	)
	flag.Parse()

	ctx := context.Background()

	// init lorem service
	var svc token_grpc.TokenServiceInterface
	svc = token_grpc.TokenService{}

	// creating Endpoints struct
	endpoints := token_grpc.TokenServiceEndpoints{
		GenerateEndpoint: token_grpc.MakeTokenServiceGenerateEndpoint(svc),
		VerifyEndpoint:   token_grpc.MakeTokenServiceVerifyEndpoint(svc),
		AffectEndpoint:   token_grpc.MakeTokenServicAffectEndpoint(svc),
		RenewEndpoint:    token_grpc.MakeTokenServiceRenewEndpoint(svc),
	}

	//execute grpc server
	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}

		handler := token_grpc.NewGRPCServer(ctx, endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterTokenServiceServer(gRPCServer, handler)
		fmt.Printf("Service info %v", gRPCServer.GetServiceInfo())
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// Register Service to Consul
	discover.ConfigureAndAdvertise(
		&models.AddressConfig{
			ConsulAddr:          *consulAddr,
			ConsulPort:          *consulPort,
			AdvertiseAddr:       *advertiseAddr,
			AdvertisePort:       *advertisePort,
			AdvertiseHealthPort: *advertiseHealthPort},
		&models.ServiceConfig{
			ID:   "JWT",
			Name: "JWT-Service",
			Tags: []string{"jwt", "generate", "refresh", "verify"},
		})

	//notifyOnStart()
	error := <-errChan
	logger.Log("err", error)
	discover.ErrChanHC <- error
	// deregister service

}

/*
run command
go run . -consul.addr localhost -consul.port 8500 -advertise.addr localhost -advertise.port 8081 -health.port 8082

//debuggings args
"args": ["-consul.addr", "localhost", "-consul.port", "8500","-advertise.addr", "localhost","-advertise.port", "8081","-health.port", "8082"]

*/
