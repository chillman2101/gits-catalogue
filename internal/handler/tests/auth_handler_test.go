package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chillman2101/gits-catalogue/internal/handler"
	"github.com/chillman2101/gits-catalogue/internal/mocks"
	"github.com/chillman2101/gits-catalogue/internal/service"
	res_err "github.com/chillman2101/gits-catalogue/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthRouter(svc *mocks.MockAuthService) *gin.Engine {
	r := gin.New()
	h := handler.NewAuthHandler(svc)
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.Refresh)
	r.POST("/auth/logout", func(c *gin.Context) {
		c.Set("user_id", float64(1))
		h.Logout(c)
	})
	return r
}

func TestAuthHandler_Register(t *testing.T) {
	t.Run("201 success", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		r := setupAuthRouter(svc)

		svc.On("Register", "test@example.com", "password123").Return(nil)

		body := `{"email":"test@example.com","password":"password123","confirm_password":"password123"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("400 password mismatch", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		r := setupAuthRouter(svc)

		body := `{"email":"test@example.com","password":"password123","confirm_password":"different"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "validation failed", resp["message"])
		// confirmpassword must match password (ToLower strips the underscore from struct field name)
		data := resp["data"].([]interface{})
		assert.Contains(t, data[0].(string), "confirmpassword")
		svc.AssertNotCalled(t, "Register")
	})

	t.Run("400 invalid email", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		r := setupAuthRouter(svc)

		body := `{"email":"not-an-email","password":"password123","confirm_password":"password123"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		svc.AssertNotCalled(t, "Register")
	})

	t.Run("400 password too short", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		r := setupAuthRouter(svc)

		body := `{"email":"test@example.com","password":"short","confirm_password":"short"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		svc.AssertNotCalled(t, "Register")
	})
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("200 success", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		r := setupAuthRouter(svc)

		pair := &service.TokenPair{AccessToken: "access-token", RefreshToken: "refresh-token"}
		svc.On("Login", "test@example.com", "password123").Return(pair, nil)

		body := `{"email":"test@example.com","password":"password123"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, "access-token", data["access_token"])
		assert.Equal(t, "refresh-token", data["refresh_token"])
		svc.AssertExpectations(t)
	})

	t.Run("401 invalid credentials", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		r := setupAuthRouter(svc)

		svc.On("Login", "test@example.com", "wrong").Return(nil, res_err.ErrUserInvalidCredentials)

		body := `{"email":"test@example.com","password":"wrong"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("400 missing fields", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		r := setupAuthRouter(svc)

		body := `{"email":"test@example.com"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		svc.AssertNotCalled(t, "Login")
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	t.Run("200 success", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		r := setupAuthRouter(svc)

		svc.On("Logout", uint(1)).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("500 service error", func(t *testing.T) {
		svc := new(mocks.MockAuthService)
		r := setupAuthRouter(svc)

		svc.On("Logout", uint(1)).Return(res_err.ErrUserNotFound)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})
}
