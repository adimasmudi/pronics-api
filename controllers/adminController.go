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

type adminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) *adminHandler{
	return &adminHandler{adminService}
}

func (h *adminHandler) Register(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.RegisterAdminInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Register Failed", http.StatusBadRequest, "error", err)
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	registeredAdmin, err := h.adminService.Register(ctx, input)

	if err != nil{
		response := helper.APIResponse("Register Failed", http.StatusBadRequest, "error", err)
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Admin registration success", http.StatusOK, "success", registeredAdmin)
	c.Status(http.StatusOK).JSON(response)
	return nil
	
}

func (h *adminHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.LoginAdminInput

	//validate the request body
	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "error", err)
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	logedinAdmin, token,  err := h.adminService.Login(ctx,input)

	if err != nil{
		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "error", err)
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Login success", http.StatusOK, "success", &fiber.Map{"user" : logedinAdmin, "token" : token})
	c.Status(http.StatusOK).JSON(response)
	return nil
}