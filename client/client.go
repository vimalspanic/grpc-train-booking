package main

import (
	"context"
	"log"
	"time"

	pb "grpc-train-booking/generated"

	"google.golang.org/grpc"
)

func main() {
        conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
        if err != nil {
                log.Fatalf("Failed to connect: %v", err)
        }
        defer conn.Close()

        client := pb.NewTrainBookingClient(conn)

        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        res, err := client.PurchaseTicket(ctx, &pb.PurchaseRequest{
                FirstName: "Vimal",
                LastName:  "Srinivasan",
                Email:     "vimal.srini@example.com",
        })
        log.Printf("Ticket Purchased: %+v", res)
        if err != nil {
                log.Fatalf("Error purchasing ticket: %v", err)
        }

        // Get Receipt
        receipt, err := client.GetReceipt(ctx, &pb.UserRequest{Email: "vimal.srini@example.com"})
        if err != nil {
                log.Fatalf("Error fetching receipt: %v", err)
        }
        log.Printf("Receipt Details: %+v\n", receipt)

        Userbysection , err :=client.GetUsersBySection(ctx, &pb.SectionRequest{Section: receipt.Section})
        if err != nil{
                log.Fatalf("Error fetching user bysection :%v ", err)
        }
        log.Printf("\nUser Details :%v", Userbysection)
}
