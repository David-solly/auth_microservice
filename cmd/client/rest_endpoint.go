package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/David-solly/auth_microservice/pkg/api/v1/models"
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	"github.com/go-chi/chi"
)

var (
	router = chi.NewRouter()
	port   = os.Getenv("PORT")
)

var (
	tUser = models.User{
		ID:       1,
		Username: "binyamin.dev@gmail.com",
		Password: "binyamin.dev@gmail.com",
	}
)

func greetingHandler(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodGet {
		return nil, http.ErrBodyNotAllowed
	}
	return models.ResponseObject{Code: http.StatusOK, Message: "Api is up"}, nil
}

func verifyHandler(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPost {
		return nil, http.ErrBodyNotAllowed
	}
	token, ok := extractAuthToken(r)
	if !ok {
		return models.ResponseObject{Error: token, Code: http.StatusUnauthorized}, nil

	}
	serviceID := chi.URLParam(r, "serviceID")
	return token_grpc.TokenVerifyRequest{Token: token, Service: serviceID}, nil

}

func logoutHandler(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPost {
		return nil, http.ErrBodyNotAllowed
	}
	token, ok := extractAuthToken(r)
	if !ok {
		return models.ResponseObject{Error: token, Code: http.StatusUnauthorized}, nil

	}
	return models.TokenAffectRequest{Token: token, DesiredState: models.TokenState_LOGOUT}, nil

}

func errorHandler(w http.ResponseWriter, response models.ResponseObject) {
	err, _ := json.Marshal(response)
	w.Write(err)
}

func loginHandler(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPost {
		return nil, http.ErrBodyNotAllowed
	}
	user := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return models.ResponseObject{Error: "Bad request fromat", Code: http.StatusBadRequest}, nil
	}

	if user.Username == tUser.Username && user.Password == tUser.Password {

		claims := make(map[string]string)
		claims["id"] = fmt.Sprintf("%v", tUser.ID)
		return token_grpc.TokenRequest{Claims: claims}, nil

	}

	return models.ResponseObject{Error: "Sorry, the login credentials don't match any records", Code: http.StatusNoContent}, nil

}

func registerHandler(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPost {
		return nil, http.ErrBodyNotAllowed
	}
	user := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return models.ResponseObject{Error: "Bad request fromat", Code: http.StatusBadRequest}, nil

	}

	if user.PII.Mobile == "" {
		return models.ResponseObject{Error: "Sorry, please include a valid mobile number", Code: http.StatusUnprocessableEntity}, nil
	}

	if unique, e := verifyUniqueUser(user.Username); e != nil {
		return models.ResponseObject{Error: "Sorry, could not process request", Code: http.StatusInternalServerError}, nil

	} else if !unique {
		return models.ResponseObject{Error: "Sorry, the username already exist on our system.", Code: http.StatusAlreadyReported}, nil
	}

	claims := user.MapClaims()
	claims["id"] = fmt.Sprintf("%v", 2)
	return token_grpc.TokenRequest{Claims: claims}, nil

}

func verifyUniqueUser(username string) (bool, error) {
	return username != tUser.Username, nil
}

func extractAuthToken(r *http.Request) (string, bool) {
	token := r.Header.Get("AUTH_API_KEY")

	if token == "" {
		return "Unauthorized access, no api key found with your request", false
	}

	return token, true
}
