package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

type AccessTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

type TokenServiceInterface interface {
	Generate(ctx context.Context, claims map[string]string) (*AccessTokens, error)
	// VerifyToken(ctx context.Context, in *TokenVerifyRequest, opts ...grpc.CallOption) (*TokenVerifyResponse, error)
	// RenewTokens(ctx context.Context, in *TokenRenewRequest, opts ...grpc.CallOption) (*TokenResponse, error)
	// AffectToken(ctx context.Context, in *TokenAffectRequest, opts ...grpc.CallOption) (*TokenAffectResponse, error)
}

type TokenService struct {
}

func (ts TokenService) Generate(ctx context.Context, claims map[string]string) (*AccessTokens, error) {

	td := TokenDetails{}
	td.AtExpiry = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpiry = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	// create access token
	atClaims := mergeClaims(claims)
	atClaims["authorize"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["exp"] = td.AtExpiry
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	//create refresh token
	rtClaims := mergeClaims(claims)

	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["exp"] = td.RtExpiry
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtoken, rerr := rt.SignedString([]byte(os.Getenv("JWT_R_SECRET")))
	if rerr != nil {
		return nil, rerr
	}

	td.AccessToken = token
	td.RefreshToken = rtoken

	tokens := storeJWTMeta(td)
	fmt.Printf("token details \nAT:{%v}\nRT:{%v}\n", token, rtoken)

	return tokens, nil

}

func mergeClaims(claims map[string]string) jwt.MapClaims {
	c := jwt.MapClaims{}
	for claim, value := range claims {
		c[claim] = value
	}
	return c
}

func storeJWTMeta(td TokenDetails) *AccessTokens {
	fmt.Printf("Storing tokens : %v", td)
	// storage ...
	return &AccessTokens{AccessToken: td.AccessToken, RefreshToken: td.RefreshToken}
}
