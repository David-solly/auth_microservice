package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/David-solly/auth_microservice/pkg/api/v1/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/twinj/uuid"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

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

type TokenServiceInterface interface {
	Generate(ctx context.Context, claims map[string]string) (*models.AccessTokens, error)
	VerifyToken(ctx context.Context, tokenToverify TokenVerifyRequest) (*models.TokenVerifyResponse, error)
	// RenewTokens(ctx context.Context, in *TokenRenewRequest, opts ...grpc.CallOption) (*TokenResponse, error)
	// AffectToken(ctx context.Context, in *TokenAffectRequest, opts ...grpc.CallOption) (*TokenAffectResponse, error)
}

type TokenService struct {
}

func (ts TokenService) Generate(ctx context.Context, claims map[string]string) (*models.AccessTokens, error) {

	td := models.TokenDetails{}
	td.AtExpiry = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpiry = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	// create access token
	atClaims := mergeClaims(claims)
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

	tokens, err := createAuth(rtClaims["id"].(string), &td)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("token details \nAT:{%v}\nRT:{%v}\n", token, rtoken)

	return tokens, nil

}

func (ts TokenService) VerifyToken(ctx context.Context, tokenToverify TokenVerifyRequest) (*models.TokenVerifyResponse, error) {

	return nil, nil
}

func verifyAndGetTokenClaims(token string) (*models.TokenVerifyResponse, *models.ResponseObject, bool) {
	tokenAuth, tokenClaims, err := ExtractTokenMetadata(token)
	if err != nil {
		return nil, &models.ResponseObject{Error: "Unauthorized token", Code: http.StatusUnauthorized}, false
	}

	userID, err := FetchAuth(tokenAuth)
	if err != nil {
		return nil, &models.ResponseObject{Error: "Unauthorized for resource", Code: http.StatusUnauthorized}, false
	}

	return &models.TokenVerifyResponse{UserID: userID, Claims: &tokenClaims}, nil, true
}

func mergeClaims(claims map[string]string) jwt.MapClaims {
	c := jwt.MapClaims{}
	for claim, value := range claims {
		c[claim] = value
	}
	return c
}

func createAuth(userid string, td *models.TokenDetails) (*models.AccessTokens, error) {
	at := time.Unix(td.AtExpiry, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpiry, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUUID, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return nil, errAccess
	}
	errRefresh := client.Set(td.RefreshUUID, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return nil, errRefresh
	}

	// fmt.Printf("Storing tokens : %v", td)
	// storage ...
	return &models.AccessTokens{AccessToken: td.AccessToken, RefreshToken: td.RefreshToken}, nil

}

func deleteAuth(uuid string) (int64, error) {
	idDeleted, err := client.Del(uuid).Result()
	if err != nil {
		return -1, err
	}

	return idDeleted, nil
}

// FetchAuth : ensure the token hasn't expired
func FetchAuth(authD *models.AccessDetails) (uint64, error) {
	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func ExtractTokenMetadata(tokenString string) (*models.AccessDetails, jwt.MapClaims, error) {
	token, err := VerifyTokenIntegrity(tokenString)
	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, nil, err
		}

		claimID := claims["id"]
		userID, err := strconv.ParseUint(claimID.(string), 10, 64)
		if err != nil {
			return nil, nil, err
		}

		return &models.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userID,
		}, claims, nil
	}
	return nil, nil, err
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
