package transaction

import (
	"time"
)

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTransactionFormatter struct {
	ID        int               `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func CampaignTransactionFormat(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}

	return formatter
}

func CampaignTransactionsFormat(transaction []Transaction) []CampaignTransactionFormatter {
	if len(transaction) == 0 {
		return []CampaignTransactionFormatter{}
	}

	var transactionFormatter []CampaignTransactionFormatter

	for _, trans := range transaction {
		formatter := CampaignTransactionFormat(trans)
		transactionFormatter = append(transactionFormatter, formatter)
	}

	return transactionFormatter
}

func UserTransactionFormat(transaction Transaction) UserTransactionFormatter {
	campaignFormatter := CampaignFormatter{
		Name:     transaction.Campaign.Name,
		ImageURL: "",
	}

	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	formatter := UserTransactionFormatter{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
		Campaign:  campaignFormatter,
	}

	return formatter
}

func UserTransactionsFormat(transaction []Transaction) []UserTransactionFormatter {
	if len(transaction) == 0 {
		return []UserTransactionFormatter{}
	}

	var transactionFormatter []UserTransactionFormatter

	for _, trans := range transaction {
		formatter := UserTransactionFormat(trans)
		transactionFormatter = append(transactionFormatter, formatter)
	}

	return transactionFormatter
}
