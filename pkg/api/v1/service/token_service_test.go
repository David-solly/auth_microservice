package service

import (
	"os"
	"reflect"
	"testing"

	"github.com/David-solly/auth_microservice/pkg/api/v1/models"
	"github.com/docker/docker/pkg/testutil/assert"
)

const accesstokenToTest = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjFhYjRmMGFmLTEzNGItNGMxMC05NjRhLWQwNjU5YzFiMmNkMiIsImNsZWFyYW5jZSI6ImRlbHRhIiwiZXhwIjoxNTkwMDE4NDYxLCJpZCI6IjIiLCJzZXJ2aWNlIjoiY29tLmJpZy5iZW4ifQ.z687q82iP42VWRlWyR3gFrJhsMN_6YDM3P7v5rFL_G4"
const refreshTokenToTest = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTA2MjIzNjEsImlkIjoiMiIsInJlZnJlc2hfdXVpZCI6IjM2OThiMTg2LTBkNTctNDcxZS05N2ZjLTQ3Y2ViOGZkOWRhYSJ9.IqPJzW51-ruOwSHkQmjpcoYHApddUrlgO6GfiZhrV-I"

func TestJWTExtraction(t *testing.T) {
	os.Setenv("JWT_SECRET", "superduperdupersecuresecret23232string6568_55")
	os.Setenv("JWT_R_SECRET", "superduperdupersecsjld_kdkjuresecret23232string2d545565656_")
	t.Run("VERIFY and EXTRACT jwt token", func(t *testing.T) {
		got, err := VerifyTokenIntegrity(accesstokenToTest, false)
		assert.Error(t, err, "Invalid")
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
	RedisInit()
	k := models.AccessDetails{RefreshUUID: "c635b068-8719-449d-ad66-5298f6046f51"}
	t.Run("VERIFY and REDIS READ claims from UUID", func(t *testing.T) {
		got, err := FetchRefresh(&k)
		assert.NilError(t, err)
		assert.Equal(t, reflect.TypeOf(got), reflect.TypeOf(map[string]string{}))

	})
	t.Run("VERIFY and REDIS READ claims from UUID Fail", func(t *testing.T) {
		got, err := FetchRefresh(&k)
		assert.NilError(t, err)
		assert.Equal(t, reflect.TypeOf(got), reflect.TypeOf(map[string]string{}))

	})

	t.Run("SHOuld fail and REDIS READ claims from UUID", func(t *testing.T) {
		got, err := FetchRefresh(&k)
		assert.NilError(t, err)
		assert.Equal(t, reflect.TypeOf(got), reflect.TypeOf(map[string]string{}))

	})
}

func BenchmarkExtract(t *testing.B) {
	os.Setenv("JWT_SECRET", "superduperdupersecuresecret23232string6568_55")
	t.Run("EXTRACT valid TOKEN from JWT", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			_, _ = VerifyTokenIntegrity(accesstokenToTest, false)
		}

	})
}
