package handler

import (
	"fmt"
	"funding/auth"
	"funding/helper"
	"funding/user"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service     user.Service
	AuthService auth.Service
}

func NewUserHandler(service user.Service, authService auth.Service) *userHandler {
	return &userHandler{service, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := helper.ErrorMessageResponse(errors)

		c.JSON(http.StatusUnprocessableEntity, helper.APIFailedResponse(
			"Failed To Created Account",
			http.StatusUnprocessableEntity,
			errorMessage,
		))
		return
	}

	newUser, err := h.service.RegisterUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed To Created Account",
			http.StatusBadRequest,
			gin.H{"errors": err.Error()},
		))
		return
	}

	token, err := h.AuthService.GenerateToken(newUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed To Created Account",
			http.StatusBadRequest,
			gin.H{"errors": err.Error()},
		))
		return
	}

	formatter := user.UserFormat(newUser, token)

	c.JSON(http.StatusOK, helper.APIResponse(
		"Your account has been created",
		http.StatusOK,
		formatter,
	))
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := helper.ErrorMessageResponse(errors)

		c.JSON(http.StatusUnprocessableEntity, helper.APIFailedResponse(
			"Failed To Login",
			http.StatusUnprocessableEntity,
			errorMessage,
		))
		return
	}

	loggedInUser, err := h.service.LoginUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed To Login",
			http.StatusBadRequest,
			gin.H{"errors": err.Error()},
		))
		return
	}

	token, err := h.AuthService.GenerateToken(loggedInUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed To Login",
			http.StatusBadRequest,
			gin.H{"errors": err.Error()},
		))
		return
	}
	formatter := user.UserFormat(loggedInUser, token)

	c.JSON(http.StatusOK, helper.APIResponse(
		"Successfully Login",
		http.StatusOK,
		formatter,
	))
}

func (h *userHandler) IsEmailExist(c *gin.Context) {
	var input user.IsEmailExistInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := helper.ErrorMessageResponse(errors)

		c.JSON(http.StatusUnprocessableEntity, helper.APIFailedResponse(
			"Email checking failed",
			http.StatusUnprocessableEntity,
			errorMessage,
		))
		return
	}

	isEmailValid, err := h.service.IsEmailExist(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Email Checking Failed",
			http.StatusBadRequest,
			gin.H{"errors": "Internal Server Error"},
		))
		return
	}

	data := gin.H{
		"is_available": isEmailValid,
	}

	metaMessage := "Email has been registered"
	if isEmailValid {
		metaMessage = "Email is available"
	}

	c.JSON(http.StatusOK, helper.APIResponse(
		metaMessage,
		http.StatusOK,
		data,
	))
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to upload avatar image",
			http.StatusBadRequest,
			gin.H{"is_uploaded": false},
		))
		return
	}

	currentUser := c.MustGet("current_user").(user.User)
	userID := currentUser.ID

	path := "images"
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	path = fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to upload avatar image",
			http.StatusBadRequest,
			gin.H{"is_uploaded": false},
		))
		return
	}

	_, err = h.service.SaveAvatar(userID, path)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to upload avatar image",
			http.StatusBadRequest,
			gin.H{"is_uploaded": false},
		))
		return
	}

	c.JSON(http.StatusOK, helper.APIResponse(
		"Success Upload Avatar",
		http.StatusOK,
		gin.H{"is_uploaded": true},
	))
}
