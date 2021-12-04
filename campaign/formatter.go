package campaign

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

func CampaignFormat(campaign Campaign) CampaignFormatter {
	formatter := CampaignFormatter{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		ImageURL:         "",
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
	}

	if len(campaign.CampaignImages) > 0 {
		if campaign.CampaignImages[0].FileName != "" {
			formatter.ImageURL = campaign.CampaignImages[0].FileName
		}
	}

	return formatter
}

func CampaignsFormat(campaigns []Campaign) []CampaignFormatter {
	formatter := []CampaignFormatter{}

	for _, campaign := range campaigns {
		campaignFormat := CampaignFormat(campaign)
		formatter = append(formatter, campaignFormat)
	}

	return formatter
}
