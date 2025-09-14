package main

import (
	pb "chatty/protos/message"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type server struct {
	pb.UnimplementedMessagesServer
}

func (s *server) SendMessage(_ context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	log.Printf("Received: %v", protojson.Format(request))
	return &pb.SendMessageResponse{Success: true, MessageId: "test-id", Error: ""}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":51001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMessagesServer(grpcServer, &server{})
	reflection.Register(grpcServer)
	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
