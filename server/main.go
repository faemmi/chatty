package main

import (
	database "chatty/database"
	pb "chatty/protos/message"
	utils "chatty/utils"
	"fmt"

	"context"
	"errors"
	"log"
	"net"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type messageServer struct {
	pb.UnimplementedMessagesServer
}

func NewMessageServer() *messageServer {
	s := &messageServer{}
	return s
}

func (s *messageServer) SendMessage(ctx context.Context, request *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	log.Printf("Received: %v", protojson.Format(request))

	config := utils.ReadConfig()
	_, messages, disconnect := database.Connect(config)
	defer disconnect()

	id, err := saveMessage(ctx, messages, request)

	if err != nil {
		return &pb.SendMessageResponse{Success: false, MessageId: "", Error: fmt.Sprintf("%v", err)}, nil
	}


	return &pb.SendMessageResponse{Success: true, MessageId: id.String(), Error: ""}, nil
}

func saveMessage(ctx context.Context, messages *mongo.Collection, request *pb.SendMessageRequest) (uuid.UUID, error) {
	id, err := uuid.NewV7()

	if err != nil {
		return id, err
	}

	message := bson.D{
		{Key: "_id", Value: id.String()},
		{Key: "sender_id", Value: request.SenderId},
		{Key: "receiver_id", Value: request.ReceiverId},
		{Key: "content", Value: request.Content},
	}

	result, err := messages.InsertOne(ctx, message)

	if err != nil {
		return id, err
	}

	if result.InsertedID != id.String() {
		deleteMessage(ctx, messages, result.InsertedID)
		msg := fmt.Sprintf("Inserted document ID (%s) is not equal to message ID (%s)", result.InsertedID, id)
		return id, errors.New(msg)

	}

	return id, nil
}

func deleteMessage(ctx context.Context, messages *mongo.Collection, id any) error {
	result, err := messages.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		log.Printf("No document found with the specified ID %s", id)
	} else {
		log.Printf("Deleted message with ID %s", id)
	}

	return nil
}

func main() {
	config := utils.ReadConfig()

	lis, err := net.Listen("tcp", ":" + config.Server.Port)
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

