package transaction

type service struct {
	repository Repository
}

type Service interface {
	GetTransactionsByCampaignID(input GetTransactionsByCampaignIdInput) ([]Transaction, error)
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetTransactionsByCampaignID(input GetTransactionsByCampaignIdInput) ([]Transaction, error) {
	transactions, err := s.repository.FindByCampaignID(input.ID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
