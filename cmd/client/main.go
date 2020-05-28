package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	ilog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcClient "github.com/David-solly/auth_microservice/cmd/client/client"
	models "github.com/David-solly/auth_microservice/pkg/api/v1/models"
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	ht "github.com/go-kit/kit/transport/http"
	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
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
		//fmt.Print"no cmd args")
	}

	switch cmd {
	case "Generate":
		claims := make(map[string]string)
		var claim, value string

		for r := 0; r <= len(args); r++ {
			var claim, value string
			//fmt.Printr)
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
func generateToken(ctx context.Context, service token_grpc.TokenServiceInterface, claims map[string]string) (*models.AccessTokens, error) {
	mesg, err := service.Generate(ctx, claims)
	if err != nil {
		ilog.Fatalln(err.Error())
		return nil, err
	}
	return mesg, nil
}

// call renewToken service
func renewToken(ctx context.Context, service token_grpc.TokenServiceInterface, token token_grpc.TokenRenewRequest) (*token_grpc.TokenResponse, error) {
	mesg, err := service.RenewTokens(ctx, token)
	if err != nil {
		ilog.Print(err.Error())
		return nil, err
	}
	return mesg, nil
}

// call generateToken service
func affectToken(ctx context.Context, service token_grpc.TokenServiceInterface, req models.TokenAffectRequest) (*models.TokenAffectResponse, error) {
	mesg, err := service.AffectToken(ctx, req)
	if err != nil {
		ilog.Fatalln(err.Error())
		return nil, err
	}
	if req.DesiredState == models.TokenState_LOGOUT {
		mesg.Message = "Log-out was successful."
	}
	return mesg, nil
}

// call generateToken service
func verifyToken(ctx context.Context, service token_grpc.TokenServiceInterface, tokenToverify token_grpc.TokenVerifyRequest) (*models.TokenVerifyResponse, interface{}) {
	mesg, err := service.VerifyToken(ctx, tokenToverify)
	if err != nil {
		if ob, k := err.(models.ResponseObject); k {
			return nil, ob
		}
		if errSt, k := err.(*string); k {
			return nil, &models.ResponseObject{Error: *errSt, Code: 401}
		}

		return nil, fmt.Errorf("Unknown error in verify with object %v", err)
	}
	return mesg, nil
}

//parse consul message to endpoint
func generateTokenFactory(_ context.Context, method, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		// ####### ----START----for http #############
		// if !strings.HasPrefix(instance, "http") {
		// 	instance = "http://" + instance
		// }

		fmt.Println("@@@@@@@@ received from consul")
		fmt.Println(instance)
		fmt.Println(method)
		fmt.Println(path)

		// tgt, err := url.Parse(instance)
		// if err != nil {
		// 	return nil, nil, err
		// }
		// tgt.Path = path

		// var (
		// 	enc ht.EncodeRequestFunc
		// 	dec ht.DecodeResponseFunc
		// )
		// enc, dec = encodeLoremRequest, decodeLoremResponse
		//return ht.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
		//######### ----END---- for http endpoint

		conn1, err := dialConnection(&instance)
		if err != nil {
			return nil, nil, err
		}

		svc := makeConnection(conn1)
		return makeBalancedGenerateEndpoint(svc), conn1, nil
	}
}

//parse consule message to endpoint
func verifyTokenFactory(_ context.Context, method, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		fmt.Println("@@@@@@@@ received from consul")
		fmt.Println(instance)
		fmt.Println(method)
		fmt.Println(path)

		conn1, err := dialConnection(&instance)
		if err != nil {
			return nil, nil, err
		}

		svc := makeConnection(conn1)
		return makeBalancedVerifyEndpoint(svc), conn1, nil
	}
}

func affectTokenFactory(_ context.Context, method, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		fmt.Println("@@@@@@@@ received from consul")
		fmt.Println(instance)
		fmt.Println(method)
		fmt.Println(path)

		conn1, err := dialConnection(&instance)
		if err != nil {
			return nil, nil, err
		}

		svc := makeConnection(conn1)
		return makeBalancedAffectEndpoint(svc), conn1, nil
	}
}

