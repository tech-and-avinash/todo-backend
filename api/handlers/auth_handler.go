package handlers

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var clerkClient clerk.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	secret := os.Getenv("CLERK_SECRET_KEY")
	if secret == "" {
		log.Fatal("CLERK_SECRET_KEY is not set in the environment")
	}

	clerkClient, err = clerk.NewClient(secret)
	if err != nil {
		log.Fatalf("Failed to initialize Clerk client: %v", err)
	}
}

func extractUserFromHeader(c *gin.Context) (*clerk.User, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return nil, errMissingToken()
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return nil, errInvalidToken()
	}

	token := parts[1]
	sessionClaims, err := clerkClient.VerifyToken(token)
	if err != nil {
		log.Printf("Token verification failed: %v", err)
		return nil, err
	}

	user, err := clerkClient.Users().Read(sessionClaims.Subject)
	if err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}

	return user, nil
}

func errMissingToken() error {
	return &gin.Error{
		Err:  http.ErrNoCookie,
		Type: gin.ErrorTypePrivate,
	}
}

func errInvalidToken() error {
	return &gin.Error{
		Err:  http.ErrNotSupported,
		Type: gin.ErrorTypePrivate,
	}
}
