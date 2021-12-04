package handler

import (
	"funding/campaign"
	"funding/helper"
	"net/http"
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