func renewTokenFactory(_ context.Context, method, path string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		fmt.Println("@@@@@@@@ received from consul")
		fmt.Println(instance)
		fmt.Println(method)
		fmt.Println(path)

		conn1, err := dialConnection(&instance)
		if err != nil {
			return nil, nil, err
		}

		svc := makeConnection(conn1)
		return makeBalancedRenewEndpoint(svc), conn1, nil
	}
}

func main() {
	// TODO:
	//Remove for production, already loads on a different flow
	//####################
	if err := godotenv.Load("../../.env"); err != nil {
		ilog.Fatalln("Server could not load environmental variables")
	}

	port = os.Getenv("PORT")
	//####################

	servieURL = os.Getenv("SERVICE_URL")
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

	ctx := context.Background()

	tags := []string{"jwt generate", "yea-buddy!!!"}
	passingOnly := true
	duration := 500 * time.Millisecond
	var generateEndpoint endpoint.Endpoint
	var verifyEndpoint endpoint.Endpoint
	var affectEndpoint endpoint.Endpoint
	var renewEndpoint endpoint.Endpoint

	serviceName := "JWT-Service"
	instancer := consulsd.NewInstancer(client, logger, serviceName, tags, passingOnly)
	{
		factory := affectTokenFactory(ctx, "POST", "/affect")
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, duration, balancer)
		affectEndpoint = retry
	}
	{
		factory := generateTokenFactory(ctx, "POST", "/login")
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, duration, balancer)
		generateEndpoint = retry
	}
	{
		factory := verifyTokenFactory(ctx, "POST", "/verify")
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(1, duration, balancer)
		verifyEndpoint = retry
	}

	{
		factory := renewTokenFactory(ctx, "POST", "/renew")
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(1, duration, balancer)
		renewEndpoint = retry
	}

	loginHandle := ht.NewServer(
		generateEndpoint,
		loginHandler,
		encodeResponse,
	)

	logoutHandle := ht.NewServer(
		affectEndpoint,
		logoutHandler,
		encodeAffectResponse,
	)

	registerhandle := ht.NewServer(
		generateEndpoint,
		registerHandler,
		encodeResponse,
	)

	greethandle := ht.NewServer(
		generateEndpoint,
		greetingHandler,
		encodeResponse,
	)

	verifyhandle := ht.NewServer(
		verifyEndpoint,
		verifyHandler,
		encodeVerifyResponse,
	)
	renewhandle := ht.NewServer(
		renewEndpoint,
		refreshHandler,
		encodeRenewResponse,
	)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set request timeout
	r.Use(middleware.Timeout(500 * time.Millisecond))

	// Api endpoints
	r.Handle("/", greethandle)
	r.Handle("/v1/login", loginHandle)
	r.Handle("/v1/logout", logoutHandle)
	r.Handle("/v1/refresh", renewhandle)
	r.Handle("/v1/register", registerhandle)
	r.Handle("/v1/verify/{serviceID}", verifyhandle)

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
		req, k := request.(token_grpc.TokenRequest)
		if !k {
			if resp, ok := request.(models.ResponseObject); ok {

				return resp, nil
			}
		}
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

// Endpoints are a primary abstraction in go-kit. An endpoint represents a single RPC (method in our service interface)
func makeBalancedRenewEndpoint(svc token_grpc.TokenServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, k := request.(token_grpc.TokenRenewRequest)
		if !k {
			if resp, ok := request.(models.ResponseObject); ok {

				return resp, nil
			}
		}
		if req.RefreshToken == "" {
			return nil, errors.New("No Refresh Token provided to encode")
		}

		v, err := renewToken(ctx, svc, req)
		if err != nil {
			return v, err
		}
		return v, nil
	}
}

