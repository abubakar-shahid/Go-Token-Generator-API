// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/rs/cors"
// )

// type Response struct {
// 	Token string `json:"token"`
// }

// type UserInfo struct {
// 	Identity string `json:"identity"`
// 	Name     string `json:"name"`
// 	Metadata string `json:"metadata"`
// }

// type RequestBody struct {
// 	RoomName string   `json:"roomName"`
// 	UserInfo UserInfo `json:"userInfo"`
// }

// func GenerateToken(userInfo UserInfo) (string, error) {
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	claims := token.Claims.(jwt.MapClaims)

// 	claims["identity"] = userInfo.Identity
// 	claims["name"] = userInfo.Name
// 	claims["metadata"] = userInfo.Metadata
// 	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

// 	tokenString, err := token.SignedString([]byte("secret"))

// 	if err != nil {
// 		fmt.Errorf("Something Went Wrong: %s", err.Error())
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Client Requested for Token:")
// 	var requestBody RequestBody
// 	err := json.NewDecoder(r.Body).Decode(&requestBody)
// 	log.Println("data received: ", requestBody)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	tokenString, err := GenerateToken(requestBody.UserInfo)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	response := Response{Token: tokenString}
// 	json.NewEncoder(w).Encode(response)
// 	log.Println("Token Sent!")
// }

// func main() {
// 	c := cors.New(cors.Options{
// 		AllowedOrigins: []string{"http://localhost:3000"},
// 		AllowedMethods: []string{http.MethodPost},
// 		AllowedHeaders: []string{"Content-Type"},
// 	})

// 	http.Handle("/get-token", c.Handler(http.HandlerFunc(GetTokenHandler)))
// 	log.Println("Server running on port 8080...")
// 	http.ListenAndServe(":8080", nil)
// }

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/cors"
)

type Response struct {
	Token string `json:"token"`
}

type UserInfo struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	Metadata string `json:"metadata"`
}

type RequestBody struct {
	RoomName string   `json:"roomName"`
	UserInfo UserInfo `json:"userInfo"`
}

var (
	apiKey    = os.Getenv("devkey")
	apiSecret = os.Getenv("secret")
)

func GenerateToken(userInfo UserInfo, apiKey, apiSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"identity": userInfo.Identity,
		"name":     userInfo.Name,
		"metadata": userInfo.Metadata,
		"exp":      time.Now().Add(time.Minute * 30).Unix(),
		"apiKey":   apiKey,
	})

	tokenString, err := token.SignedString([]byte(apiSecret))
	if err != nil {
		return "", fmt.Errorf("Failed to generate token: %v", err)
	}

	return tokenString, nil
}

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Client requested for token:")
	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	log.Println("data received: ", requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := GenerateToken(requestBody.UserInfo, apiKey, apiSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Println("Token sent!")
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{http.MethodPost},
		AllowedHeaders: []string{"Content-Type"},
	})

	certFile := "ssl/cert.pem"
	keyFile := "ssl/key.pem"

	http.Handle("/generate-token", c.Handler(http.HandlerFunc(GetTokenHandler)))

	log.Println("Server running on https://localhost:8443...")
	if err := http.ListenAndServeTLS(":8443", certFile, keyFile, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
