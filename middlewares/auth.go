package middlewares

import (
	"errors"
	"pronics-api/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(c *fiber.Ctx)(string,error){
	
	authorizationHeader := c.Get("Authorization")
	
	if !strings.Contains(authorizationHeader, "Bearer"){
		return "",errors.New("you have to use bearer")
	}

	tokenArray := strings.Split(authorizationHeader," ")
	if len(tokenArray) < 2{
		return "",errors.New("can't Find token")
	}

	tokenString := tokenArray[1]

	token, err := helper.ValidateToken(tokenString)

	if err != nil{
		return "",errors.New("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid{
		return "",errors.New("token is not valid")
	}

	id := claims["ID"]
	
	return id.(string), nil
}
