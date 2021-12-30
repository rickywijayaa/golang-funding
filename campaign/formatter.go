package campaign

import "strings"

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

type CampaignDetailFormatter struct {
	ID               int                      `json:"id"`
	Name             string                   `json:"name"`
	ShortDescription string                   `json:"short_description"`
	Description      string                   `json:"description"`
	ImageURL         string                   `json:"image_url"`
	GoalAmount       int                      `json:"goal_amount"`
	CurrentAmount    int                      `json:"current_amount"`
	BackerCount      int                      `json:"backer_count"`
	UserID           int                      `json:"user_id"`
	Slug             string                   `json:"slug"`
	Perks            []string                 `json:"perks"`
	User             CampaignUserFormatter    `json:"user"`
	Images           []CampaignImageFormatter `json:"images"`
}

type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type CampaignImageFormatter struct {
	Filename  string `json:"image_url"`
	IsPrimary int    `json:"is_primary"`
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

func CampaignDetailFormat(campaign Campaign) CampaignDetailFormatter {
	var perks []string
	var perksList = strings.Split(campaign.Perks, ",")

	for _, perk := range perksList {
		perks = append(perks, strings.TrimSpace(perk))
	}

	user := campaign.User
	campaignUserFormatter := CampaignUserFormatter{
		Name:     user.Name,
		ImageUrl: user.AvatarFileName,
	}

	var images []CampaignImageFormatter

	campaignImages := campaign.CampaignImages

	for _, img := range campaignImages {
		imageFormatter := CampaignImageFormatter{
			Filename:  img.FileName,
			IsPrimary: img.IsPrimary,
		}

		images = append(images, imageFormatter)
	}

	formatter := CampaignDetailFormatter{
		ID:               campaign.ID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		ImageURL:         "",
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		BackerCount:      campaign.BackerCount,
		Slug:             campaign.Slug,
		Perks:            perks,
		User:             campaignUserFormatter,
		Images:           images,
	}

	if len(campaign.CampaignImages) > 0 {
		if campaign.CampaignImages[0].FileName != "" {
			formatter.ImageURL = campaign.CampaignImages[0].FileName
		}
	}

	return formatter
}
