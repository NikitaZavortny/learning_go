package handler

import (
	"auth-server/internal/model"
	"auth-server/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtSecret   string
}

func NewAuthHandler(authService *service.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:   jwtSecret,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем refresh token в httpOnly cookie
	h.setRefreshTokenCookie(c, tokens.RefreshToken)

	// Возвращаем только access token в теле ответа
	c.JSON(http.StatusCreated, gin.H{
		"access_token": tokens.AccessToken,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем refresh token в httpOnly cookie
	h.setRefreshTokenCookie(c, tokens.RefreshToken)

	// Возвращаем только access token в теле ответа
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	// Получаем refresh token из cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token not found"})
		return
	}

	tokens, err := h.authService.RefreshTokens(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Обновляем refresh token в cookie
	h.setRefreshTokenCookie(c, tokens.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Удаляем refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	c.SetCookie("refresh_token", token, 7*24*3600, "/", "", false, true) // httpOnly, 7 дней
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
	})
}

func (h *AuthHandler) Activate(c *gin.Context) {

	link := c.Param("link")
	fmt.Println(link)
	_, err := h.authService.Activate(link)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "activated",
	})
}
