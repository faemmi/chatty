package main_test

import (
	pb "chatty/protos/message"
	"chatty/utils"
	"context"
	"log"
	"net"
	"testing"

	server "chatty/server"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)



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
	if res.Success != true {
		t.Errorf("Expected successful response, got %q", res.Error)
	}

	assert.Equal(t, utils.IsValidUUID(res.MessageId), true, "MessageId is not a valid UUID")
}

func TestGetMessages(t *testing.T) {
	ctx := context.Background()
	client, closer := newMessageServerClient()
	defer closer()

	send_request := &pb.SendMessageRequest{
		SenderId:   "test-sender-id",
		ReceiverId: "test-receiver-id",
		Content:    "test-content",
	}

	// Call the RPC
	send_response, err := client.SendMessage(ctx, send_request)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	// Assert the response
	if send_response.Success != true {
		t.Errorf("Expected successful response, got %q", send_response.Error)
	}

	get_request := &pb.GetMessagesRequest{
		UserId: "test-receiver-id",
	}

	stream, err := client.GetMessages(ctx, get_request)
	if err != nil {
		t.Fatalf("GetMessages stream failed: %v", err)
	}
	get_response, err := stream.Recv()
	if err != nil {
		t.Fatalf("Receiving a message from the stream failed: %v", err)
	}

	assert.Equal(t, get_response.Id, send_response.MessageId, "Received message ID does not equal to message ID in sent message")
	assert.Equal(t, get_response.SenderId, send_request.SenderId, "Received sender ID neq to sender ID in sent message")
	assert.Equal(t, get_response.ReceiverId, send_request.ReceiverId, "Received receiver ID neq to receiver ID in send message")
	assert.Equal(t, get_response.Content, send_request.Content, "Received content neq to receiver ID in sent message")

}

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
