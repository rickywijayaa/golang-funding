package transaction

type GetTransactionsByCampaignIdInput struct {
	ID int `uri:"id" binding:"required"`
}
