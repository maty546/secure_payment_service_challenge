package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type authHandler struct {
	username string
	password string
}

var jwtKey = []byte("secret")

type IAuthHandler interface {
	GetToken(c *gin.Context)
	GetAuthMiddleware() gin.HandlerFunc
}

func NewLogin(username string, password string) IAuthHandler {
	return authHandler{username, password}
}

func (h authHandler) GetToken(c *gin.Context) {
	var creds LoginRequest

	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if creds.User != h.username || creds.Password != h.password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": creds.User,
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h authHandler) GetAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		tokenStr := authHeader[7:]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
