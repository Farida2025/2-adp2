package handler

import (
	"net/http"
	"order-service/internal/usecase"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createOrder *usecase.CreateOrder
	getOrder    *usecase.GetOrder
	cancelOrder *usecase.CancelOrder
}

func NewHandler(co *usecase.CreateOrder, goUc *usecase.GetOrder, ca *usecase.CancelOrder) *Handler {
	return &Handler{
		createOrder: co,
		getOrder:    goUc,
		cancelOrder: ca,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/orders", h.CreateOrder)
	r.GET("/orders/:id", h.GetOrder)
	r.PATCH("/orders/:id/cancel", h.CancelOrder)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var cmd usecase.CreateOrderCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	orderID, err := h.createOrder.Execute(c.Request.Context(), cmd)
	if err != nil {
		if strings.Contains(err.Error(), "payment") || err.Error() == "payment declined" {
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
	order, err := h.getOrder.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *Handler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	err := h.cancelOrder.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order cancelled successfully"})
}
