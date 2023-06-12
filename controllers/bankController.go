package controllers

import (
	"context"
	"net/http"
	"pronics-api/configs"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type bankHandler struct {
	bankService services.BankService
}

func NewBankHandler(bankService services.BankService) *bankHandler {
	return &bankHandler{bankService}
}

func (h *bankHandler) Save(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.AddBankInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add bank failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	fileName := ""

	logoBank, err := c.FormFile("logo_bank")

	if err != nil {
		response := helper.APIResponse("Add bank failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	if logoBank != nil{
		blobFile, err := logoBank.Open()
	
		if err != nil {
			response := helper.APIResponse("Add bank failed", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}

		fileName = helper.GenerateFilename(logoBank.Filename)
	
		err = configs.StorageInit("bank").UploadFile(blobFile, fileName)

		if err != nil {
			response := helper.APIResponse("Add bank failed", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}
	}

	addedBank, err := h.bankService.SaveBank(ctx, input, fileName)

	if err != nil {
		response := helper.APIResponse("Add bank failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add bank success", http.StatusOK, "success", addedBank)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

// get all
func (h *bankHandler) FindAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	allBank, err := h.bankService.FindAll(ctx)

	if err != nil{
		response := helper.APIResponse("Failed to get all bank", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get all bank success", http.StatusOK, "success", allBank)
	c.Status(http.StatusOK).JSON(response)
	return nil
}
