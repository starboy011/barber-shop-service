package main

import (
	"fmt"
	"net/http"
)

func main() {

	fmt.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
