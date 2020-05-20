package service

import (
	"os"
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

const accesstokenToTest = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjFhYjRmMGFmLTEzNGItNGMxMC05NjRhLWQwNjU5YzFiMmNkMiIsImNsZWFyYW5jZSI6ImRlbHRhIiwiZXhwIjoxNTkwMDE4NDYxLCJpZCI6IjIiLCJzZXJ2aWNlIjoiY29tLmJpZy5iZW4ifQ.z687q82iP42VWRlWyR3gFrJhsMN_6YDM3P7v5rFL_G4"
const refreshTokenToTest = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTA2MjIzNjEsImlkIjoiMiIsInJlZnJlc2hfdXVpZCI6IjM2OThiMTg2LTBkNTctNDcxZS05N2ZjLTQ3Y2ViOGZkOWRhYSJ9.IqPJzW51-ruOwSHkQmjpcoYHApddUrlgO6GfiZhrV-I"

func TestJWTExtraction(t *testing.T) {
	os.Setenv("JWT_SECRET", "superduperdupersecuresecret23232string6568_55")
	t.Run("VERIFY and EXTRACT jwt token", func(t *testing.T) {
		got, err := VerifyTokenIntegrity(accesstokenToTest)
		assert.NilError(t, err)
		assert.NotNil(t, got)

	})

}
