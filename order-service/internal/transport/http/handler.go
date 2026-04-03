package http

import (
	"net/http"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createOrder     *usecase.CreateOrder
	getRecentOrders *usecase.GetRecentOrders
}

func NewHandler(createOrder *usecase.CreateOrder, getRecentOrders *usecase.GetRecentOrders) *Handler {
	return &Handler{createOrder: createOrder,
		getRecentOrders: getRecentOrders}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/orders", h.CreateOrder)
	r.GET("/orders/:id", h.GetOrder)
	r.PATCH("/orders/:id/cancel", h.CancelOrder)
	r.GET("/orders/recent", h.GetRecentOrders)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var cmd usecase.CreateOrderCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	orderID, err := h.createOrder.Execute(c.Request.Context(), cmd)
	if err != nil {
		if err.Error() == "payment service unavailable" || err.Error() == "payment declined" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order_id": orderID, "status": "Paid"})
}

func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"order_id": id,
		"message":  "order details (simplified)",
	})
}

func (h *Handler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message":  "order cancelled (if pending)",
		"order_id": id,
	})
}

func (h *Handler) GetRecentOrders(c *gin.Context) {
	var cmd usecase.GetRecentOrderCommand
	if err := c.ShouldBindQuery(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	orders, err := h.getRecentOrders.Execute(c.Request.Context(), cmd)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)

}
