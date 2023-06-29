package controllers

import (
	"context"
	"net/http"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type orderDetailHandler struct {
	orderDetailService services.OrderDetailService
}

func NewOrderDetailHandler(orderDetailService services.OrderDetailService) *orderDetailHandler{
	return &orderDetailHandler{orderDetailService}
}

func (h *orderDetailHandler) AddOrUpdateOrderDetail(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderId, _:= primitive.ObjectIDFromHex(c.Params("orderId"))

	var input inputs.AddOrUpdateOrderDetailInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add or update order detail failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	AddOrUpdateOrderDetail, err := h.orderDetailService.AddOrUpdateOrderDetail(ctx, orderId, input)

	if err != nil{
		response := helper.APIResponse("Add or Update order detail failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add or update Update order detail success", http.StatusOK, "success", AddOrUpdateOrderDetail)
	c.Status(http.StatusOK).JSON(response)
	return nil
}