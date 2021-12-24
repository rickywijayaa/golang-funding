package campaign

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(UserID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(id GetCampaignDetailInput, input CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(UserID int) ([]Campaign, error) {
	if UserID != 0 {
		campaigns, err := s.repository.FindByUserID(UserID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)

	if campaign.ID == 0 {
		return campaign, errors.New("no campaign found under this ID")
	}

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	slugDefination := fmt.Sprintf("%s %d", input.Name, input.User.ID)

	campaign := Campaign{
		Name:             input.Name,
		Description:      input.Description,
		ShortDescription: input.ShortDescription,
		Perks:            input.Perks,
		GoalAmount:       input.GoalAmount,
		UserID:           input.User.ID,
		Slug:             slug.Make(slugDefination),
	}

	existingCampaign, _ := s.repository.FindByUserID(input.User.ID)
	for _, exist := range existingCampaign {
		if strings.TrimSpace(exist.Name) == strings.TrimSpace(input.Name) {
			return Campaign{}, errors.New("cant create a same campaign name")
		}
	}

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *service) UpdateCampaign(id GetCampaignDetailInput, input CreateCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(id.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != input.User.ID {
		return campaign, errors.New("failed to update campaign , due not owning this campaign")
	}

	slugDefination := fmt.Sprintf("%s %d", input.Name, input.User.ID)

	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.Slug = slug.Make(slugDefination)

	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	isPrimary := 0
	if *input.IsPrimary {
		isPrimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, errors.New("failed update image")
		}
	}

	campaignImage := CampaignImage{
		CampaignID: input.CampaignID,
		IsPrimary:  isPrimary,
		FileName:   fileLocation,
	}

	newCampaignImage, err := s.repository.CreateImage(campaignImage)

	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}
