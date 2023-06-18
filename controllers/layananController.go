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

type layananHandler struct {
	layananService services.LayananService
}

func NewLayananHandler(layananService services.LayananService) *layananHandler {
	return &layananHandler{layananService}
}

func (h *layananHandler) Save(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.AddLayananInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add layanan failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	addedLayanan, err := h.layananService.SaveLayanan(ctx, input)

	if err != nil{
		response := helper.APIResponse("Add layanan Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add layanan success", http.StatusOK, "success", addedLayanan)
	c.Status(http.StatusOK).JSON(response)
	return nil
	
}

func (h *layananHandler) FindAll(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	allLayanan, err := h.layananService.FindAll(ctx)

	if err != nil{
		response := helper.APIResponse("Failed to get all layanan", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get all layanan success", http.StatusOK, "success", allLayanan)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *layananHandler) FindById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	layananId,_ := primitive.ObjectIDFromHex(c.Params("layananId"))

	layanan, err := h.layananService.FindById(ctx, layananId)

	if err != nil{
		response := helper.APIResponse("Failed to get layanan", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get layanan success", http.StatusOK, "success", layanan)
	c.Status(http.StatusOK).JSON(response)
	return nil
}