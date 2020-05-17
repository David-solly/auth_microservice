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
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Server could not load environmental variables")
	}

	token_grpc.RedisInit()

	print("Starting grpc server...")

	var (
		gRPCAddr = flag.String("grpc", ":8083",
			"gRPC listen address")
	)

	flag.Parse()
	var (
		gRPCAddrHealth = strings.Replace(*gRPCAddr, ":8083", ":8084", 1)
	)
	flag.Parse()
	ctx := context.Background()

	// init lorem service
	var svc token_grpc.TokenServiceInterface
	svc = token_grpc.TokenService{}

	errChan := make(chan error)

	// creating Endpoints struct
	endpoints := token_grpc.TokenServiceEndpoints{
		Endpoint: token_grpc.MakeTokenServiceEndpoint(svc),
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
		println("started on")
		fmt.Printf("%v", gRPCServer.GetServiceInfo())
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
	notifyOnStart()
	fmt.Println(<-errChan)
}

func notifyOnStart() {

	///braodcast start command
	http.PostForm("localhost:8080/v1/discover/service/BDg7tZZ2WhKPYJgsBeCgBokhUDshPcwNG1P0seddkHnsTbVB4iCTjxctoUjmQ8F1Dg3xCqgewnP5PbGmXMs4zrVvpVFYQQR43wHb", url.Values{})

}
