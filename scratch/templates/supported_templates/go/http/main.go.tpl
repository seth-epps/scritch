package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleGet))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("error starting server: %v", err)
	}
}

func handleGet(w http.ResponseWriter, req *http.Request) {
	log.Printf("Recieved req [%v]", req)
	res := Response{"Hello!"}
	json.NewEncoder(w).Encode(res)
}
