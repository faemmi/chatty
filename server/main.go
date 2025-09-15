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

type messageServer struct {
	pb.UnimplementedMessagesServer
}

func NewMessageServer() (*messageServer) {
	s := &messageServer{}
	return s
}

func (s *messageServer) SendMessage(_ context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	log.Printf("Received: %v", protojson.Format(request))
	return &pb.SendMessageResponse{Success: true, MessageId: "test-id", Error: ""}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":51001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterMessagesServer(server, &messageServer{})
	reflection.Register(server)
	log.Printf("Server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
