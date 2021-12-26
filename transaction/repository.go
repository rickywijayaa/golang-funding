package transaction

import "gorm.io/gorm"

type Repository interface {
	FindByCampaignID(campaignID int) ([]Transaction, error)
	FindByUserID(userID int) ([]Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindByCampaignID(campaignID int) ([]Transaction, error) {
	var transaction []Transaction

	err := r.db.Preload("User").
		Where("campaign_id = ?", campaignID).
		Order("created_at desc").
		Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) FindByUserID(userID int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("Campaign.CampaignImages", "is_primary = 1").
		Where("user_id = ?", userID).
		Order("id desc").
		Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
