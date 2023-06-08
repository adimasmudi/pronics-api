package controllers

import (
	"context"
	"net/http"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type wilayahCakupanHandler struct {
	wilayahCakupanService services.WilayahCakupanService
}

func NewwilayahCakupanHandler(wilayahCakupanService services.WilayahCakupanService) *wilayahCakupanHandler{
	return &wilayahCakupanHandler{wilayahCakupanService}
}

func (h *wilayahCakupanHandler) Save(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.AddWilayahCakupanInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add wilayah cakupan failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	addedwilayahCakupan, err := h.wilayahCakupanService.Save(ctx, input)

	if err != nil{
		response := helper.APIResponse("Add wilayah cakupan Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add wilayah cakupan success", http.StatusOK, "success", addedwilayahCakupan)
	c.Status(http.StatusOK).JSON(response)
	return nil
	
}

// get all wilayahCakupan
func (h *wilayahCakupanHandler) FindAll(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	allwilayahCakupan, err := h.wilayahCakupanService.FindAll(ctx)

	if err != nil{
		response := helper.APIResponse("Failed to get all wilayah cakupan", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get all wilayah success", http.StatusOK, "success", allwilayahCakupan)
	c.Status(http.StatusOK).JSON(response)
	return nil
}