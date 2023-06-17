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

type bidangHandler struct {
	bidangService services.BidangService
}

func NewbidangHandler(bidangService services.BidangService) *bidangHandler{
	return &bidangHandler{bidangService}
}

func (h *bidangHandler) Save(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.AddBidangInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add bidang failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	var creatorId primitive.ObjectID
	if c.Locals("currentUserID").(string) != ""{
		currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))
		creatorId = currentUserId
	}

	addedbidang, err := h.bidangService.SaveBidang(ctx, input, creatorId)

	if err != nil{
		response := helper.APIResponse("Add bidang Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add bidang success", http.StatusOK, "success", addedbidang)
	c.Status(http.StatusOK).JSON(response)
	return nil
	
}

// get all bidang
func (h *bidangHandler) FindAll(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	allbidang, err := h.bidangService.FindAll(ctx)

	if err != nil{
		response := helper.APIResponse("Failed to get all bidang", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get all bidang success", http.StatusOK, "success", allbidang)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

// update bidang
func (h *bidangHandler) UpdateBidang(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	bidangId,_ := primitive.ObjectIDFromHex(c.Params("bidangId"))

	var input inputs.AddBidangInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Update bidang failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	updatedBidang, err := h.bidangService.UpdateBidang(ctx, currentUserId,bidangId, input)

	if err != nil{
		response := helper.APIResponse("Update bidang failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Update bidang success", http.StatusOK, "success", updatedBidang)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

// delete bidang
func (h *bidangHandler) DeleteBidang(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bidangId,_ := primitive.ObjectIDFromHex(c.Params("bidangId"))

	deletedBidang, err := h.bidangService.DeleteBidang(ctx, bidangId)

	if err != nil{
		response := helper.APIResponse("Delete bidang failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Delete bidang success", http.StatusOK, "success", deletedBidang)
	c.Status(http.StatusOK).JSON(response)
	return nil
}