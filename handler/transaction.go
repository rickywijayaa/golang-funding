package handler

import (
	"funding/helper"
	"funding/transaction"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	Service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignsTransaction(c *gin.Context) {
	var input transaction.GetTransactionsByCampaignIdInput

	err := c.ShouldBindUri(input.ID)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := helper.ErrorMessageResponse(errors)

		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed To Get Campaign Transaction",
			http.StatusBadRequest,
			errorMessage,
		))
		return
	}

	transaction, err := h.Service.GetTransactionsByCampaignID(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to get campaigns transaction",
			http.StatusBadRequest,
			nil,
		))
		return
	}

	formatter := transaction

	c.JSON(http.StatusOK, helper.APIResponse(
		"List of campaigns",
		http.StatusOK,
		formatter,
	))
}
