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

type kategoriHandler struct {
	kategoriService services.KategoriBidangService
}

func NewKategoriHandler(kategoriService services.KategoriBidangService) *kategoriHandler{
	return &kategoriHandler{kategoriService}
}

func (h *kategoriHandler) Save(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.AddKategoriInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add kategori failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	addedKategori, err := h.kategoriService.Save(ctx, input)

	if err != nil{
		response := helper.APIResponse("Add kategori Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add kategori success", http.StatusOK, "success", addedKategori)
	c.Status(http.StatusOK).JSON(response)
	return nil
	
}

// get all kategori
func (h *kategoriHandler) FindAll(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	allKategori, err := h.kategoriService.FindAll(ctx)

	if err != nil{
		response := helper.APIResponse("Failed to get all category", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get all category success", http.StatusOK, "success", allKategori)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

// get all kategori with Bidang
func (h *kategoriHandler) GetKategoriWithBidang(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	allKategoriWithBidang, err := h.kategoriService.GetKategoriWithBidang(ctx)

	if err != nil{
		response := helper.APIResponse("Failed to get all category", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get all category success", http.StatusOK, "success", allKategoriWithBidang)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

// get kategori by id
func (h *kategoriHandler) GetKategoriById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	kategoriId, _ := primitive.ObjectIDFromHex(c.Params("kategoriId"))

	kategori, err := h.kategoriService.GetKategoriById(ctx, kategoriId)

	if err != nil{
		response := helper.APIResponse("Failed to get category", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get category success", http.StatusOK, "success", kategori)
	c.Status(http.StatusOK).JSON(response)
	return nil
}