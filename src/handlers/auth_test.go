package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ia-boilerplate/src/infrastructure"
	"ia-boilerplate/src/logger"
	"ia-boilerplate/src/repository"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupAuthRouter() *gin.Engine {
	logger.Init()
	inf := infrastructure.NewInfrastructure(logger.Log)
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	repo := &repository.Repository{DB: db, Logger: logger.Log, Infrastructure: inf}
	repo.DB.AutoMigrate(&repository.RoleUser{}, &repository.User{}, &repository.DeviceDetails{}, &repository.Medicine{}, &repository.ICDCie{})
	os.Setenv("START_USER_EMAIL", "test@example.com")
	os.Setenv("START_USER_PW", "pass123")
	os.Setenv("ACCESS_SECRET_KEY", "mockAccessSecretKey")
	os.Setenv("REFRESH_SECRET_KEY", "mockRefreshSecretKey")
	os.Setenv("ACCESS_TOKEN_TTL", "15")
	os.Setenv("REFRESH_TOKEN_TTL", "10080")
	os.Setenv("JWT_ISSUER", "mockIssuer")

	repo.SeedInitialRole()
	repo.SeedInitialUser()
	h := NewHandler(repo, zap.NewNop(), inf)
	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/login", h.Login)
	r.POST("/access-token/refresh", h.AccessTokenByRefreshToken)
	return r
}

func TestLoginSuccess(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "pass123"})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["accessToken"])
	assert.NotEmpty(t, resp["refreshToken"])
	assert.Equal(t, "test@example.com", resp["userEmail"])
}

func TestLoginInvalidPassword(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "wrong"})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func TestRefreshTokenSuccess(t *testing.T) {
	r := setupAuthRouter()
	w1 := httptest.NewRecorder()
	body1, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "pass123"})
	req1, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	var resp1 map[string]string
	json.Unmarshal(w1.Body.Bytes(), &resp1)
	refresh := resp1["refreshToken"]
	w2 := httptest.NewRecorder()
	body2, _ := json.Marshal(map[string]string{"refreshToken": refresh})
	req2, _ := http.NewRequest("POST", "/access-token/refresh", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)
	var resp2 map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.NotEmpty(t, resp2["accessToken"])
}
