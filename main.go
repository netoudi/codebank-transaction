package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/netoudi/codebank-transaction/infrastructure/grpc/server"
	"github.com/netoudi/codebank-transaction/infrastructure/kafka"
	"github.com/netoudi/codebank-transaction/infrastructure/repository"
	"github.com/netoudi/codebank-transaction/usecase"
	"log"
)

func main() {
	db := setupDb()
	defer db.Close()
	producer := setupKafkaProducer()
	transactionUseCase := setupTransactionUseCase(db, producer)
	serveGrpc(transactionUseCase)
}

func setupTransactionUseCase(db *sql.DB, producer kafka.KafkaProducer) usecase.TransactionUseCase {
	return usecase.NewTransactionUseCase(repository.NewTransactionRepositoryDb(db), producer)
}

func setupKafkaProducer() kafka.KafkaProducer {
	p := kafka.NewKafkaProducer()
	p.SetupProducer("host.docker.internal:9094")
	return p
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

func serveGrpc(transactionUseCase usecase.TransactionUseCase) {
	grpcServer := server.NewGrpcServer()
	grpcServer.TransactionUseCase = transactionUseCase
	fmt.Println("🚀 running gRPC server...")
	grpcServer.Serve()
}
