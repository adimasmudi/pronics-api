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

type komentarHandler struct {
	komentarService services.KomentarService
}

func NewKomentarHandler(komentarService services.KomentarService) *komentarHandler {
	return &komentarHandler{komentarService}
}

func (h *komentarHandler) AddKomentar(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))
	orderId, _ := primitive.ObjectIDFromHex(c.Params("orderId"))

	var input inputs.KomentarInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Add komentar failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	form, err := c.MultipartForm()

	if err != nil{
		response := helper.APIResponse("Error upload files", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	images := form.File["gambar_komentar"]

	fileNames := []string{}

	for _, image := range images{
		blobFile, err := image.Open()

		if err != nil{
			response := helper.APIResponse("Error upload image", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}

		fileName := helper.GenerateFilename(image.Filename)

		err = configs.StorageInit("komentar").UploadFile(blobFile, fileName)

		if err != nil{
			response := helper.APIResponse("Error upload image", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}

		fileNames = append(fileNames, fileName)
	}

	addedKomentar, err := h.komentarService.AddKomentar(ctx,currentUserId, orderId,input, fileNames)

	if err != nil{
		response := helper.APIResponse("Error add komentar", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Add komentar success", http.StatusOK, "success", addedKomentar)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *komentarHandler) KomentarDetail(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderId, _ := primitive.ObjectIDFromHex(c.Params("orderId"))

	detailKomentar, err := h.komentarService.SeeKomentar(ctx, orderId)

	if err != nil{
		response := helper.APIResponse("Get detail komentar failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Get detail komentarsuccess", http.StatusOK, "success", detailKomentar)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *komentarHandler) UpdateKomentar(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))
	komentarId, _ := primitive.ObjectIDFromHex(c.Params("komentarId"))

	var input inputs.KomentarInput

	if err := c.BodyParser(&input); err != nil {
		response := helper.APIResponse("Update komentar failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	form, err := c.MultipartForm()

	if err != nil{
		response := helper.APIResponse("Error upload files", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	images := form.File["gambar_komentar"]

	fileNames := []string{}

	for _, image := range images{
		blobFile, err := image.Open()

		if err != nil{
			response := helper.APIResponse("Error upload image", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}

		fileName := helper.GenerateFilename(image.Filename)

		err = configs.StorageInit("komentar").UploadFile(blobFile, fileName)

		if err != nil{
			response := helper.APIResponse("Error upload image", http.StatusBadRequest, "error", err.Error())
			c.Status(http.StatusBadRequest).JSON(response)
			return nil
		}

		fileNames = append(fileNames, fileName)
	}

	updatedKomentar, err := h.komentarService.UpdateKomentar(ctx,currentUserId, komentarId,input, fileNames)

	if err != nil{
		response := helper.APIResponse("Error update komentar", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Update komentar success", http.StatusOK, "success", updatedKomentar)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

func (h *komentarHandler) ResponseKomentar(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentUserId, _ := primitive.ObjectIDFromHex(c.Locals("currentUserID").(string))
	komentarId, _ := primitive.ObjectIDFromHex(c.Params("komentarId"))

	tipe := c.Query("type")

	likedKomentar, err := h.komentarService.ResponseKomentar(ctx, currentUserId, komentarId, tipe)

	if err != nil{
		response := helper.APIResponse("Error response komentar", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Response komentar success", http.StatusOK, "success", likedKomentar)
	c.Status(http.StatusOK).JSON(response)
	return nil
}

// delete bidang
func (h *komentarHandler) DeleteKomentar(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	komentarId,_ := primitive.ObjectIDFromHex(c.Params("komentarId"))

	deletedKomentar, err := h.komentarService.DeleteKomentar(ctx, komentarId)

	if err != nil{
		response := helper.APIResponse("Delete komentar failed", http.StatusBadRequest, "error", err.Error())
		c.Status(http.StatusBadRequest).JSON(response)
		return nil
	}

	response := helper.APIResponse("Delete komentar success", http.StatusOK, "success", deletedKomentar)
	c.Status(http.StatusOK).JSON(response)
	return nil
}