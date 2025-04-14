package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"ymir/api/controllers"
	"ymir/api/db"
	"ymir/api/middleware"
	"ymir/api/models"
	"ymir/api/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/brianvoe/gofakeit/v7"
)

func TestMain(m *testing.M) {
	db.ConnectToDb()
	db.SyncDatabase()

	m.Run()
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestNotAuthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetUpRouter()

	router.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)

	req, _ := http.NewRequest("GET", "/auth/validate", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestAuthenticatedWithInvalidToken(t *testing.T) {
	router := SetUpRouter()

	router.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)

	req, _ := http.NewRequest("GET", "/auth/validate", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthenticatedWithValidTokenButUserNotFound(t *testing.T) {
	router := SetUpRouter()

	router.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)

	req, _ := http.NewRequest("GET", "/auth/validate", nil)

	token, _ := utils.GenerateToken(models.User{Email: "ymir@duckduckgo.fr", Password: "bird"})

	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthenticatedWithValidTokenAndUserFound(t *testing.T) {
	router := SetUpRouter()

	router.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)

	user := models.User{Email: gofakeit.Email(), Username: gofakeit.Username(), Password: gofakeit.Password(true, true, true, true, false, 14)}

	db.DB.Create(&user)

	token, _ := utils.GenerateToken(user)

	print(user.ID)

	req, _ := http.NewRequest("GET", "/auth/validate", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	print(resp.Body.String())
	print(resp.Code)

	assert.Equal(t, http.StatusOK, resp.Code)

	db.DB.Delete(&user)
}

func TestFullAuthentificationFlowWithCookie(t *testing.T) {
	router := SetUpRouter()

	router.POST("/auth/login", controllers.Login)
	router.GET("/auth/validate", middleware.RequireAuth, controllers.Validate)

	user := models.User{Email: gofakeit.Email(), Username: gofakeit.Username(), Password: gofakeit.Password(true, true, true, true, false, 14)}

	db.DB.Create(&user)

	/*
		LOGIN
	*/

	loginRequest := models.LoginRequest{Username: user.Username, Password: user.Password}
	jsonValue, _ := json.Marshal(loginRequest)

	// Login user
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var tokenResponse models.TokenResponse
	json.NewDecoder(resp.Body).Decode(&tokenResponse)
	assert.NotEmpty(t, tokenResponse.Token)

	// Validate token locally
	token, err := utils.ParseToken(tokenResponse.Token)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	// Validate user
	req, _ = http.NewRequest("GET", "/auth/validate", nil)
	req.Header.Set("Cookie", resp.Header().Get("Set-Cookie")) // Copy cookie from login response

	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	db.DB.Delete(&user)
}
