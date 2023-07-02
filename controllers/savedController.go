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

// get all saved mitra
func (h *savedHandler) ShowAllSaved(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	search := c.Query("search")
	daerah := c.Query("daerah")
	bidang := c.Query("bidang")
	urut := c.Query("urut")
	alamatCustomer := c.Query("alamatCustomer")

	searchFilter := make(map[string] string)

	searchFilter["search"] = search
	searchFilter["daerah"] = daerah
	searchFilter["bidang"] = bidang
	searchFilter["urut"] = urut
	searchFilter["alamatCustomer"] = alamatCustomer

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	katalogMitra, err := h.savedService.ShowAll(ctx, currentUserId, searchFilter)

	if err != nil {
		response := helper.APIResponse("Get all saved mitra failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get all saved mitra success", http.StatusOK, "success", katalogMitra)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

// delete from saved
func (h *savedHandler) DeleteSaved(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	savedId,_ := primitive.ObjectIDFromHex(c.Params("savedId"))

	deletedSaved, err := h.savedService.DeleteSaved(ctx, savedId)

	if err != nil{
		response := helper.APIResponse("Delete item from saved failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Delete item from saved success", http.StatusOK, "success", deletedSaved)
	c.Status(http.StatusOK).JSON(response)
	return nil
}
