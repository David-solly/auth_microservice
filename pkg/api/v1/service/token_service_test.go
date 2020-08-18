package service

import (
	"os"
	"reflect"
	"testing"

	"github.com/David-solly/auth_microservice/pkg/api/v1/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/docker/docker/pkg/testutil/assert"
)

var (
	expiredToken       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjFhYjRmMGFmLTEzNGItNGMxMC05NjRhLWQwNjU5YzFiMmNkMiIsImNsZWFyYW5jZSI6ImRlbHRhIiwiZXhwIjoxNTkwMDE4NDYxLCJpZCI6IjIiLCJzZXJ2aWNlIjoiY29tLmJpZy5iZW4ifQ.z687q82iP42VWRlWyR3gFrJhsMN_6YDM3P7v5rFL_G4"
	accesstokenToTest  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjFhYjRmMGFmLTEzNGItNGMxMC05NjRhLWQwNjU5YzFiMmNkMiIsImNsZWFyYW5jZSI6ImRlbHRhIiwiZXhwIjoxNTkwMDE4NDYxLCJpZCI6IjIiLCJzZXJ2aWNlIjoiY29tLmJpZy5iZW4ifQ.z687q82iP42VWRlWyR3gFrJhsMN_6YDM3P7v5rFL_G4"
	refreshTokenToTest = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTA2MjIzNjEsImlkIjoiMiIsInJlZnJlc2hfdXVpZCI6IjM2OThiMTg2LTBkNTctNDcxZS05N2ZjLTQ3Y2ViOGZkOWRhYSJ9.IqPJzW51-ruOwSHkQmjpcoYHApddUrlgO6GfiZhrV-I"
	claimsToRenew      = map[string]string{"service": "testing.service", "id": "123"}
	UUID, invalidUUID  string
)

func setEnvirons() {
	os.Setenv("JWT_SECRET", "superduperdupersecuresecret23232string6568_55")
	os.Setenv("JWT_R_SECRET", "superduperdupersecsjld_kdkjuresecret23232string2d545565656_")
	RedisInit()
	initTokensVars()

}
func initTokensVars() {

	tokens, _ := generateTokenPair(claimsToRenew)
	accesstokenToTest = tokens.AccessToken
	refreshTokenToTest = tokens.RefreshToken

	meta, _, _ := ExtractTokenMetadata(refreshTokenToTest, true)
	UUID = meta.RefreshUUID
	invalidUUID = "0c881322-f4fe-49c0-8619-7ada6bb642a4"

}

func TestRedisConnection(t *testing.T) {
	t.Run("INIT redis db connection", func(t *testing.T) {
		k, err := RedisInit()
		assert.NilError(t, err)
		assert.Equal(t, k, "PONG")
	})

}

func TestJWTGeneration(t *testing.T) {
	setEnvirons()
	t.Run("GENERATE jwt token pair", func(t *testing.T) {
		tokens, err := generateTokenPair(claimsToRenew)
		assert.NilError(t, err)
		assert.Equal(t, reflect.ValueOf(tokens).IsNil(), false)
		assert.Equal(t, tokens.AccessToken != "", true)
		assert.Equal(t, tokens.RefreshToken != "", true)

	})
}

func TestJWTExtraction(t *testing.T) {
	setEnvirons()
	t.Run("VERIFY and EXTRACT jwt token", func(t *testing.T) {
		got, err := VerifyTokenIntegrity(accesstokenToTest, false)
		assert.NilError(t, err)
		assert.NotNil(t, got)

	})
	t.Run("VERIFY and EXTRACT REFRESH jwt token", func(t *testing.T) {
		got, err := VerifyTokenIntegrity(refreshTokenToTest, true)
		assert.NilError(t, err)
		assert.NotNil(t, got)

	})

	t.Run("FAIL to VERIFY and EXTRACT REFRESH jwt token", func(t *testing.T) {
		got, err := VerifyTokenIntegrity(refreshTokenToTest, false)
		assert.Error(t, err, "Invalid")
		assert.NotNil(t, got)

	})

	t.Run("VERIFY and EXTRACT EMPTY jwt token", func(t *testing.T) {
		got, err := VerifyTokenIntegrity("", false)
		assert.Error(t, err, "Invalid")
		assert.Equal(t, true, reflect.ValueOf(got).IsNil())

	})

}

