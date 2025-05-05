package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for login and other public routes
		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" {
			c.Next()
			return
		}

		// Get the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"details": "Token parsing failed",
				"action":  "Please login again"})
			c.Abort()
			return
		}

		// Check token validity and expiration
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Token is invalid",
				"action": "Please login again"})
			c.Abort()
			return
		}

		// Check expiration
		exp, ok := claims["exp"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Token expiration is invalid",
				"action": "Please login again"})
			c.Abort()
			return
		}

		// Compare expiration time
		if time.Now().Unix() > int64(exp) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Token has expired",
				"action": "Please login again"})
			c.Abort()
			return
		}

		// Token is valid, proceed
		c.Next()
	}
}
