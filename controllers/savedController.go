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

type savedHandler struct {
	savedService services.SavedService
}

func NewSavedHandler(savedService services.SavedService) *savedHandler{
	return &savedHandler{savedService}
}

func (h *savedHandler) Save(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	
	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	mitraId, _ := primitive.ObjectIDFromHex(c.Params("mitraId"))
		

	addedSaved, err := h.savedService.Save(ctx, currentUserId, mitraId)

	if err != nil{
		response := helper.APIResponse("Add to saved failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add to saved success", http.StatusOK, "success", addedSaved)
	c.Status(http.StatusOK).JSON(response)
	return nil
}
