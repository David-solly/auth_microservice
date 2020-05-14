package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"

)

// ResponseObject ...
// general response object
type ResponseObject struct {
	Error   string        `json:"error,omitempty"`
	Code    int           `json:"code,omitempty"`
	Tokens  *token_grpc.AccessTokens `json:"tokens,omitempty"`
	Token   string        `json:"token,omitempty"`
	Message string        `json:"message,omitempty"`
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
var address *string;

func rest(grpcAddr *string) {
	address=grpcAddr
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

	http.ListenAndServe(":"+port, r)
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	greet, _ := json.Marshal(ResponseObject{Code: http.StatusOK, Message: "Api is up"})
	w.Write(greet)
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
		ctx := context.Background()
		conn:=dialConnection(address)
		tokenService := makeConnection(conn)
		defer conn.Close()
		claims:=make(map[string]string)
		claims["id"]=fmt.Sprintf("%v",tUser.ID)

		tokens, err := generateToken(ctx,tokenService,claims)
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
		ctx := context.Background()
	conn := dialConnection(address)
	tokenService := makeConnection(conn)
	defer conn.Close()
		claims:=make(map[string]string)
		claims["id"]=fmt.Sprintf("%v",tUser.ID)

		tokens, err := generateToken(ctx,tokenService,claims)
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

