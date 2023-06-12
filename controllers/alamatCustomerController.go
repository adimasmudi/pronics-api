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

type alamatCustomerHandler struct {
	alamatCustomerService services.AlamatCustomerService
}

func NewAlamatCustomerHandler(alamatCustomerService services.AlamatCustomerService) *alamatCustomerHandler{
	return &alamatCustomerHandler{alamatCustomerService}
}

func (h *alamatCustomerHandler) Save(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	var input inputs.AddAlamatCustomerInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add alamat failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	addedAlamat, err := h.alamatCustomerService.SaveAlamat(ctx, input, currentUserId)

	if err != nil{
		response := helper.APIResponse("Add alamat Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add alamat success", http.StatusOK, "success", addedAlamat)
	c.Status(http.StatusOK).JSON(response)
	return nil
	
}