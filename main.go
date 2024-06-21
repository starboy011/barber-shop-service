package main

import (
	"fmt"
	"net/http"

	"github.com/starboy011/barber-shop-service/internal/handler"
)

func main() {

	fmt.Println("Starting server on :8081")
	http.HandleFunc("/shops", handler.GetShops)
	http.HandleFunc("/upload-images", handler.UploadImagesHandler)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
