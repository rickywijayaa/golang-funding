package handler

import (
	"funding/helper"
	"funding/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *userHandler {
	return &userHandler{service}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}

		errorMessage := gin.H{"errors": errors}

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
			nil,
		))
		return
	}

	formatter := user.UserFormat(newUser, "token")

	c.JSON(http.StatusOK, helper.APIResponse(
		"Your account has been created",
		http.StatusOK,
		"success",
		formatter,
	))
}
