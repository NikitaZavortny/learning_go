package config

import (
	"os"
	"strings"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	ServerPort       string
	Host_Mail        string
	Mail_Passwd      string
	CORSAllowOrigins string
}

func Load() *Config {
	return &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "1"),
		DBName:           getEnv("DB_NAME", "auth_db"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		Host_Mail:        getEnv("HOST_MAIL", ""),
		Mail_Passwd:      getEnv("MAIL_PASSWD", ""),
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "http://*"),
	}
}

func (c *Config) GetCORSAllowOrigins() []string {
	if c.CORSAllowOrigins == "" {
		return []string{"*"}
	}
	return strings.Split(c.CORSAllowOrigins, ",")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
