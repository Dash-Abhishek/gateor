package plugin

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthenticator(t *testing.T) {

	validator := jwtAuthenticator{
		secret: "mysecret",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  time.Now().Unix(),
	})

	secretKey := "mysecret"
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	t.Logf("Generated Token: %s", tokenString)
	if !validator.validateToken(fmt.Sprintf("Bearer %s", tokenString)) {
		t.Error("should have validated successfully")
	}

}
