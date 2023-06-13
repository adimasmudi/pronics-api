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

// func (h *rekeningHandler) ChangeDetailRekening()