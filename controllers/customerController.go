package controllers

import (
	"context"
	"net/http"
	"pronics-api/helper"
	"pronics-api/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type customerHandler struct {
	customerService services.CustomerService
}

func NewCustomerHandler(customerService services.CustomerService) *customerHandler{
	return &customerHandler{customerService}
}

func (h *customerHandler) GetProfile(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	customer, err := h.customerService.GetCustomerProfile(ctx,currentUserId)

	if err != nil{
		response := helper.APIResponse("Can't get customer profile", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("get profil customer success", http.StatusOK, "success", customer)
	c.Status(http.StatusOK).JSON(response)
	return nil

}