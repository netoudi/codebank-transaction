package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/netoudi/codebank-transaction/infrastructure/grpc/server"
	"github.com/netoudi/codebank-transaction/infrastructure/kafka"
	"github.com/netoudi/codebank-transaction/infrastructure/repository"
	"github.com/netoudi/codebank-transaction/usecase"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

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
	p.SetupProducer(os.Getenv("KAFKA_BOOTSTRAP_SERVERS"))
	return p
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
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
