package v1.service

import (
	"context"
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
	Error   string        `json:"error,omitempty"`
	Code    int           `json:"code,omitempty"`
}

func(at *AccessTokens)isResponse(){}

func(at *ServiceError)isResponse(){}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUUID   string `json:"access_uuid"`
	RefreshUUID  string `json:"refresh_uuid"`
	AtExpiry     int64  `json:"at_expiry"`
	RtExpiry     int64  `json:"rt_expiry"`
}

type TokenClaim struct{
	Claim string `json:"claim"`
	Value interface{} `json:"value"`
}

type TokenServiceInterface interface {
	Generate(ctx context.Context, claims map[string]interface{}) (*AccessTokens, error)
	// VerifyToken(ctx context.Context, in *TokenVerifyRequest, opts ...grpc.CallOption) (*TokenVerifyResponse, error)
	// RenewTokens(ctx context.Context, in *TokenRenewRequest, opts ...grpc.CallOption) (*TokenResponse, error)
	// AffectToken(ctx context.Context, in *TokenAffectRequest, opts ...grpc.CallOption) (*TokenAffectResponse, error)
}

type TokenService struct {
}

func (ts *TokenService) Generate(ctx context.Context, claims map[string]interface{}) (*AccessTokens, error) {

	td := TokenDetails{}
	td.AtExpiry = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpiry = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	// create access token
	atClaims := jwt.MapClaims{}
	if ok:=mergeClaims(atClaims,claims);!ok{
		return nil, error.New("Error adding claims to new token")
	}
	atClaims["authorize"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["exp"] = td.AtExpiry
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	//create refresh token
	rtClaims := jwt.MapClaims{}
	
	if ok:=mergeClaims(rtClaims,claims);!ok{
		return nil, error.New("Error adding claims to new token")
	}
	rtClaims["user_id"] = id
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["exp"] = td.RtExpiry
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtoken, rerr := rt.SignedString([]byte(os.Getenv("JWT_R_SECRET")))
	if rerr != nil {
		return nil, rerr
	}

	td.AccessToken = token
	td.RefreshToken = rtoken

	tokens := storeJWTMeta(&td)

	return tokens, nil

}

func mergeClaims(items *jwt.MapClaims,claims )bool{
	for claim,value:= range(claims){
		items[claim]=value
	}
	return true
}

func storeJWTMeta(td *TokenDetails) *AccessTokens {
	fmt.Printf("Storing tokens : %v",td)
	// storage ...
	return &AccessTokens{AccessToken: td.AccessToken, RefreshToken: td.RefreshToken}
}
