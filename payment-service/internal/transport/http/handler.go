package http

import (
	"net/http"
	"payment-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createPayment *usecase.CreatePayment
}

func NewHandler(createPayment *usecase.CreatePayment) *Handler {
	return &Handler{createPayment: createPayment}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/payments", h.CreatePayment)
	r.GET("/payments/:order_id", h.GetPayment)
}

func (h *Handler) CreatePayment(c *gin.Context) {
	var cmd usecase.CreatePaymentCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.createPayment.Execute(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetPayment(c *gin.Context) {
	orderID := c.Param("order_id")
	c.JSON(http.StatusOK, gin.H{"order_id": orderID, "status": "exists"})
}
