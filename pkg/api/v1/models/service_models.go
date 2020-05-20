package models

import "github.com/dgrijalva/jwt-go"

type AccessTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

type ServiceError struct {
	Error string `json:"error,omitempty"`
	Code  int    `json:"code,omitempty"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUUID   string `json:"access_uuid"`
	RefreshUUID  string `json:"refresh_uuid"`
	AtExpiry     int64  `json:"at_expiry"`
	RtExpiry     int64  `json:"rt_expiry"`
}

type TokenClaim struct {
	Claim string      `json:"claim"`
	Value interface{} `json:"value"`
}

type TokenVerifyResponse struct {
	UserID uint64         `json:"id,omitempty"`
	Status TokenStatus    `json:"status,omitempty"`
	Claims *jwt.MapClaims `json:"claims,omitempty"`
	Error  ServiceError   `json:"error,omitempty"`
}

type TokenStatus int32

const (
	TokenStatus_INVALID    TokenStatus = 0
	TokenStatus_AUTHORIZED TokenStatus = 1
	TokenStatus_RESTRICTED TokenStatus = 2
	TokenStatus_EXPIRED    TokenStatus = 3
)
