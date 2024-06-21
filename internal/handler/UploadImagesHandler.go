package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadImagesHandler(w http.ResponseWriter, r *http.Request) {
	err := saveFileInServer(w, r)
	if err != nil {
		// Log the error or handle it appropriately
		fmt.Printf("error in saveing file in server: %v\n", err)
		DeleteUploadedFiles(r)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Files uploaded successfully"))
}

func saveFileInServer(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return fmt.Errorf("invalid request method")
	}

	// Parse the multipart form data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "error parsing form data", http.StatusBadRequest)
		return fmt.Errorf("error parsing form data: %v", err)
	}

	// Retrieve shopId from request
	shopId := r.FormValue("shopId")
	if shopId == "" {
		http.Error(w, "missing shopId parameter", http.StatusBadRequest)
		return fmt.Errorf("missing shopId parameter")
	}

	// Ensure the uploads directory exists
	uploadsPath := "./uploads"
	if _, err := os.Stat(uploadsPath); os.IsNotExist(err) {
		if err := os.Mkdir(uploadsPath, os.ModePerm); err != nil {
			http.Error(w, "error creating uploads directory", http.StatusInternalServerError)
			return fmt.Errorf("error creating uploads directory: %v", err)
		}
	}

	// Create a directory for the shopId if it doesn't exist
	shopDir := filepath.Join(uploadsPath, shopId)
	if _, err := os.Stat(shopDir); os.IsNotExist(err) {
		if err := os.Mkdir(shopDir, os.ModePerm); err != nil {
			http.Error(w, "error creating shop directory", http.StatusInternalServerError)
			return fmt.Errorf("error creating shop directory: %v", err)
		}
	}

	// Create a directory named 'home-image' inside the shop directory
	homeImageDir := filepath.Join(shopDir, "home-image")
	if _, err := os.Stat(homeImageDir); os.IsNotExist(err) {
		if err := os.Mkdir(homeImageDir, os.ModePerm); err != nil {
			http.Error(w, "error creating home-image directory", http.StatusInternalServerError)
			return fmt.Errorf("error creating home-image directory: %v", err)
		}
	}

	// Retrieve files from the form data
	files := r.MultipartForm.File["images"]
	if len(files) > 5 {
		http.Error(w, "too many files uploaded. Maximum is 5", http.StatusBadRequest)
		return fmt.Errorf("too many files uploaded. Maximum is 5")
	}

	// Define allowed file extensions
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	for _, fileHeader := range files {
		// Validate file extension
		ext := filepath.Ext(fileHeader.Filename)
		if !allowedExts[ext] {
			http.Error(w, fmt.Sprintf("unsupported file type: %s", ext), http.StatusBadRequest)
			return fmt.Errorf("unsupported file type: %s", ext)
		}

		// Open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "error opening file", http.StatusInternalServerError)
			return fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()

		// Create a new file in the home-image directory
		dst, err := os.Create(filepath.Join(homeImageDir, fileHeader.Filename))
		if err != nil {
			http.Error(w, "error saving file", http.StatusInternalServerError)
			return fmt.Errorf("error saving file: %v", err)
		}
		defer dst.Close()

		// Copy the file content to the destination
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "error copying file", http.StatusInternalServerError)
			return fmt.Errorf("error copying file: %v", err)
		}
	}

	return nil
}

func DeleteUploadedFiles(r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB
	files := r.MultipartForm.File["images"]
	for _, fileHeader := range files {
		dstPath := filepath.Join("./uploads", r.FormValue("shopId"), "home-image", fileHeader.Filename)
		os.Remove(dstPath)
	}
}
