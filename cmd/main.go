package main

import (
	"demo-basic-ai-chat-bot/web/app"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Create API v1 group
	v1 := router.Group("/api/v1")
	{
		// Register POST endpoint for /responses
		v1.POST("/responses", app.PostQuestion)
		v1.GET("/responses", app.GetBotResponse)
		v1.GET("/responses/version", app.CheckBotVersion)
		// Register POST endpoint for /openai/vision
		v1.POST("/openai/vision", app.PostOpenAIVision)
		v1.POST("/google/vision", app.GoogleChatService)
	}

	if err := router.Run(":8089"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
