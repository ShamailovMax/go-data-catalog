package main

import (
	"log"
	"go-data-catalog/internal/config"
	"go-data-catalog/internal/handlers"
	"go-data-catalog/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	artifactRepo := postgres.NewArtifactRepository(db)
	
	artifactHandler := handlers.NewArtifactHandler(artifactRepo)
	
	r := gin.Default()

	r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
        c.Next()
    })
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
	
	artifacts := r.Group("/artifacts")
	{
		artifacts.GET("", artifactHandler.GetArtifacts)
		artifacts.POST("", artifactHandler.CreateArtifact)
	}
	
	log.Println("Server starting on :" + cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}