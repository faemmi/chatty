package main_test

import (
	pb "chatty/protos/message"
	"context"
	"log"
	"net"
	"testing"

	server "chatty/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func newMessageServerClient() (pb.MessagesClient, func()) {
	lis := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	pb.RegisterMessagesServer(grpcServer, server.NewMessageServer())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to server: %v", err)
		}
	}()
	dialServer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(dialServer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}
	closer := func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing client: %v", err)
		}
		grpcServer.Stop()
	}
	client := pb.NewMessagesClient(conn)
	return client, closer

}

func TestSendMessage(t *testing.T) {
	ctx := context.Background()
	client, closer := newMessageServerClient()
	defer closer()

	req := &pb.SendMessageRequest{
		SenderId:   "test-sender-id",
		ReceiverId: "test-receiver-id",
		Content:    "test-content",
	}

	// Call the RPC
	res, err := client.SendMessage(ctx, req)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	// Assert the response
	expected := &pb.SendMessageResponse{
		Success:   true,
		MessageId: "test-message-id",
		Error:     "",
	}
	if res.Success != true {
		t.Errorf("Expected response %q, got %q", expected, res.Error)
	}
}
