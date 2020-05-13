package main

import (
	"fmt"

	pb "github.com/David-solly/auth_microservice/pkg/api/v1"
	pbs "github.com/David-solly/auth_microservice/pkg/api/v1/service"
)

func main() {
	k := &pbs.TokenServiceEndpoints{}
	l := pb.ServiceError{}
	fmt.Printf("Hello world server %v %v", k, l)
}
