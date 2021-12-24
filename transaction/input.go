package transaction

import "funding/user"

type GetTransactionsByCampaignIdInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
