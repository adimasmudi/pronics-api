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

type rekeningHandler struct {
	rekeningService services.RekeningService
}

func NewRekeningHandler(rekeningService services.RekeningService) *rekeningHandler {
	return &rekeningHandler{rekeningService}
}

func (h *rekeningHandler) GetDetailRekening(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	rekening, err := h.rekeningService.GetDetailRekening(ctx, currentUserId)

	if err != nil {
		response := helper.APIResponse("Can't get rekening detail", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("get rekening detail success", http.StatusOK, "success", rekening)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *rekeningHandler) ChangeDetailRekening(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	var input inputs.UpdateRekeningInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Update rekening failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	updatedRekening, err := h.rekeningService.UpdateRekening(ctx, currentUserId, input)

	if err != nil{
		response := helper.APIResponse("Update rekening failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Update rekening success", http.StatusOK, "success", updatedRekening)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *rekeningHandler) AddRekening(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	var input inputs.UpdateRekeningInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add rekening failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	addedRekening, err := h.rekeningService.SaveRekening(ctx, currentUserId, input)

	if err != nil{
		response := helper.APIResponse("Add rekening failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add rekening success", http.StatusOK, "success", addedRekening)
	c.Status(http.StatusOK).JSON(response)
	return nil
}