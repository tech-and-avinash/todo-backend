package routes

import (
	"nomadule-backend/api/handlers"
	"nomadule-backend/azure"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, azureClient *azure.AzureStorageClient) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://todo-nomadule.netlify.app", "https://todo.nomadule.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// User routes
	userHandler := handlers.NewUserHandler(db, azureClient)
	userGroup := r.Group("/users")
	{
		userGroup.POST("", userHandler.CreateUser)
		userGroup.GET("", userHandler.ListUsers)
		userGroup.GET("/:id", userHandler.GetUser)
		userGroup.PUT("/:id", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
		userGroup.POST("/:id/image", userHandler.UploadProfileImage)
	}

	// Note routes
	noteHandler := handlers.NewNoteHandler(db)
	noteGroup := r.Group("/notes")
	{
		noteGroup.POST("", noteHandler.CreateNote)
		noteGroup.GET("", noteHandler.GetAllNotes)
		noteGroup.GET("/:id", noteHandler.GetNoteByID)
		noteGroup.PUT("/:id", noteHandler.UpdateNote)
		noteGroup.DELETE("/:id", noteHandler.DeleteNote)
	}

	// File routes
	fileHandler := handlers.NewFileHandler(azureClient)
	fileGroup := r.Group("/files")
	{
		fileGroup.POST("/upload", fileHandler.UploadFile)
		fileGroup.GET("", fileHandler.ListFiles)
		fileGroup.DELETE("/:filename", fileHandler.DeleteFile)
	}

	contactHandler := handlers.NewContactHandler(db)
	contactGroup := r.Group("/contacts")
	{
		contactGroup.POST("", contactHandler.CreateContact)
		contactGroup.GET("", contactHandler.GetAllContacts)
		contactGroup.GET("/:id", contactHandler.GetContact)
		contactGroup.PUT("/:id", contactHandler.UpdateContact)
		contactGroup.DELETE("/:id", contactHandler.DeleteContact)
	}

	// // Auth routes
	// authHandler := handlers.NewAuthHandler(db)
	// authGroup := r.Group("/auth")
	// {
	// 	authGroup.POST("/login", authHandler.Login)
	// }
	return r
}
