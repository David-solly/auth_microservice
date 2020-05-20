package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
	hc "github.com/David-solly/auth_microservice/pkg/api/v1/hc"
	token_consul "github.com/David-solly/auth_microservice/pkg/api/v1/registration"
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Server could not load environmental variables")
	}

	token_grpc.RedisInit()

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

	var (
		gRPCAddr = flag.String("grpc", ""+*advertiseAddr+":"+*advertisePort,
			"gRPC listen address")
	)

	flag.Parse()
	var (
		gRPCAddrHealth = strings.Replace(*gRPCAddr, ":"+*advertisePort, ":"+*advertiseHealthPort, 1)
	)
	flag.Parse()

	// Register Service to Consul
	registar := token_consul.RegisterService(*consulAddr,
		*consulPort,
		*advertiseAddr,
		*advertisePort,
		*advertiseHealthPort)

	ctx := context.Background()

	// init lorem service
	var svc token_grpc.TokenServiceInterface
	svc = token_grpc.TokenService{}

	errChan := make(chan error)

	// creating Endpoints struct
	endpoints := token_grpc.TokenServiceEndpoints{
		GenerateEndpoint: token_grpc.MakeTokenServiceGenerateEndpoint(svc),
		VerifyEndpoint:   token_grpc.MakeTokenServiceVerifyEndpoint(svc),
	}

	var svcH token_grpc.Health
	svcH = token_grpc.HealthService{}
	// svcH = token_grpc.LoggingMiddlewareHealth(logger)(svcH)

	check := token_grpc.MakeHealthServiceCheckEndpoint(svcH)
	watch := token_grpc.MakeHealthServiceWatchEndpoint(svcH)

	// check = token_grpc.NewTokenBucketLimitter(rlbucket)(check)
	endpointsHealth := token_grpc.EndpointsConsul{
		ConsulHealthCheckEndpoint: check,
		ConsulHealthWatchEndpoint: watch,
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
		registar.Register()
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		fmt.Println("starting health server")
		listener, err := net.Listen("tcp", gRPCAddrHealth)
		if err != nil {
			errChan <- err
			return
		}
		handler := token_grpc.NewGRPCServerHealth(ctx, endpointsHealth)
		gRPCServer := grpc.NewServer()
		hc.RegisterHealthServer(gRPCServer, handler)
		fmt.Printf("Service info %v", gRPCServer.GetServiceInfo())
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	//notifyOnStart()
	error := <-errChan
	// deregister service
	registar.Deregister()
	log.Fatalln(error)

}

func notifyOnStart() {

	///braodcast start command
	http.PostForm("localhost:8080/v1/discover/service/BDg7tZZ2WhKPYJgsBeCgBokhUDshPcwNG1P0seddkHnsTbVB4iCTjxctoUjmQ8F1Dg3xCqgewnP5PbGmXMs4zrVvpVFYQQR43wHb", url.Values{})

}

/*
run command
go run . -consul.addr localhost -consul.port 8500 -advertise.addr localhost -advertise.port 8081 -health.port 8082

*/
