package main

import (
	"log"
	"mcdonalds-bot-backend/internal/handlers"
	"mcdonalds-bot-backend/internal/scheduler"
	"mcdonalds-bot-backend/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	s := scheduler.NewScheduler()
	hub := websocket.NewHub()
	h := handlers.NewHandler(s, hub)

	go hub.Run()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/ws", h.HandleWebSocket)

	api := r.Group("/api")
	{
		api.POST("/orders", h.CreateOrder)
		api.GET("/orders", h.GetOrders)
		api.GET("/orders/pending", h.GetPendingOrders)
		api.GET("/orders/completed", h.GetCompletedOrders)

		api.POST("/bots", h.AddBot)
		api.DELETE("/bots", h.RemoveBot)
		api.GET("/bots", h.GetBots)

		api.GET("/stats", h.GetStats)
		api.POST("/reset", h.Reset)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
