package models

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
	UserID uint64     `json:"id"`
	Status int        `json:"status"`
	Claims jwt.Claims `json:"claims,omitempty"`
}