package controllers

import (
	"context"
	"net/http"
	"pronics-api/constants"
	"pronics-api/helper"
	"pronics-api/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type orderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *orderHandler{
	return &orderHandler{orderService}
}

func (h *orderHandler) CreateTemporaryOrder(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	mitraId, _:= primitive.ObjectIDFromHex(c.Params("mitraId"))

	addedTemporaryOrder, err := h.orderService.CreateTemporaryOrder(ctx, currentUserId, mitraId)

	if err != nil{
		if err.Error() == constants.TemporaryOrderExistMessage{
			response := helper.APIResponse(constants.TemporaryOrderExistMessage, http.StatusOK, "success", addedTemporaryOrder)
			c.Status(http.StatusOK).JSON(response)
			return nil
		}
		response := helper.APIResponse("Add temporary order Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Added temporary order success", http.StatusOK, "success", addedTemporaryOrder)
	c.Status(http.StatusOK).JSON(response)
	return nil
}