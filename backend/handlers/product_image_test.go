package handlers

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUploadProductImage(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(workingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/products/images", UploadProductImage)

	jpegData := createTestJPEG(t)
	response := submitProductImage(t, router, jpegData)
	if response.Code != http.StatusCreated {
		t.Fatalf("upload status = %d, want %d: %s", response.Code, http.StatusCreated, response.Body.String())
	}

	var result struct {
		ImageURL string `json:"image_url"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(result.ImageURL, "/uploads/products/") || !strings.HasSuffix(result.ImageURL, ".jpg") {
		t.Fatalf("unexpected image URL: %q", result.ImageURL)
	}
	filename := filepath.Base(result.ImageURL)
	if _, err := os.Stat(filepath.Join("uploads", "products", filename)); err != nil {
		t.Fatalf("uploaded image was not saved: %v", err)
	}

	oversizedImage := append(append([]byte{}, jpegData...), bytes.Repeat([]byte{0}, maxProductImageSize)...)
	oversizedResponse := submitProductImage(t, router, oversizedImage)
	if oversizedResponse.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized upload status = %d, want %d", oversizedResponse.Code, http.StatusRequestEntityTooLarge)
	}
}

func createTestJPEG(t *testing.T) []byte {
	t.Helper()
	picture := image.NewRGBA(image.Rect(0, 0, 2, 2))
	picture.Set(0, 0, color.White)
	var output bytes.Buffer
	if err := jpeg.Encode(&output, picture, nil); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}

func submitProductImage(t *testing.T, router http.Handler, data []byte) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	file, err := writer.CreateFormFile("image", "product.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/products/images", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
