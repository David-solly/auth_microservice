package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	grpcClient "github.com/David-solly/auth_microservice/cmd/client/client"
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	var (
		grpcAddr = flag.String("addr", ":8081",
			"gRPC address")
	)
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))

	if err != nil {
		log.Fatalln("gRPC dial:", err)
	}
	defer conn.Close()

	tokenService := grpcClient.New(conn)
	args := flag.Args()
	var cmd string
	cmd, args = pop(args)

	switch cmd {
	case "Generate":
		claims := make(map[string]string)

		for r := 0; r <= len(args); r++ {
			var claim, value string
			print(r)
			claim, args = pop(args)
			value, args = pop(args)
			claims[claim] = value
		}

		lorem(ctx, tokenService, claims)
	default:
		log.Fatalln("unknown command", cmd)
	}
}

// parse command line argument one by one
func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return s[0], s[1:]
}

// call lorem service
func lorem(ctx context.Context, service token_grpc.TokenServiceInterface, claims map[string]string) {
	mesg, err := service.Generate(ctx, claims)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(mesg)
}
