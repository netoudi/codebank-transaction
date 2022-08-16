package server

import (
	"github.com/netoudi/codebank-transaction/infrastructure/grpc/pb"
	"github.com/netoudi/codebank-transaction/infrastructure/grpc/service"
	"github.com/netoudi/codebank-transaction/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type GrpcServer struct {
	TransactionUseCase usecase.TransactionUseCase
}

func NewGrpcServer() GrpcServer {
	return GrpcServer{}
}

func (g GrpcServer) Serve() {
	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("could not listen tcp port")
	}
	transactionService := service.NewTransactionService()
	transactionService.TransactionUseCase = g.TransactionUseCase
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterPaymentServiceServer(grpcServer, transactionService)
	grpcServer.Serve(listen)
}
