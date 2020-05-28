package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/David-solly/auth_microservice/pkg/api/v1/models"
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
	"github.com/go-chi/chi"
)

var (
	router    = chi.NewRouter()
	port      = os.Getenv("PORT")
	servieURL string
)

func greetingHandler(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodGet {
		return nil, http.ErrBodyNotAllowed
	}
	return models.ResponseObject{Code: http.StatusOK, Message: "Api is up"}, nil
}

func verifyHandler(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPost {
		return nil, http.ErrNotSupported
	}
	token, ok := extractAuthToken(r)
	if !ok {
		return models.ResponseObject{Error: token, Code: http.StatusUnauthorized}, nil

	}
	serviceID := chi.URLParam(r, "serviceID")
	return token_grpc.TokenVerifyRequest{Token: token, Service: serviceID}, nil

}

func refreshHandler(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodGet {
		return nil, http.ErrBodyNotAllowed
	}
	token, ok := extractAuthToken(r)
	if !ok {
		return models.ResponseObject{Error: token, Code: http.StatusUnauthorized}, nil

	}

	return token_grpc.TokenRenewRequest{RefreshToken: token}, nil

}

func logoutHandler(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodGet {
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
		return nil, http.ErrNotSupported
	}
	user := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return models.ResponseObject{Error: "Bad request fromat", Code: http.StatusBadRequest}, nil
	}

	//make a http post to a second service
	resp, err := http.PostForm(servieURL+"login",
		url.Values{"username": {user.Username}, "password": {user.Password}})
	if err != nil {
		return models.ResponseObject{Error: "Unreachable service", Code: http.StatusNoContent}, nil
	}
	defer resp.Body.Close()

	user = models.User{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		d, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return models.ResponseObject{Error: err.Error(), Code: http.StatusNoContent}, nil
		}
		return models.ResponseObject{Error: string(d), Code: http.StatusBadRequest}, nil
	}

	if user.Error != "" {
		return models.ResponseObject{Error: user.Error, Code: user.Code}, nil
	}

	if user.ID < 1 {
		return models.ResponseObject{Error: "Sorry, there's been and error processing this request", Code: http.StatusNoContent}, nil
	}

	claims := make(map[string]string)
	for k, v := range user.Claims {
		claims[k] = v.(string)
	}
	claims["id"] = fmt.Sprintf("%v", user.ID)
	return token_grpc.TokenRequest{Claims: claims}, nil
	// return models.ResponseObject{Error: "Sorry, the login credentials don't match any records", Code: http.StatusNoContent}, nil

}

func registerHandler(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodPost {
		return nil, http.ErrNotSupported
	}
	user := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return models.ResponseObject{Error: "Bad request fromat", Code: http.StatusBadRequest}, nil

	}

	if user.PII.Mobile == "" {
		return models.ResponseObject{Error: "Sorry, please include a valid mobile number", Code: http.StatusUnprocessableEntity}, nil
	}

	//make a http post to a second service
	resp, err := http.PostForm(servieURL+"signup",
		url.Values{"username": {user.Username}, "password": {user.Password}})
	if err != nil {
		return models.ResponseObject{Error: "Unreachable service", Code: http.StatusNoContent}, nil
	}
	defer resp.Body.Close()

	user = models.User{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		d, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return models.ResponseObject{Error: err.Error(), Code: http.StatusNoContent}, nil
		}
		return models.ResponseObject{Error: string(d), Code: http.StatusBadRequest}, nil
	}

	if user.Error != "" {
		return models.ResponseObject{Error: user.Error, Code: user.Code}, nil
	}

	if user.ID < 1 {
		return models.ResponseObject{Error: "Sorry, there's been and error processing this request", Code: http.StatusNoContent}, nil
	}

	claims := make(map[string]string)
	for k, v := range user.Claims {
		claims[k] = v.(string)
	}
	claims["id"] = fmt.Sprintf("%v", user.ID)
	return token_grpc.TokenRequest{Claims: claims}, nil

}

func extractAuthToken(r *http.Request) (string, bool) {
	token := r.Header.Get("AUTH_API_KEY")

	if token == "" {
		return "Unauthorized access, no api key found with your request", false
	}

	return token, true
}
