package handler

import (
	"fmt"
	"funding/campaign"
	"funding/helper"
	"funding/user"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	Service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.Service.GetCampaigns(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to get campaigns",
			http.StatusBadRequest,
			nil,
		))
		return
	}

	formatter := campaign.CampaignsFormat(campaigns)

	c.JSON(http.StatusOK, helper.APIResponse(
		"List of campaigns",
		http.StatusOK,
		formatter,
	))
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := helper.ErrorMessageResponse(errors)

		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed To Get Detail Campaign",
			http.StatusBadRequest,
			errorMessage,
		))
		return
	}

	campaignDetail, err := h.Service.GetCampaignByID(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed To Get Detail Campaign",
			http.StatusBadRequest,
			gin.H{"errors": err.Error()},
		))
		return
	}

	formatter := campaign.CampaignDetailFormat(campaignDetail)

	c.JSON(http.StatusOK, helper.APIResponse(
		"Success Get Detail Campaign",
		http.StatusOK,
		formatter,
	))
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := helper.ErrorMessageResponse(errors)

		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed Create Campaign",
			http.StatusBadRequest,
			errorMessage,
		))
		return
	}

	currentUser := c.MustGet("current_user").(user.User)
	input.User = currentUser

	newCampaign, err := h.Service.CreateCampaign(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed Create Campaign",
			http.StatusBadRequest,
			gin.H{"errors": err.Error()},
		))
		return
	}

	formatter := campaign.CampaignFormat(newCampaign)

	c.JSON(http.StatusOK, helper.APIResponse(
		"Success Create Campaign",
		http.StatusOK,
		formatter,
	))
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := helper.ErrorMessageResponse(errors)

		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed To Update Campaign",
			http.StatusBadRequest,
			errorMessage,
		))
		return
	}

	var input campaign.CreateCampaignInput
	currentUser := c.MustGet("current_user").(user.User)
	input.User = currentUser

	err = c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := helper.ErrorMessageResponse(errors)

		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to Update Campaign",
			http.StatusBadRequest,
			errorMessage,
		))
		return
	}

	updatedCampaign, err := h.Service.UpdateCampaign(inputID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to Update Campaign",
			http.StatusBadRequest,
			gin.H{"errors": err.Error()},
		))
		return
	}

	formatter := campaign.CampaignFormat(updatedCampaign)

	c.JSON(http.StatusOK, helper.APIResponse(
		"Success Update Campaign",
		http.StatusOK,
		formatter,
	))
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput
	err := c.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := helper.ErrorMessageResponse(errors)

		c.JSON(http.StatusUnprocessableEntity, helper.APIFailedResponse(
			"Failed to upload campaign image",
			http.StatusUnprocessableEntity,
			errorMessage,
		))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to upload campaign image",
			http.StatusBadRequest,
			gin.H{"is_uploaded": false},
		))
		return
	}

	currentUser := c.MustGet("current_user").(user.User)
	userID := currentUser.ID
	input.User = currentUser

	path := "images"
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	path = fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to upload campaign image",
			http.StatusBadRequest,
			gin.H{"is_uploaded": false},
		))
		return
	}

	_, err = h.Service.SaveCampaignImage(input, path)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.APIFailedResponse(
			"Failed to upload campaign image due not owning this campaign",
			http.StatusBadRequest,
			gin.H{"is_uploaded": false},
		))
		return
	}

	c.JSON(http.StatusOK, helper.APIResponse(
		"Success upload campaign image",
		http.StatusOK,
		gin.H{"is_uploaded": true},
	))
}
