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

type layananMitraHandler struct {
	layananMitraService services.LayananMitraService
}

func NewLayananMitraHandler(layananMitraService services.LayananMitraService) *layananMitraHandler{
	return &layananMitraHandler{layananMitraService}
}

func (h *layananMitraHandler) Save(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.AddLayananInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add layanan mitra failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	var creatorId primitive.ObjectID
	if c.Locals("currentUserID").(string) != ""{
		currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))
		creatorId = currentUserId
	}

	addedLayananMitra, err := h.layananMitraService.Save(ctx, input, creatorId)

	if err != nil{
		response := helper.APIResponse("Add layanan mitra Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add layanan mitra success", http.StatusOK, "success", addedLayananMitra)
	c.Status(http.StatusOK).JSON(response)
	return nil
	
}

func (h *layananMitraHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	layananMitraId,_ := primitive.ObjectIDFromHex(c.Params("layananMitraId"))

	deletedLayananMitra, err := h.layananMitraService.Delete(ctx, layananMitraId)

	if err != nil{
		response := helper.APIResponse("Delete Layanan Mitra failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Delete Layanan Mitra success", http.StatusOK, "success", deletedLayananMitra)
	c.Status(http.StatusOK).JSON(response)
	return nil
}