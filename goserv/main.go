package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"github.com/BolvicBolvicovic/scraper/request"
)



func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", request.GetRoot)

	fmt.Println("Server starts on: localhost:3333")
	err := http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
