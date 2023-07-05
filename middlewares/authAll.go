package middlewares

import (
	"context"
	"net/http"
	"pronics-api/helper"
	"time"

	"github.com/gofiber/fiber/v2"
)

type authAll struct{
	
}

func AuthAll() *authAll{
	return &authAll{}
}

func (a *authAll) AuthAll(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	authId, err := Auth(c)

	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", err.Error())
		c.Status(http.StatusUnauthorized).JSON(response)
		return nil
	}

	c.Locals("currentUserID",authId)

	return c.Next()
}