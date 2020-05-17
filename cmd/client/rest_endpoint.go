package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// ResponseObject ...
// general response object
type ResponseObject struct {
	Error   string                   `json:"error,omitempty"`
	Code    int                      `json:"code,omitempty"`
	Tokens  *token_grpc.AccessTokens `json:"tokens,omitempty"`
	Token   string                   `json:"token,omitempty"`
	Message string                   `json:"message,omitempty"`
}

type User struct {
	ID       uint64  `json:"id"`
	Username string  `json:"email"`
	Password string  `json:"password"`
	PII      UserPII `json:"details,omitempty"`
}

type UserPII struct {
	Mobile    string `json:"mobile,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type ServiceRequest struct {
	UserID  uint64      `json:"user_id"`
	Payload interface{} `json:"payload,omitempty"`
}

var (
	router = chi.NewRouter()
	port   = os.Getenv("PORT")
)

var (
	tUser = User{
		ID:       1,
		Username: "binyamin.dev@gmail.com",
		Password: "binyamin.dev@gmail.com",
	}
)

var address *string

var conn *grpc.ClientConn

func initGrpc() {

	conn1, err := dialConnection(address)
	if err != nil {
		// errorHandler(w, ResponseObject{Error: "Sorry, could not process request at this time", Code: http.StatusInternalServerError})
		log.Panicln("Failed to connect to Token service instance from api gateway")
		return
	}

	if conn1 == nil {
		log.Panicf("Error etablishing initial connection to grpc server : %v", *address)
		return
	}

	conn = conn1
	fmt.Printf("Succefully connected to grpc server: %v", conn)
}

func rest(grpcAddr *string) {
	address = grpcAddr
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("\nRecovered From [", r, "] \n Running a limited service!!!")
			}
		}()
		initGrpc()
	}()

	// TODO:
	//Remove for production, already loads on a different flow
	//####################
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Server could not load environmental variables")
	}

	port = os.Getenv("PORT")
	//####################
	fmt.Printf("auth micro\nServer is running on port %v", port)
	buildRoutes(router, port)
}

func buildRoutes(r *chi.Mux, port string) {

	go RedisInit()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set request timeout
	r.Use(middleware.Timeout(10 * time.Second))

	// Api endpoints
	r.Get("/", greetingHandler)

	r.Post("/login", loginHandler)

	r.Post("/register", registerHandler)

	r.Post("/service/{serviceID}", serviceHandler)

	r.Post("/v1/discover/service/{serviceID}", serviceConnectionHandler)

	http.ListenAndServe(":"+port, r)
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	greet, _ := json.Marshal(ResponseObject{Code: http.StatusOK, Message: "Api is up"})
	w.Write(greet)
}

func serviceConnectionHandler(w http.ResponseWriter, r *http.Request) {

	serviceID := chi.URLParam(r, "serviceID")
	if serviceID == os.Getenv("GRPC_ONLINE_CODE") {
		if conn == nil {
			go func(w http.ResponseWriter) {
				defer func(w http.ResponseWriter) {

					if r := recover(); r != nil {
						fmt.Println("\nRecovered From [", r, "] \n Running a limited service!!!")
					}
				}(w)
				initGrpc()
				confirmService, _ := json.Marshal(ResponseObject{Code: http.StatusOK, Message: fmt.Sprintf("Service %v is ... restarting", serviceID)})
				w.Write(confirmService)
			}(w)
		}
	}

	if serviceID == os.Getenv("GRPC_OFFLINE_CODE") {
		if conn != nil {
			if conn.GetState() == connectivity.Idle || conn.GetState() == connectivity.Connecting || conn.GetState() == connectivity.Ready {
				conn.Close()
				conn = nil
				confirmService, _ := json.Marshal(ResponseObject{Code: http.StatusOK, Message: fmt.Sprintf("Service %v is disconnecting", serviceID)})
				w.Write(confirmService)

			}
		}
	}

}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	token, ok := extractAuthToken(r)
	if !ok {
		errorHandler(w, ResponseObject{Error: token, Code: http.StatusUnauthorized})
		return
	}
	tokenAuth, err := ExtractTokenMetadata(token)
	if err != nil {
		errorHandler(w, ResponseObject{Error: "Unauthorized token", Code: http.StatusUnauthorized})
		return
	}
	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		errorHandler(w, ResponseObject{Error: "Unauthorized for resource", Code: http.StatusUnauthorized})
		return
	}
	td := ServiceRequest{}
	td.UserID = userId

	serviceID := chi.URLParam(r, "serviceID")
	confirmService, _ := json.Marshal(ResponseObject{Code: http.StatusOK, Message: fmt.Sprintf("Service %v is up - authorized for id :%s", serviceID, td.UserID)})
	w.Write(confirmService)
}

func errorHandler(w http.ResponseWriter, response ResponseObject) {
	err, _ := json.Marshal(response)
	w.Write(err)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user := User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errorHandler(w, ResponseObject{Error: "Bad request fromat", Code: http.StatusBadRequest})
		return
	}

	if user.Username == tUser.Username && user.Password == tUser.Password {
		if conn == nil {
			errorHandler(w, ResponseObject{Error: "Sorry, could not process request at this time. Please try again later", Code: http.StatusInternalServerError})
			return
		}

		if conn.GetState() != connectivity.Ready {
			errorHandler(w, ResponseObject{Error: "Sorry, could not process request at this time", Code: http.StatusInternalServerError})
			return
		}

		tokenService := makeConnection(conn)

		claims := make(map[string]string)
		claims["id"] = fmt.Sprintf("%v", tUser.ID)

		tokens, err := generateToken(context.Background(), tokenService, claims)
		if err != nil {
			errorHandler(w, ResponseObject{Error: "Sorry, could not process request", Code: http.StatusUnprocessableEntity})
			return
		}

		_ = json.NewEncoder(w).Encode(ResponseObject{Tokens: tokens})
		return
	}

	errorHandler(w, ResponseObject{Error: "Sorry, the login credentials don't match any records", Code: http.StatusNoContent})
	return

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	user := User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errorHandler(w, ResponseObject{Error: "Bad request fromat", Code: http.StatusBadRequest})
		return
	}

	if user.PII.Mobile == "" {
		errorHandler(w, ResponseObject{Error: "Sorry, please include a valid mobile number", Code: http.StatusUnprocessableEntity})
		return
	}

	if unique, e := verifyUniqueUser(user.Username); e != nil {
		errorHandler(w, ResponseObject{Error: "Sorry, could not process request", Code: http.StatusInternalServerError})
		return
	} else if !unique {
		errorHandler(w, ResponseObject{Error: "Sorry, the username already exist on our system.", Code: http.StatusAlreadyReported})
		return

	} else {
		if conn == nil {
			errorHandler(w, ResponseObject{Error: "Sorry, could not process request at this time. Please try again later", Code: http.StatusInternalServerError})
			return
		}

		if conn.GetState() != connectivity.Ready {
			errorHandler(w, ResponseObject{Error: "Sorry, could not process request at this time", Code: http.StatusInternalServerError})
			return
		}
		tokenService := makeConnection(conn)

		claims := make(map[string]string)
		claims["id"] = fmt.Sprintf("%v", 2)

		tokens, err := generateToken(context.Background(), tokenService, claims)
		if err != nil {
			errorHandler(w, ResponseObject{Error: "Sorry, could not process request", Code: http.StatusUnprocessableEntity})
			return
		}
		_ = json.NewEncoder(w).Encode(ResponseObject{Tokens: tokens})
		return
	}

}

func verifyUniqueUser(username string) (bool, error) {
	return username != tUser.Username, nil
}

func VerifyTokenIntegrity(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Token signing error, unexpected method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(tokenString string) error {
	token, err := VerifyTokenIntegrity(tokenString)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func extractAuthToken(r *http.Request) (string, bool) {
	token := r.Header.Get("AUTH_API_KEY")
	if token == "" {
		return "Unauthorized access, no api key found with your request", false
	}
	return token, true

}

var client *redis.Client

func RedisInit() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis server - Online ..........")
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

func FetchAuth(authD *AccessDetails) (uint64, error) {
	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func ExtractTokenMetadata(tokenString string) (*AccessDetails, error) {
	token, err := VerifyTokenIntegrity(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		claimID := claims["id"]
		userId, err := strconv.ParseUint(claimID.(string), 10, 64)
		if err != nil {
			return nil, err
		}

		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}
