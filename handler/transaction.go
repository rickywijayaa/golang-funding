package handler

import (
	"funding/helper"
	"funding/transaction"
	"funding/user"
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

	err := c.ShouldBindUri(&input)
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

	currentUser := c.MustGet("current_user").(user.User)
	input.User = currentUser

	transactions, err := h.Service.GetTransactionsByCampaignID(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to get campaigns transaction",
			http.StatusBadRequest,
			nil,
		))
		return
	}

	formatter := transaction.CampaignTransactionsFormat(transactions)

	c.JSON(http.StatusOK, helper.APIResponse(
		"List of campaigns transactions",
		http.StatusOK,
		formatter,
	))
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("current_user").(user.User)
	userID := currentUser.ID

	transactions, err := h.Service.GetTransactionsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to get transaction",
			http.StatusBadRequest,
			nil,
		))
		return
	}

	formatter := transaction.UserTransactionsFormat(transactions)

	c.JSON(http.StatusOK, helper.APIResponse(
		"List of transactions",
		http.StatusOK,
		formatter,
	))
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := helper.ErrorMessageResponse(errors)

		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed To Create Transaction",
			http.StatusBadRequest,
			errorMessage,
		))
		return
	}

	currentUser := c.MustGet("current_user").(user.User)
	input.User = currentUser

	transactions, err := h.Service.CreateTransaction(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed To Create Transaction",
			http.StatusBadRequest,
			nil,
		))
		return
	}

	formatter := transaction.TransactionFormat(transactions)

	c.JSON(http.StatusOK, helper.APIResponse(
		"Success Create Transaction",
		http.StatusOK,
		formatter,
	))
}
