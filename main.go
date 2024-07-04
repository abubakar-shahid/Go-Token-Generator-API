package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/cors"
)

type Response struct {
	Token string `json:"token"`
}

type RequestBody struct {
	RoomName string `json:"roomName"`
	UserInfo struct {
		APIKey    string `json:"apiKey"`
		APISecret string `json:"apiSecret"`
	} `json:"userInfo"`
}

func GenerateToken(apiKey, apiSecret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["apiKey"] = apiKey
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte(apiSecret))

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	log.Println("Token Sent!")
	return tokenString, nil
}

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := GenerateToken(requestBody.UserInfo.APIKey, requestBody.UserInfo.APISecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{Token: tokenString}
	json.NewEncoder(w).Encode(response)
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{http.MethodPost},
		AllowedHeaders: []string{"Content-Type"},
	})

	http.Handle("/get-token", c.Handler(http.HandlerFunc(GetTokenHandler)))
	log.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
