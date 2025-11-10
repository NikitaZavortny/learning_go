package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig конфигурация для CORS
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// DefaultCORSConfig дефолтная конфигурация CORS
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// CORS middleware
func CORS(config *CORSConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Проверяем разрешен ли origin
		if len(config.AllowOrigins) == 0 {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			for _, allowedOrigin := range config.AllowOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		// Устанавливаем остальные CORS заголовки
		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if len(config.AllowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
		}

		if len(config.AllowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		}

		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%.0f", config.MaxAge.Seconds()))
		}

		// Обрабатываем preflight запросы
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
