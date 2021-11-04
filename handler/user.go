package handler

import (
	"funding/helper"
	"funding/user"
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, err)
	}

	newUser, err := h.service.RegisterUser(input)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	formatter := user.UserFormat(newUser, "token")

	c.JSON(http.StatusOK, helper.APIResponse(
		"Your account has been created",
		http.StatusOK,
		"success",
		formatter,
	))
}
