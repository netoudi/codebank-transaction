package usecase

import (
	"encoding/json"
	"github.com/netoudi/codebank-transaction/domain"
	"github.com/netoudi/codebank-transaction/dto"
	"github.com/netoudi/codebank-transaction/infrastructure/kafka"
	"os"
)

type TransactionUseCase struct {
	TransactionRepository domain.TransactionRepository
	KafkaProducer         kafka.KafkaProducer
}

func NewTransactionUseCase(repository domain.TransactionRepository, producer kafka.KafkaProducer) TransactionUseCase {
	return TransactionUseCase{TransactionRepository: repository, KafkaProducer: producer}
}

func (u *TransactionUseCase) ProcessTransaction(input dto.Transaction) (domain.Transaction, error) {
	card := u.hydrateCreditCard(input)
	cardBalanceAndLimit, err := u.TransactionRepository.GetCreditCard(*card)
	if err != nil {
		return domain.Transaction{}, err
	}
	card.Id = cardBalanceAndLimit.Id
	card.Limit = cardBalanceAndLimit.Limit
	card.Balance = cardBalanceAndLimit.Balance
	t := u.newTransaction(input, cardBalanceAndLimit)
	t.ProcessAndValidate(card)
	err = u.TransactionRepository.SaveTransaction(*t, *card)
	if err != nil {
		return domain.Transaction{}, err
	}
	input.Id = t.Id
	input.CreatedAt = t.CreatedAt
	transactionJson, err := json.Marshal(input)
	if err != nil {
		return domain.Transaction{}, err
	}
	err = u.KafkaProducer.Publish(string(transactionJson), os.Getenv("KAFKA_TRANSACTIONS_TOPIC"))
	if err != nil {
		return domain.Transaction{}, err
	}
	return *t, nil
}

func (u *TransactionUseCase) hydrateCreditCard(input dto.Transaction) *domain.CreditCard {
	c := domain.NewCreditCard()
	c.Name = input.Name
	c.Number = input.Number
	c.ExpirationMonth = input.ExpirationMonth
	c.ExpirationYear = input.ExpirationYear
	c.Cvv = input.Cvv
	return c
}

func (u *TransactionUseCase) newTransaction(input dto.Transaction, card domain.CreditCard) *domain.Transaction {
	t := domain.NewTransaction()
	t.CreditCardId = card.Id
	t.Amount = input.Amount
	t.Store = input.Store
	t.Description = input.Description
	return t
}
