package api

import (
	"github.com/gin-gonic/gin"
	"github.com/aryanc1027/api-kafka-golang/internal/config"
	"github.com/aryanc1027/api-kafka-golang/internal/kafka"
)


func SetupRouter(cfg *config.Config, producer kafka.Producer, consumer kafka.Consumer) *gin.Engine {
	router := gin.Default()

	rateLimiter := NewRateLimiter(1, 5) // 1 request per second with burst of 5
	router.Use(RateLimitMiddleware(rateLimiter))

	handler := NewHandler(producer, consumer)


	router.Use(corsMiddleware())
	router.Use(authMiddleware(cfg))


	api := router.Group("/api")
	{

		stream := api.Group("/stream")
		{
			stream.POST("/start", handler.StartStream)
			stream.POST("/:stream_id/send", handler.SendData)
			stream.GET("/:stream_id/results", handler.GetResults)
		}


	}

	return router
}


func corsMiddleware() gin.HandlerFunc {
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


func authMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(401, gin.H{"error": "API key is missing"})
			c.Abort()
			return
		}

		if apiKey != cfg.APIKey {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}