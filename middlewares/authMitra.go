package middlewares

import (
	"context"
	"net/http"
	"pronics-api/helper"
	"pronics-api/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type mitraAuth struct{
	mitraRepository repositories.MitraRepository
}

func MitraAuth(mitraRepository repositories.MitraRepository) *mitraAuth{
	return &mitraAuth{mitraRepository}
}

func (a *mitraAuth) AuthMitra(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	authId, err := Auth(c)

	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", err.Error())
		c.Status(http.StatusUnauthorized).JSON(response)
		return nil
	}

	id, err := primitive.ObjectIDFromHex(authId)

	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", err.Error())
		c.Status(http.StatusUnauthorized).JSON(response)
		return nil
	}

	_, err = a.mitraRepository.GetMitraByIdUser(ctx, id)

	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", "You're not the mitra")
		c.Status(http.StatusUnauthorized).JSON(response)
		return nil
	}

	c.Locals("currentUserID",authId)

	return c.Next()
}