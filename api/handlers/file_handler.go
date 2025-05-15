package handlers

import (
	"fmt"
	"log"
	"net/http"
	"nomadule-backend/azure"
	"os"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	azureClient *azure.AzureStorageClient
}

func NewFileHandler(azureClient *azure.AzureStorageClient) *FileHandler {
	return &FileHandler{
		azureClient: azureClient,
	}
}

// POST /files/upload
func (h *FileHandler) UploadFile(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("File upload error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Failed to open uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process file"})
		return
	}
	defer file.Close()

	err = h.azureClient.UploadFile(user.ID, fileHeader.Filename, file)
	if err != nil {
		log.Printf("Azure upload failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	// Construct file URL
	fileURL := fmt.Sprintf(
		"https://%s.blob.core.windows.net/user-files/user-%s/%s",
		os.Getenv("AZURE_STORAGE_ACCOUNT"),
		user.ID,
		fileHeader.Filename,
	)

	contentType := fileHeader.Header.Get("Content-Type")

	// ✅ Respond with file metadata
	c.JSON(http.StatusOK, gin.H{
		"message":     "File uploaded successfully",
		"fileName":    fileHeader.Filename,
		"url":         fileURL,
		"contentType": contentType,
	})
}

// GET /files
func (h *FileHandler) ListFiles(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	files, err := h.azureClient.ListFiles(user.ID)
	if err != nil {
		log.Printf("Failed to list files for user %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not list files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// DELETE /files/:filename
func (h *FileHandler) DeleteFile(c *gin.Context) {
	user, err := extractUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	err = h.azureClient.DeleteFile(user.ID, filename)
	if err != nil {
		log.Printf("Delete failed for user %s file %s: %v", user.ID, filename, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
