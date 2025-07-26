package main

import (
	exchangepb "exchange/api/v1"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ordersservice "exchange/services/orders"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	orders, err := ordersservice.New()
	if err != nil {
		log.Fatalf("Failed to create orders service: %v", err)
	}
	exchangepb.RegisterOrdersServiceServer(s, orders)

	// enable reflection
	reflection.Register(s)

	log.Println("gRPC server is listening on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
