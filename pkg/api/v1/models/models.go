package service.models

import (
	token_grpc "github.com/David-solly/auth_microservice/pkg/api/v1/service"
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
	ID       uint64                 `json:"id"`
	Username string                 `json:"email"`
	Password string                 `json:"password"`
	PII      UserPII                `json:"details,omitempty"`
	Claims   map[string]interface{} `json:"claims,omitempty"`
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
