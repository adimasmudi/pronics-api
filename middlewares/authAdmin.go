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
type adminAuth struct{
	adminRepository repositories.AdminRepository
}

func AdminAuth(adminRepository repositories.AdminRepository) *adminAuth{
	return &adminAuth{adminRepository}
}

func (a *adminAuth) AuthAdmin(c *fiber.Ctx) error {
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

	_, err = a.adminRepository.GetAdminById(ctx, id)

	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", "You're not an admin")
		c.Status(http.StatusUnauthorized).JSON(response)
		return nil
	}

	c.Locals("currentUserID",authId)

	return c.Next()
}