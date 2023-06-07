package middlewares

import (
	"errors"
	"net/http"
	"pronics-api/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(c *fiber.Ctx)error{
	
	authorizationHeader := c.Get("Authorization")
	
	if !strings.Contains(authorizationHeader, "Bearer"){
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", errors.New("you have to use bearer"))
		c.Status(http.StatusUnauthorized).JSON(response)
		return nil
	}

	tokenArray := strings.Split(authorizationHeader," ")
	if len(tokenArray) < 2{
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", errors.New("can't Find token"))
		c.Status(http.StatusUnauthorized).JSON(response)
		return nil
	}

	tokenString := tokenArray[1]

	token, err := helper.ValidateToken(tokenString)

	if err != nil{
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", errors.New("token is not valid"))
		c.Status(http.StatusUnauthorized).JSON(response)
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid{
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", errors.New("token is not valid"))
		c.Status(http.StatusUnauthorized).JSON(response)
		return nil
	}

	id := claims["ID"].(string)
	
	c.Locals("currentUserID",id)

	return c.Next()
}
