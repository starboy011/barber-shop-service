package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadImagesHandler(w http.ResponseWriter, r *http.Request) {
	err := saveFileInServer(w, r)
	if err != nil {
		// Log the error or handle it appropriately
		fmt.Printf("error in saveing file in server: %v\n", err)
		DeleteUploadedFiles(r)
		w.Write([]byte("error in saveing file in server"))
		return
	}
	err = UploadFileInS3(w, r)
	if err != nil {
		// Log the error or handle it appropriately
		fmt.Printf("error in uploading file to S3: %v\n", err)
		DeleteUploadedFiles(r)
		w.Write([]byte("error in uploading file to S3:"))
		return
	}
	// Respond to the client

	DeleteUploadedFiles(r)
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
	r.ParseMultipartForm(10 << 20)
	files := r.MultipartForm.File["images"]
	shopId := r.FormValue("shopId")
	shopDir := filepath.Join("./uploads", shopId, "home-image")
	for _, fileHeader := range files {
		dstPath := filepath.Join(shopDir, fileHeader.Filename)
		os.Remove(dstPath)
	}
	os.RemoveAll(shopDir)
	parentDir := filepath.Join("./uploads", shopId)
	os.Remove(parentDir)

}

func UploadFileInS3(w http.ResponseWriter, r *http.Request) error {
	shopId := r.FormValue("shopId")
	if shopId == "" {
		return fmt.Errorf("shopId parameter is required")
	}

	// Get bucket name from environment variable
	bucket := "barber-shop-home-images"
	if bucket == "" {
		return fmt.Errorf("BARBER_SHOP_HOME_IMAGES_BUCKET_NAME environment variable not set")
	}
	// Extract the filename with extension
	files := r.MultipartForm.File["images"]
	if len(files) > 5 {
		http.Error(w, "too many files uploaded. Maximum is 5", http.StatusBadRequest)
		return fmt.Errorf("too many files uploaded. Maximum is 5")
	}

	for _, fileHeader := range files {
		// Validate file extension
		filename := fileHeader.Filename

		filePath := fmt.Sprintf("./uploads/%s/home-image/%s", shopId, filename)

		// Initialize AWS session in the default region configuration
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("eu-north-1"), // Specify your AWS Region
		})
		if err != nil {
			fmt.Println("Failed to set session", err)
			return fmt.Errorf("failed to set session")
		}
		// Create S3 service client
		svc := s3.New(sess)

		// Open the file
		fileToUpload, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Failed to open file", err)
			return fmt.Errorf("failed to open file")
		}
		defer fileToUpload.Close()

		// Upload the file to S3
		_, err = svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fmt.Sprintf("%s/home-images/%s", shopId, filename)), // Object key
			Body:   fileToUpload,
		})
		if err != nil {
			fmt.Println("Failed to upload data to S3", err)
			return fmt.Errorf("failed to upload data to S3")
		}
	}

	// Example file path, replace this with your dynamic file path logic

	fmt.Println("File uploaded successfully")
	return nil
}
