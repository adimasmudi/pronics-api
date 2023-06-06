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

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *userHandler{
	return &userHandler{userService}
}

func (h *userHandler) Register(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.RegisterUserInput

	if err := c.BodyParser(&input); err != nil {
		errorMessage := &fiber.Map{
			"Error": err.Error(),
		}
		response := helper.APIResponse("Register Failed", http.StatusBadRequest, "error", errorMessage)
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	registeredUser, err := h.userService.Register(ctx, input)

	if err != nil{
		errorMessage := &fiber.Map{
			"Error": err.Error(),
		}
		response := helper.APIResponse("Register Failed", http.StatusBadRequest, "error", errorMessage)
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("User registration success", http.StatusOK, "success", registeredUser)
	c.Status(http.StatusOK).JSON(response)
	return nil
	
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.LoginUserInput

	//validate the request body
	if err := c.BodyParser(&input); err != nil {
		errorMessage := &fiber.Map{
			"Error": err.Error(),
		}
		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "error", errorMessage)
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	_, token,  err := h.userService.Login(ctx,input)

	if err != nil{
		errorMessage := &fiber.Map{
			"Error": err.Error(),
		}
		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "error", errorMessage)
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Login success", http.StatusOK, "success", &fiber.Map{ "token" : token})
	c.Status(http.StatusOK).JSON(response)
	return nil
}