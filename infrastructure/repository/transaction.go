package repository

import (
	"database/sql"
	"errors"
	"github.com/netoudi/codebank-transaction/domain"
)

type TransactionRepositoryDb struct {
	db *sql.DB
}

func NewTransactionRepositoryDb(db *sql.DB) *TransactionRepositoryDb {
	return &TransactionRepositoryDb{db: db}
}

func (t *TransactionRepositoryDb) SaveTransaction(transaction domain.Transaction, card domain.CreditCard) error {
	stmt, err := t.db.Prepare(`insert into transactions (id, credit_card_id, amount, status, description, store, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		transaction.Id,
		transaction.CreditCardId,
		transaction.Amount,
		transaction.Status,
		transaction.Description,
		transaction.Store,
		transaction.CreatedAt,
	)
	if err != nil {
		return err
	}
	if transaction.Status == "approved" {
		err = t.updateBalance(card)
		if err != nil {
			return err
		}
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepositoryDb) GetCreditCard(creditCard domain.CreditCard) (domain.CreditCard, error) {
	var c domain.CreditCard
	stmt, err := t.db.Prepare("select id, balance, balance_limit from credit_cards where number=$1")
	if err != nil {
		return c, err
	}
	if err = stmt.QueryRow(creditCard.Number).Scan(&c.Id, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("credit card does not exists")
	}
	return c, nil
}

func (t *TransactionRepositoryDb) CreateCreditCard(card domain.CreditCard) error {
	stmt, err := t.db.Prepare(`insert into credit_cards (id, name, number, expiration_month, expiration_year, cvv, balance, balance_limit, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		card.Id,
		card.Name,
		card.Number,
		card.ExpirationMonth,
		card.ExpirationYear,
		card.Cvv,
		card.Balance,
		card.Limit,
		card.CreatedAt,
	)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionRepositoryDb) updateBalance(card domain.CreditCard) error {
	_, err := t.db.Exec("update credit_cards set balance = $1 where id = $2", card.Balance, card.Id)
	if err != nil {
		return err
	}
	return nil
}
