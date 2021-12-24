package helper

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Success string `json:"success"`
}

func APIResponse(message string, code int, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Success: "Success",
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func APIFailedResponse(message string, code int, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Success: "Failed",
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, fmt.Sprintf("Field %v is %v", e.Field(), e.ActualTag()))
	}

	return errors
}

func ErrorMessageResponse(errors []string) map[string]interface{} {
	return gin.H{"errors": errors}
}
