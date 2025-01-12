package plugin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type jwtAuthenticator struct {
	secret     string
	NextPlugin PluginInterface
}

func NewJwtAuthenticator(secret string) PluginInterface {
	return &jwtAuthenticator{
		secret: secret,
	}
}

func (j *jwtAuthenticator) Handle(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("Authorization")

	if !j.validateToken(tokenString) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

}

func (j jwtAuthenticator) AddNext(nextPlugin PluginInterface) {

	j.NextPlugin = nextPlugin

}

func (j jwtAuthenticator) validateToken(tokenString string) bool {

	if tokenString == "" {
		return false
	}

	tokenParts := strings.Split(tokenString, "Bearer ")
	if len(tokenParts) != 2 {
		return false
	}

	jwtToken := tokenParts[1]

	// validated Token
	parsedToken, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		fmt.Println("token:", token.Header, token.Method, token.Signature, token.Signature)
		return []byte(j.secret), nil
	})

	if err != nil {
		fmt.Println("Invalid token :", err)
		return false
	}

	// parsedToken.Valid
	if !parsedToken.Valid {
		return false
	}

	if parsedToken.Method.Alg() != "HS256" {
		return false
	}

	return true

}