func TestReadingRefreshData(t *testing.T) {
	setEnvirons()
	k := models.AccessDetails{RefreshUUID: UUID}
	t.Run("VERIFY and REDIS READ claims from UUID", func(t *testing.T) {
		got, err := FetchRefresh(&k)
		assert.NilError(t, err)
		assert.Equal(t, reflect.TypeOf(got), reflect.TypeOf(map[string]string{}))
		assert.Equal(t, got["service"], "testing.service")

	})
	t.Run("VERIFY and REDIS READ claims from UUID Fail", func(t *testing.T) {
		got, err := FetchRefresh(&models.AccessDetails{RefreshUUID: invalidUUID})
		assert.Error(t, err, "nil")
		assert.Equal(t, reflect.TypeOf(got), reflect.TypeOf(map[string]string{}))

	})

}

func TestRefreshToken(t *testing.T) {
	setEnvirons()
	validRefreshToken := refreshTokenToTest
	ids := &models.AccessDetails{}
	claimsToRenew := map[string]string{}
	t.Run("VERIFY Token then Extract UUID", func(t *testing.T) {
		got, claims, err := ExtractTokenMetadata(validRefreshToken, true)

		assert.NilError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, reflect.TypeOf(claims), reflect.TypeOf(jwt.MapClaims{}))
		assert.Equal(t, true, len(claims) > 0)
		assert.Equal(t, got.AccessUuid != "", true)
		assert.Equal(t, got.RefreshUUID != "", true)
		// assert.Equal(t, got.AccessUuid, "27fba8b4-1368-4c91-89ae-3495ae6a8df6") //To test specific uuid values
		// assert.Equal(t, got.RefreshUUID, "471db4f4-313a-4ddc-a79e-2b4bb8e5412c") //test specific uuid values
		ids = got
		assert.NotNil(t, ids)
	})

	t.Run("VERIFY Token FAIL and return INVALID", func(t *testing.T) {
		got, claims, err := ExtractTokenMetadata(validRefreshToken, false)
		assert.Error(t, err, "Invalid")
		assert.Equal(t, true, reflect.ValueOf(got).IsNil())
		assert.Equal(t, reflect.TypeOf(claims), reflect.TypeOf(jwt.MapClaims{}))

	})

	t.Run("GET REFRESH CLAIMS and REDIS READ claims from UUID", func(t *testing.T) {
		got, err := FetchRefresh(ids)
		assert.NilError(t, err)
		assert.Equal(t, reflect.TypeOf(got), reflect.TypeOf(map[string]string{}))
		assert.Equal(t, true, len(got) > 0)
		claimsToRenew = got
		assert.NotNil(t, claimsToRenew)
	})

	t.Run("DELETE old AUTH TOKEN", func(t *testing.T) {
		got, err := deleteAuth(ids.AccessUuid)
		assert.NilError(t, err)
		assert.Equal(t, true, got > 0)
	})

	t.Run("DELETE old REFRESH TOKEN", func(t *testing.T) {
		got, err := deleteAuth(ids.RefreshUUID)
		assert.NilError(t, err)
		assert.Equal(t, true, got > 0)
	})

	t.Run("GENERATE new TOKEN PAIR", func(t *testing.T) {
		tokens, err := generateTokenPair(claimsToRenew)
		assert.NilError(t, err)
		assert.Equal(t, reflect.ValueOf(tokens).IsNil(), false)
		assert.Equal(t, tokens.AccessToken != "", true)
		assert.Equal(t, tokens.RefreshToken != "", true)
	})

	t.Run("VERIFY DELETED and REDIS READ claims from UUID", func(t *testing.T) {
		got, err := FetchRefresh(ids)
		assert.Error(t, err, "nil")
		assert.Equal(t, reflect.ValueOf(got).IsNil(), true)
		assert.Equal(t, false, len(got) > 0)
		assert.NotNil(t, claimsToRenew)
	})

}

func BenchmarkExtract(t *testing.B) {
	setEnvirons()
	t.Run("EXTRACT valid TOKEN from JWT", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			_, _ = VerifyTokenIntegrity(accesstokenToTest, false)
		}

	})
}
