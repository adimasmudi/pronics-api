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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// bank detail
func (h *bankHandler) FindById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bankId,_ := primitive.ObjectIDFromHex(c.Params("bankId"))

	bank, err := h.bankService.GetDetail(ctx, bankId)

	if err != nil{
		response := helper.APIResponse("Failed to get detail bank", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get detail bank success", http.StatusOK, "success", bank)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

// update bank
func (h *bankHandler) UpdateBank(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bankId,_ := primitive.ObjectIDFromHex(c.Params("bankId"))

	var input inputs.AddBankInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Update bank failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	fileName := ""

	logoBank, _:= c.FormFile("logo_bank")

	if logoBank != nil{
		blobFile, err := logoBank.Open()
	
		if err != nil {
			response := helper.APIResponse("Update bank failed", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}

		fileName = helper.GenerateFilename(logoBank.Filename)
	
		err = configs.StorageInit("bank").UploadFile(blobFile, fileName)

		if err != nil {
			response := helper.APIResponse("Update bank failed", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}
	}

	updatedBank, err := h.bankService.UpdateBank(ctx, bankId, input, fileName)

	if err != nil {
		response := helper.APIResponse("Update bank failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Update bank success", http.StatusOK, "success", updatedBank)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

// delete bank
