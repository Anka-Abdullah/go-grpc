package main

import (
	"go-grpc-crud/app/handler"
	"go-grpc-crud/config/database"
	"go-grpc-crud/proto/book"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	database.Connect()
	bookHandler := handler.NewBookHandler(database.DB)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	book.RegisterBookServiceServer(grpcServer, bookHandler)

	log.Println("gRPC server is running on port :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
