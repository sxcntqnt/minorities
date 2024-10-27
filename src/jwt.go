package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/uber/h3-go/v4"
	"time"
)

func ExampleLatLngToCell() {
	latLng := h3.NewLatLng(37.775938728915946, -122.41795063018799)
	resolution := 9 // between 0 (biggest cell) and 15 (smallest cell)

	cell := h3.LatLngToCell(latLng, resolution)

	fmt.Printf("%s", cell)
	// Output:
	// 8928308280fffff
}

func GenerateToken() (string, error) {
	// Set the secret key
	secretKey := []byte("kadaddy")

	// Define the claims
	claims := jwt.StandardClaims{
		Issuer:    "my-app",
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func main() {
	GenerateToken()
	tokenString, err := GenerateToken()
	if err != nil {
		// Handle error
	}
	fmt.Printf("jwttoken:%s", tokenString)
	// Use the tokenString as the JWT token in the Authorization header
	// e.g. "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

}
