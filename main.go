package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/netoudi/codebank-transaction/domain"
	"github.com/netoudi/codebank-transaction/infrastructure/repository"
	"github.com/netoudi/codebank-transaction/usecase"
	"log"
)

func main() {
	fmt.Println("Hello world!")
	db := setupDb()
	defer db.Close()

	card := domain.NewCreditCard()
	card.Name = "John Doe"
	card.Number = "123456789"
	card.ExpirationMonth = 12
	card.ExpirationYear = 2022
	card.Cvv = 123
	card.Balance = 1000
	card.Limit = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*card)
	if err != nil {
		fmt.Println(err)
	}
}

func setupTransactionUseCase(db *sql.DB) usecase.TransactionUseCase {
	return usecase.NewTransactionUseCase(repository.NewTransactionRepositoryDb(db))
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db",
		"5432",
		"root",
		"root",
		"codebank",
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connection to database")
	}
	return db
}
