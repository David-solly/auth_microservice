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

func makeConnection(conn *grpc.ClientConn) token_grpc.TokenServiceInterface {

	tokenService := grpcClient.New(conn)
	return tokenService
}

func dialConnection(grpcAddr *string) *grpc.ClientConn {

	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))

	if err != nil {
		log.Fatalln("gRPC dial:", err)
		return nil
	}
	return conn
}

func runCmd(add *string) {
	ctx := context.Background()
	conn := dialConnection(add)
	tokenService := makeConnection(conn)
	defer conn.Close()
	args := flag.Args()
	var cmd string
	cmd, args = pop(args)

	switch cmd {
	case "Generate":
		claims := make(map[string]string)
		var claim, value string

		for r := 0; r <= len(args); r++ {
			var claim, value string
			print(r)
			claim, args = pop(args)
			value, args = pop(args)
			claims[claim] = value
		}

		claim, args = pop(args)
		value, args = pop(args)
		claims[claim] = value

		generateToken(ctx, tokenService, claims)
	default:
		log.Fatalln("unknown command", cmd)
	}
}

func main() {
	var (
		grpcAddr = flag.String("addr", ":8081",
			"gRPC address")
	)
	flag.Parse()

	// go runCmd(grpcAddr)
	rest(grpcAddr)
}

// parse command line argument one by one
func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return s[0], s[1:]
}

// call generateToken service
func generateToken(ctx context.Context, service token_grpc.TokenServiceInterface, claims map[string]string) (*token_grpc.AccessTokens, error) {
	mesg, err := service.Generate(ctx, claims)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}
	fmt.Println(mesg)
	return mesg, nil
}
