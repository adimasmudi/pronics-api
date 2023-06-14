package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"pronics-api/configs"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		response := helper.APIResponse("Register Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	registeredUser, err := h.userService.Register(ctx, input)

	if err != nil{
		response := helper.APIResponse("Register Failed", http.StatusBadRequest, "error", err.Error())
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
		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	user, token,  err := h.userService.Login(ctx,input)

	if err != nil{
		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	userFormat := formatters.FormatUser(user)
	fmt.Println(user)
	fmt.Println(userFormat)
	response := helper.APIResponse("Login success", http.StatusOK, "success", &fiber.Map{ "user" : userFormat,"token" : token})
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *userHandler) Callback(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if c.FormValue("state") != os.Getenv("oAuth_String") {
		response := helper.APIResponse("Can't login to your account", http.StatusBadRequest, "error", errors.New("token login invalid"))
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	token, err := configs.GoogleOAuthConfig().Exchange(context.Background(), c.FormValue("code"))
	if err != nil {
		response := helper.APIResponse("code exchange failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		response := helper.APIResponse("failed getting user info", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		response := helper.APIResponse("Failed reading response body", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	var googleUser helper.GoogleUser

	json.Unmarshal([]byte(string(contents)), &googleUser)

	user,loginToken, err := h.userService.Signup(ctx,googleUser)

	if err != nil{
		response := helper.APIResponse("Signup User Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}
	userFormat := formatters.FormatUser(user)
	responses := helper.APIResponse("Signup User Success", http.StatusOK, "success", &fiber.Map{"user" : userFormat,"token" : loginToken})
	c.Status(http.StatusOK).JSON(responses)
	return nil

}

func (h *userHandler) RegisterMitra(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input inputs.RegisterMitraInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Register Mitra Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	eKTP, err := c.FormFile("e_ktp")

	if err != nil{
		response := helper.APIResponse("Register Mitra Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	blobFile, err := eKTP.Open()

	if err != nil{
		response := helper.APIResponse("Register Mitra Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	// generate kode unik untuk fileName
	fileName := helper.GenerateFilename(eKTP.Filename)
	
	err = configs.StorageInit("ektp").UploadFile(blobFile, fileName)

	if err != nil{
		response := helper.APIResponse("Register Mitra Failed Upload Image E-ktp", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	registeredUser, err := h.userService.RegisterMitra(ctx, input, fileName)

	if err != nil{
		response := helper.APIResponse("Register Mitra Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Mitra registration success", http.StatusOK, "success", registeredUser)
	c.Status(http.StatusOK).JSON(response)
	return nil
	
}

func (h *userHandler) ChangePassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))

	var input inputs.ChangePasswordUserInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Change Password Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	updatedPassword, err := h.userService.ChangePassword(ctx, currentUserId, input)

	if err != nil{
		response := helper.APIResponse("Change Password Failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Change Password success", http.StatusOK, "success", updatedPassword)
	c.Status(http.StatusOK).JSON(response)
	return nil
}