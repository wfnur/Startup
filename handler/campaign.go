package handler

import (
	"net/http"
	"startup/campaign"
	"startup/helper"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaign, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List Of Campaign", http.StatusOK, "Error", campaign)
	c.JSON(http.StatusOK, response)

}
