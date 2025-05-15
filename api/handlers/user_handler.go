package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"nomadule-backend/azure"
	"nomadule-backend/models"
	"nomadule-backend/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	repo        *repositories.UserRepository
	azureClient *azure.AzureStorageClient
}

func NewUserHandler(db *gorm.DB, azureClient *azure.AzureStorageClient) *UserHandler {
	return &UserHandler{
		repo:        repositories.NewUserRepository(db),
		azureClient: azureClient,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("JSON bind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Incoming request: %+v", req)

	if req.ClerkID == "" && req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required for manual sign-up"})
		return
	}

	// Check if user already exists
	existingUser, err := h.repo.FindByClerkID(req.ClerkID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error checking existing user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}
	if existingUser != nil {
		log.Println("User already exists:", existingUser)
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	var hashedPassword string
	if req.Password != "" {
		pw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Password hashing error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		hashedPassword = string(pw)
	}

	user := &models.User{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		ClerkID:   req.ClerkID,
		ImageURL:  req.ImageUrl,
		Password:  hashedPassword,
	}
	log.Printf("Creating user object: %+v", user)

	if err := h.repo.Create(user); err != nil {
		log.Println("Database create error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in database"})
		return
	}

	log.Println("User created successfully with ID:", user.ID)
	c.JSON(http.StatusCreated, gin.H{"id": user.ID})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update only allowed fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.ImageURL = req.ImageURL
	user.UpdatedAt = time.Now()

	if err := h.repo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UploadProfileImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File required"})
		return
	}
	defer file.Close()

	// Upload to Azure
	err = h.azureClient.UploadFile(id.String(), header.Filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	blobPath := fmt.Sprintf("user-%s/%s", id.String(), header.Filename)
	imageURL := fmt.Sprintf("https://%s.blob.core.windows.net/user-files/%s",
		os.Getenv("AZURE_STORAGE_ACCOUNT"),
		blobPath,
	)

	// Update ImageURL in DB
	err = h.repo.UpdateImageURL(id, imageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile image uploaded", "image_url": imageURL})
}
