package handlers

import (
	"log"
	"mcdonalds-bot-backend/internal/models"
	"mcdonalds-bot-backend/internal/scheduler"
	"mcdonalds-bot-backend/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

type Handler struct {
	scheduler *scheduler.Scheduler
	hub       *websocket.Hub
}

func NewHandler(s *scheduler.Scheduler, h *websocket.Hub) *Handler {
	handler := &Handler{
		scheduler: s,
		hub:       h,
	}

	s.AddEventListener(func(event scheduler.Event) {
		h.Broadcast(event)
	})

	return handler
}

var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	h.hub.ServeWs(conn)
}

type CreateOrderRequest struct {
	ID   string `json:"id" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderType models.OrderType
	if req.Type == "VIP" {
		orderType = models.OrderTypeVIP
	} else {
		orderType = models.OrderTypeNormal
	}

	if err := h.scheduler.CreateOrder(req.ID, orderType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

type AddBotRequest struct {
	ID string `json:"id" binding:"required"`
}

func (h *Handler) AddBot(c *gin.Context) {
	var req AddBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.scheduler.AddBot(req.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Bot added successfully"})
}

type RemoveBotRequest struct {
	ID string `json:"id" binding:"required"`
}

func (h *Handler) RemoveBot(c *gin.Context) {
	var req RemoveBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.scheduler.RemoveBot(req.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bot removed successfully"})
}

func (h *Handler) GetBots(c *gin.Context) {
	bots := h.scheduler.GetBots()
	c.JSON(http.StatusOK, gin.H{"bots": bots})
}

func (h *Handler) GetOrders(c *gin.Context) {
	orders := h.scheduler.GetAllOrders()
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *Handler) GetPendingOrders(c *gin.Context) {
	orders := h.scheduler.GetPendingOrders()
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *Handler) GetCompletedOrders(c *gin.Context) {
	orders := h.scheduler.GetCompletedOrders()
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *Handler) GetStats(c *gin.Context) {
	stats := h.scheduler.GetStats()
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) Reset(c *gin.Context) {
	h.scheduler.Reset()
	c.JSON(http.StatusOK, gin.H{"message": "System reset successfully"})
}
