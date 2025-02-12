package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "grpc-train-booking/generated"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTrainBookingServer
	mu    sync.Mutex
	users map[string]*pb.Receipt
}

func (s *server) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seatNumber := fmt.Sprintf("S-%d", len(s.users)+1)
	section := "A"
	if len(s.users)%2 == 0 {
		section = "B"
	}

	receipt := &pb.Receipt{
		From:      "London",
		To:        "France",
		User:      req.FirstName + " " + req.LastName,
		Email:     req.Email,
		Seat:      seatNumber,
		Section:   section,
		PricePaid: 20.0,
	}
	s.users[req.Email] = receipt
	log.Printf("Ticket purchased: %+v", receipt)
	return receipt, nil
}

func (s *server) GetReceipt(ctx context.Context, req *pb.UserRequest) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, exists := s.users[req.Email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
        log.Println("Receipt", receipt)
	return receipt, nil
}

func (s *server) GetUsersBySection(ctx context.Context, req *pb.SectionRequest) (*pb.UserList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []*pb.Receipt
	for _, receipt := range s.users {
		if receipt.Section == req.Section {
			users = append(users, receipt)
		}
	}
        log.Println("Userbysection", &pb.UserList{Users: users})
	return &pb.UserList{Users: users}, nil
}

func (s *server) RemoveUser(ctx context.Context, req *pb.UserRequest) (*pb.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[req.Email]; !exists {
		return nil, fmt.Errorf("user not found")
	}

	delete(s.users, req.Email)
	return &pb.Response{Message: "User removed successfully"}, nil
}

func (s *server) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	receipt, exists := s.users[req.Email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	receipt.Seat = req.NewSeat
	receipt.Section = req.NewSection
        log.Println("Modified seat section", receipt.Seat,receipt.Section)
	return &pb.Response{Message: "Seat modified successfully"}, nil
}

func main() {
	users := make(map[string]*pb.Receipt)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := &server{users: users}

	pb.RegisterTrainBookingServer(grpcServer, s)

	log.Println("gRPC server running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}