package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/dghubble/oauth1"
    "github.com/golang-jwt/jwt"
)

const (
    requestTokenURL = "https://<OAuth-provider-URL>/request_token"
    accessTokenURL  = "https://<OAuth-provider-URL>/access_token"
    callbackURL     = "<callback-URL>"
    consumerKey     = "<consumer-key>"
    consumerSecret  = "<consumer-secret>"
    jwtSecret       = "<jwt-secret>"
)

type AccessTokenResponse struct {
    AccessToken       string `json:"access_token"`
    AccessTokenSecret string `json:"access_token_secret"`
}

func main() {
    // Create OAuth1.0a configuration
    config := oauth1.Config{
        ConsumerKey:    consumerKey,
        ConsumerSecret: consumerSecret,
        CallbackURL:    callbackURL,
        Endpoint: oauth1.Endpoint{
            RequestTokenURL: requestTokenURL,
            AccessTokenURL:  accessTokenURL,
        },
    }

    // Create HTTP client for OAuth requests
    httpClient := config.Client(oauth1.NoContext, &oauth1.Token{
        Token:  "<request-token>",
        Secret: "<request-token-secret>",
    })

    // Handler for the main endpoint
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Verify JWT from Authorization header
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(jwtSecret), nil
        })

        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            json.NewEncoder(w).Encode(claims)
        } else {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
        }
    })

    // Handler for OAuth callback
    http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
        // Make request to protected resource
        resp, err := httpClient.Get("https://<API-URL>/protected_resource")
        if err != nil {
            http.Error(w, "Failed to access protected resource", http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()

        // Create JWT token with access token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "access_token": "<access-token>", // Replace with actual access token
        })

        tokenString, err := token.SignedString([]byte(jwtSecret))
        if err != nil {
            http.Error(w, "Failed to create JWT", http.StatusInternalServerError)
            return
        }

        // Return JWT token
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "token": tokenString,
        })
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
