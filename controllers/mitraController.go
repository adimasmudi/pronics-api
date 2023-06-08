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

type mitraHandler struct {
	mitraService services.MitraService
}

func NewMitraHandler(mitraService services.MitraService) *mitraHandler{
	return &mitraHandler{mitraService}
}

func (h *mitraHandler) GetProfile(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	mitra, err := h.mitraService.GetMitraProfile(ctx,currentUserId)

	if err != nil{
		response := helper.APIResponse("Can't get mitra profile", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("get profil mitra success", http.StatusOK, "success", mitra)
	c.Status(http.StatusOK).JSON(response)
	return nil

}

