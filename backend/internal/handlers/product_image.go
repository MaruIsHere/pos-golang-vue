package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxProductImageSize = 300 * 1024
const maxProductImageDimension = 1200

func UploadProductImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxProductImageSize+64*1024)
	if err := c.Request.ParseMultipartForm(maxProductImageSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gambar terlalu besar atau format unggahan tidak valid"})
		return
	}
	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}

	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pilih gambar produk terlebih dahulu"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxProductImageSize+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca gambar"})
		return
	}
	if len(data) > maxProductImageSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Ukuran gambar hasil kompresi maksimal 300 KB"})
		return
	}
	if !strings.HasPrefix(http.DetectContentType(data), "image/jpeg") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gambar harus berupa JPEG hasil kompresi"})
		return
	}
	config, err := jpeg.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width < 1 || config.Height < 1 ||
		config.Width > maxProductImageDimension || config.Height > maxProductImageDimension {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dimensi gambar maksimal 1200 x 1200 piksel"})
		return
	}
	if _, _, err := image.Decode(bytes.NewReader(data)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File gambar JPEG tidak valid"})
		return
	}

	var fileID [16]byte
	if _, err := rand.Read(fileID[:]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat nama file gambar"})
		return
	}
	const uploadDir = "uploads/products"
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyiapkan penyimpanan gambar"})
		return
	}
	filename := hex.EncodeToString(fileID[:]) + ".jpg"
	filePath := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(filePath, data, 0640); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"image_url": fmt.Sprintf("/uploads/products/%s", filename)})
}
