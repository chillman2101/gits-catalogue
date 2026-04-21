package handler

import (
	"github.com/chillman2101/gits-catalogue/internal/response"
	"github.com/chillman2101/gits-catalogue/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

type registerRequest struct {
	Email           string `json:"email" binding:"required,email,max=100"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      registerRequest   true  "Register data"
// @Success      201   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	if err := h.service.Register(req.Email, req.Password); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, "registered successfully", nil)
}

// Login godoc
// @Summary      Login and get JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      loginRequest      true  "Login credentials"
// @Success      200   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	pair, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.OK(c, "login successful", gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh godoc
// @Summary      Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      refreshRequest    true  "Refresh token"
// @Success      200   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	pair, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.OK(c, "token refreshed successfully", gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
}

// Logout godoc
// @Summary      Logout and invalidate token
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := uint(c.GetFloat64("user_id"))
	if err := h.service.Logout(userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "logged out successfully", nil)
}
