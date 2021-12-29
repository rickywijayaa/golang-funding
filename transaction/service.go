package transaction

import (
	"errors"
	"fmt"
	"funding/campaign"
	"funding/payment"
	"strconv"
	"strings"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

type Service interface {
	GetTransactionsByCampaignID(input GetTransactionsByCampaignIdInput) ([]Transaction, error)
	GetTransactionsByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetTransactionsByCampaignID(input GetTransactionsByCampaignIdInput) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("cannot see transaction that not own")
	}

	transactions, err := s.repository.FindByCampaignID(input.ID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionsByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.FindByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{
		Amount:     input.Amount,
		CampaignID: input.CampaignID,
		UserID:     input.User.ID,
		Status:     "PENDING",
	}

	isAnyTransactions, err := s.repository.FindOneTransaction()
	if err != nil {
		return Transaction{}, err
	}

	if isAnyTransactions.ID != 0 {
		lastOrderID, err := s.repository.FindLastOrderID()

		if err != nil {
			return lastOrderID, err
		}

		if lastOrderID.ID != 0 && lastOrderID.Code != "" {
			lastOrderNumber := strings.Split(lastOrderID.Code, "-")
			resultOrderNumberToInt, _ := strconv.Atoi(lastOrderNumber[1])
			transaction.Code = fmt.Sprintf("ORDER-%v", resultOrderNumberToInt+1)
		}
	} else {
		transaction.Code = "ORDER-1"
	}

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
		Code:   newTransaction.Code,
	}

	paymentURL, err := s.paymentService.GetPaymentUrl(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}
