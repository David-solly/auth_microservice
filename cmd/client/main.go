package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	ilog "log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcClient "github.com/David-solly/auth_microservice/cmd/client/client"
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	ht "github.com/go-kit/kit/transport/http"

	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func makeConnection(conn *grpc.ClientConn) token_grpc.TokenServiceInterface {

	tokenService := grpcClient.New(conn)
	return tokenService
}

func dialConnection(grpcAddr *string) (*grpc.ClientConn, error) {

	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))

	if err != nil {
		ilog.Fatalln("gRPC dial:", err)
		return nil, err
	}
	return conn, nil
}

func runCmd(add *string) error {
	ctx := context.Background()
	conn, err := dialConnection(add)
	if err != nil {
		return err
	}
	tokenService := makeConnection(conn)
	//defer conn.Close()
	args := flag.Args()
	var cmd string
	cmd, args = pop(args)

	if len(args) < 2 {
		print("no cmd args")
	}

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
		ilog.Fatalln("unknown command", cmd)
	}
	return nil
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
		ilog.Fatalln(err.Error())
		return nil, err
	}
	// fmt.Println(mesg)
	return mesg, nil
}

//parse consule message to endpoint
func generateTokenFactory(_ context.Context, method, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		// if !strings.HasPrefix(instance, "http") {
		// 	instance = "http://" + instance
		// }

		fmt.Println("@@@@@@@@ received from consul")
		fmt.Println(instance)
		fmt.Println(method)
		fmt.Println(path)

		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path
		conn1, err := dialConnection(&instance)
		svc := makeConnection(conn1)
		conn = conn1
		// var (
		// 	enc ht.EncodeRequestFunc
		// 	dec ht.DecodeResponseFunc
		// )
		// enc, dec = encodeLoremRequest, decodeLoremResponse

		//return ht.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil

		return makeBalancedGenerateEndpoint(svc), conn1, nil
	}
}

func main() {
	// var (
	// 	grpcAddr = flag.String("addr", ":8081",
	// 		"gRPC address")
	// )
	// flag.Parse()

	// //go runCmd(grpcAddr)
	// rest(grpcAddr)
	var (
		consulAddr = flag.String("consul.addr", "", "consul address")
		consulPort = flag.String("consul.port", "", "consul port")
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()

		consulConfig.Address = "http://" + *consulAddr + ":" + *consulPort
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	tags := []string{"jwt generate", "yea-buddy!!!"}
	passingOnly := true
	duration := 500 * time.Millisecond
	var generateEndpoint endpoint.Endpoint

	ctx := context.Background()
	r := chi.NewRouter()

	factory := generateTokenFactory(ctx, "POST", "/login")
	serviceName := "JWT-Service"
	instancer := consulsd.NewInstancer(client, logger, serviceName, tags, passingOnly)
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(1, duration, balancer)
	generateEndpoint = retry

	// svc := token_grpc.TokenService{}

	// signinHandler := ht.NewServer(
	// 	makeUppercaseEndpoint(svc),
	// 	decodeUppercaseRequest,
	// 	encodeResponse,
	// )

	suhandle := ht.NewServer(
		generateEndpoint,
		decodeUppercaseRequest,
		encodeResponse,
	)

	// r.Handle("/login", signinHandler)
	r.Handle("/login", suhandle)

	// Interrupt handler.
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport.
	go func() {
		logger.Log("transport", "HTTP", "addr", "8080")
		errc <- http.ListenAndServe(":8080", r)
	}()

	// Run!
	logger.Log("exit", <-errc)

}

// Endpoints are a primary abstraction in go-kit. An endpoint represents a single RPC (method in our service interface)
func makeBalancedGenerateEndpoint(svc token_grpc.TokenServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(token_grpc.TokenRequest)
		if len(req.Claims) < 1 {
			return nil, errors.New("No claims to encode")
		}
		v, err := generateToken(ctx, svc, req.Claims)
		if err != nil {
			return v, err
		}
		return v, nil
	}
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	user := User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		j, _ := json.Marshal(ResponseObject{Error: "Bad request fromat", Code: http.StatusBadRequest})
		return nil, errors.New(string(j))
	}

	if user.Username == tUser.Username && user.Password == tUser.Password {
		if conn == nil {
			j, _ := json.Marshal(ResponseObject{Error: "Sorry, could not process request at this time. Please try again later", Code: http.StatusInternalServerError})
			return errors.New(string(j)), nil
		}

		request := token_grpc.TokenRequest{Claims: make(map[string]string)}
		request.Claims["email"] = user.Username
		request.Claims["id"] = fmt.Sprintf("%v", tUser.ID)

		return request, nil

		//reference code

		// var req token_grpc.TokenRequest2
		// request := token_grpc.TokenRequest{Claims: make(map[string]string)}
		// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// 	return nil, err
		// }

		// for x := 0; x < len(req.Claims); x += 2 {
		// 	request.Claims[req.Claims[x]] = req.Claims[x+1]
		// }

		// return request, nil

	}

	j, _ := json.Marshal(ResponseObject{Error: "Sorry, the login credentials don't match any records", Code: http.StatusNoContent})
	return nil, errors.New(string(j))

}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error.
type errorer interface {
	error() error
}

// encode error
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func genError(errorString string, errorCode int) *[]byte {
	v, _ := json.Marshal(ResponseObject{Error: errorString, Code: errorCode})
	return &v
}

/*
run command
go run . -consul.addr localhost -consul.port 8500

*/
