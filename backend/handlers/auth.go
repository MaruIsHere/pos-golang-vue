package handlers

import (
	"fmt"
	"net/http"
	"pos-backend/database"
	"pos-backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("pos_secret_key_2026")

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

type UpdateProfileInput struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	ProfilePhoto string `json:"profile_photo"`
	Password     string `json:"password"`
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, fmt.Sprintf("%v", userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
		return
	}

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data profil tidak valid"})
		return
	}

	input.Username = strings.TrimSpace(input.Username)
	input.Name = strings.TrimSpace(input.Name)
	if len(input.Username) < 3 || len(input.Username) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username harus berisi 3 sampai 100 karakter"})
		return
	}
	if len(input.Name) == 0 || len(input.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama harus diisi dan maksimal 100 karakter"})
		return
	}
	if len(input.ProfilePhoto) > 3*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foto profil maksimal 2MB"})
		return
	}
	if input.Password != "" && len(input.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 6 karakter"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, fmt.Sprintf("%v", userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	var duplicateCount int64
	database.DB.Model(&models.User{}).Where("LOWER(username) = LOWER(?) AND id <> ?", input.Username, user.ID).Count(&duplicateCount)
	if duplicateCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username sudah digunakan"})
		return
	}

	updates := map[string]interface{}{
		"username":      input.Username,
		"name":          input.Name,
		"profile_photo": input.ProfilePhoto,
	}
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses enkripsi password"})
			return
		}
		updates["password"] = string(hashedPassword)
	}
	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profil"})
		return
	}
	if err := database.DB.First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat profil terbaru"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui sesi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil berhasil diperbarui",
		"token":   tokenString,
		"user":    user,
	})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Select("id", "username", "name", "role", "created_at").Order("created_at desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUserProfile(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUserProfile(c *gin.Context) {
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi pengguna tidak valid"})
		return
	}
	if strings.TrimSpace(c.Param("id")) == strings.TrimSpace(fmt.Sprintf("%v", currentUserID)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gunakan menu Profil Akun untuk mengubah akun sendiri"})
		return
	}

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data profil tidak valid"})
		return
	}
	input.Username = strings.TrimSpace(input.Username)
	input.Name = strings.TrimSpace(input.Name)
	if len(input.Username) < 3 || len(input.Username) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username harus berisi 3 sampai 100 karakter"})
		return
	}
	if len(input.Name) == 0 || len(input.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama harus diisi dan maksimal 100 karakter"})
		return
	}
	if len(input.ProfilePhoto) > 3*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foto profil maksimal 2MB"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}
	var duplicateCount int64
	database.DB.Model(&models.User{}).Where("LOWER(username) = LOWER(?) AND id <> ?", input.Username, user.ID).Count(&duplicateCount)
	if duplicateCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Username sudah digunakan"})
		return
	}

	updates := map[string]interface{}{
		"username":      input.Username,
		"name":          input.Name,
		"profile_photo": input.ProfilePhoto,
	}
	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profil pengguna"})
		return
	}
	if err := database.DB.First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat profil terbaru"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profil pengguna berhasil diperbarui", "user": user})
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
		return
	}

	username := strings.TrimSpace(input.Username)
	if len(username) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username minimal 3 karakter"})
		return
	}

	if len(input.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 6 karakter"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses enkripsi password"})
		return
	}

	role := "kasir"
	reqRole := strings.ToLower(strings.TrimSpace(input.Role))
	if reqRole == "owner" || reqRole == "admin" || reqRole == "kepala_kasir" || reqRole == "kasir" {
		role = reqRole
	}

	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Role:     models.UserRole(role),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna baru berhasil ditambahkan"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	currentUserId, exists := c.Get("user_id")
	if exists && strings.TrimSpace(id) == strings.TrimSpace(fmt.Sprintf("%v", currentUserId)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak dapat menghapus akun Anda sendiri yang sedang digunakan"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}

type ChangeUserPasswordInput struct {
	Password string `json:"password" binding:"required"`
}

func ChangeUserPassword(c *gin.Context) {
	var input ChangeUserPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password wajib diisi"})
		return
	}

	if len(input.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 6 karakter"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses enkripsi password"})
		return
	}

	if err := database.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui password pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password pengguna berhasil diperbarui"})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
		return
	}

	username := strings.TrimSpace(input.Username)

	var user models.User
	if err := database.DB.Where("LOWER(username) = LOWER(?)", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 Days Token
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token akses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":            user.ID,
			"username":      user.Username,
			"name":          user.Name,
			"profile_photo": user.ProfilePhoto,
			"role":          user.Role,
		},
	})
}
