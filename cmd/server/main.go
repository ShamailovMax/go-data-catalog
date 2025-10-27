package main

import (
	"log"
	"go-data-catalog/internal/config"
	"go-data-catalog/internal/handlers"
	"go-data-catalog/internal/middleware"
	"go-data-catalog/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	
	// Инициализация БД
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	// Инициализация репозиториев
	artifactRepo := postgres.NewArtifactRepository(db)
	contactRepo := postgres.NewContactRepository(db)
	artifactFieldRepo := postgres.NewArtifactFieldRepository(db)
	
	// Инициализация handlers
	artifactHandler := handlers.NewArtifactHandler(artifactRepo)
	contactHandler := handlers.NewContactHandler(contactRepo)
	artifactFieldHandler := handlers.NewArtifactFieldHandler(artifactFieldRepo, artifactRepo)
	
	// Настройка роутера
	r := gin.New() // Используем New вместо Default чтобы сами настроить middleware
	
	// Добавляем наши middleware
	r.Use(gin.Recovery()) // Восстановление после паники
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.CORSMiddleware())
	
	// Глобальный middleware для установки Content-Type
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})
	
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
	
	// API v1 группа
	v1 := r.Group("/api/v1")
	{
		// Артефакты
		artifacts := v1.Group("/artifacts")
		{
			artifacts.GET("", artifactHandler.GetArtifacts)
			artifacts.GET("/:id", artifactHandler.GetArtifactByID)
			artifacts.POST("", artifactHandler.CreateArtifact)
			artifacts.PUT("/:id", artifactHandler.UpdateArtifact)
			artifacts.DELETE("/:id", artifactHandler.DeleteArtifact)
			// Поля артефактов (вложенные маршруты)
			artifactFields := artifacts.Group("/:id/fields")
			{
				artifactFields.GET("", artifactFieldHandler.GetFieldsByArtifact)
				artifactFields.POST("", artifactFieldHandler.CreateField)
			}
		}
		
		// Контакты
		contacts := v1.Group("/contacts")
		{
			contacts.GET("", contactHandler.GetContacts)
			contacts.GET("/:id", contactHandler.GetContactByID)
			contacts.POST("", contactHandler.CreateContact)
			contacts.PUT("/:id", contactHandler.UpdateContact)
			contacts.DELETE("/:id", contactHandler.DeleteContact)
		}
		// Отдельные маршруты для работы с полями по id
		fields := v1.Group("/fields")
		{
			fields.GET("/:id", artifactFieldHandler.GetFieldByID)
			fields.PUT("/:id", artifactFieldHandler.UpdateField)
			fields.DELETE("/:id", artifactFieldHandler.DeleteField)
		}
	}
	
	log.Println("Server starting on :" + cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
