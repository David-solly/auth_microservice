package models

import "github.com/dgrijalva/jwt-go"

type AccessTokens struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token"`
}

type AccessDetails struct {
	AccessUuid  string `json:"access_uuid,omitempty"`
	RefreshUUID string `json:"refresh_uuid,omitempty"`
	UserId      uint64 `json:"user_id"`
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
	Access ServiceAccess `json:"access,omitempty"`
	Error  ServiceError  `json:"error,omitempty"`
}

type ServiceAccess struct {
	UserID uint64         `json:"id,omitempty"`
	Status TokenStatus    `json:"status,omitempty"`
	Claims *jwt.MapClaims `json:"claims,omitempty"`
}

type TokenAffectRequest struct {
	Token        string     `json:"token"`
	DesiredState TokenState `json:"desired_state,omitempty"`
}

type TokenAffectResponse struct {
	EffectApplied bool          `json:"effect_applied,omitempty"`
	Error         *ServiceError `json:"error,omitempty"`
}

type TokenState int32

const (
	TokenState_WATCH      TokenState = 0
	TokenState_LOGOUT     TokenState = 1
	TokenState_TRACEROUTE TokenState = 2
	TokenState_FREEZE     TokenState = 3
	TokenState_UNFREEZE   TokenState = 4
)

var TokenState_name = map[int32]string{
	0: "WATCH",
	1: "LOGOUT",
	2: "TRACEROUTE",
	3: "FREEZE",
	4: "UNFREEZE",
}

var TokenState_value = map[string]int32{
	"WATCH":      0,
	"LOGOUT":     1,
	"TRACEROUTE": 2,
	"FREEZE":     3,
	"UNFREEZE":   4,
}

type TokenStatus int32

const (
	TokenStatus_INVALID    TokenStatus = 0
	TokenStatus_AUTHORIZED TokenStatus = 1
	TokenStatus_RESTRICTED TokenStatus = 2
	TokenStatus_EXPIRED    TokenStatus = 3
)
