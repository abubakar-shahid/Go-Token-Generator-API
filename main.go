package main

import (
	"Go-Token-Generator-API/api/handler" // Ensure the correct import path
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/get-token", handler.GetTokenHandler)
	fmt.Println("Server running at port 8080...")
	http.ListenAndServe(":8080", nil)
}
