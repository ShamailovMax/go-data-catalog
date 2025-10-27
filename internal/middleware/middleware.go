package middleware

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
)

// LoggerMiddleware логирует все запросы с временем выполнения
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Засекаем время начала запроса
		startTime := time.Now()
		
		// Обрабатываем запрос
		c.Next()
		
		// Вычисляем время выполнения
		duration := time.Since(startTime)
		
		// Логируем информацию о запросе
		log.Printf(
			"[%s] %s %s - Status: %d - Duration: %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			duration,
		)
	}
}

// ErrorHandlerMiddleware обрабатывает ошибки и форматирует ответы
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// Проверяем, есть ли ошибки
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Логируем ошибку
			log.Printf("Error processing request: %v", err)
			
			// Определяем тип ошибки и возвращаем соответствующий статус
			switch err.Type {
			case gin.ErrorTypePublic:
				c.JSON(c.Writer.Status(), gin.H{
					"error": err.Error(),
				})
			case gin.ErrorTypeBind:
				c.JSON(400, gin.H{
					"error": "Invalid request data",
					"details": err.Error(),
				})
			default:
				c.JSON(500, gin.H{
					"error": "Internal server error",
				})
			}
			
			c.Abort()
		}
	}
}

// CORSMiddleware добавляет заголовки CORS для работы с фронтендом
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}