// Endpoints are a primary abstraction in go-kit. An endpoint represents a single RPC (method in our service interface)
func makeBalancedAffectEndpoint(svc token_grpc.TokenServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, k := request.(models.TokenAffectRequest)
		if !k {
			if resp, ok := request.(models.ResponseObject); ok {

				return resp, nil
			}
		}
		if req.DesiredState > 5 {
			return nil, errors.New("Unknown method")
		}

		v, err := affectToken(ctx, svc, req)
		if err != nil {
			return v, err
		}
		return v, nil
	}
}

// Endpoints are a primary abstraction in go-kit. An endpoint represents a single RPC (method in our service interface)
func makeBalancedVerifyEndpoint(svc token_grpc.TokenServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, k := request.(token_grpc.TokenVerifyRequest)

		if !k {
			if resp, ok := request.(models.ResponseObject); ok {
				return resp, nil
			}
		}

		if req.Token == "" {
			return nil, errors.New("No Token to verify- please supply a valid token")
		}
		v, err := verifyToken(ctx, svc, req)
		if resp, ok := err.(models.ResponseObject); ok {
			return resp, nil
		}

		if err != nil {
			return nil, fmt.Errorf("Error in make balanced verify : %v", err)
		}
		return v, nil
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	if r, o := response.(models.ResponseObject); o {
		if r.Error != "" {
			k, _ := json.Marshal(r)

			w.WriteHeader(r.Code)
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Content-Length", fmt.Sprintf("%d", len(k)))
			return json.NewEncoder(w).Encode(r)
		}

		if r, o := response.(models.AccessTokens); o {
			k := map[string]interface{}{
				"tokens": r,
			}
			return json.NewEncoder(w).Encode(k)

		}

		if r, o := response.(*models.AccessTokens); o {
			k := map[string]interface{}{
				"tokens": r,
			}
			return json.NewEncoder(w).Encode(k)

		}

		return json.NewEncoder(w).Encode(r)

	}
	return json.NewEncoder(w).Encode(response)
}

func encodeRenewResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	if r, o := response.(*token_grpc.TokenResponse); o {
		if r.Error.Code != 0 {
			return json.NewEncoder(w).Encode(r.Error)
		}
		if r.Response.RefreshToken == "" {
			e := models.ResponseObject{Error: "Invalid Request token", Code: http.StatusBadRequest}
			k, _ := json.Marshal(e)

			w.WriteHeader(e.Code)
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Content-Length", fmt.Sprintf("%d", len(k)))
			w.Write(k)
			return json.NewEncoder(w).Encode(e)
		}

		k := map[string]interface{}{
			"tokens": r.Response,
		}
		return json.NewEncoder(w).Encode(k)

	}
	return json.NewEncoder(w).Encode(response)
}

func encodeVerifyResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	if r, o := response.(*models.TokenVerifyResponse); o {
		if r.Error.Code != 0 {
			k, _ := json.Marshal(r.Error)

			w.WriteHeader(r.Error.Code)
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Content-Length", fmt.Sprintf("%d", len(k)))
			return json.NewEncoder(w).Encode(r.Error)
		}

		return json.NewEncoder(w).Encode(r.Access)

	}
	return json.NewEncoder(w).Encode(response)
}

func encodeAffectResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	if r, o := response.(*models.TokenAffectResponse); o {
		if r.Error != nil {
			k, _ := json.Marshal(r.Error)

			w.WriteHeader(r.Error.Code)
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Content-Length", fmt.Sprintf("%d", len(k)))
			return json.NewEncoder(w).Encode(r.Error)
		}

		return json.NewEncoder(w).Encode(r)

	}
	return json.NewEncoder(w).Encode(response)
}

// encode error
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func genError(errorString string, errorCode int) *[]byte {
	v, _ := json.Marshal(models.ResponseObject{Error: errorString, Code: errorCode})
	return &v
}

/*
run command
go run . -consul.addr localhost -consul.port 8500

//debug args
  "args": ["-consul.addr", "localhost", "-consul.port", "8500"]

*/
