package tests

import (
	"context"
	"testing"

	pb "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/proto"
	"google.golang.org/grpc"
)

/**
 * Test send data to server grpc
 */
func TestServerGrpc(t *testing.T) {
	// Connect to the gRPC server
	conn, err := grpc.Dial(
		"127.0.0.1:50051",
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	// Create a new gRPC client
	client := pb.NewDispatcherClient(conn)
	// Send a message
	resp, err := client.SendMessage(context.Background(), &pb.MessageRequest{
		ConnId:  "31f09f0f-5f76-45f0-aa66-e8e607855343",
		Payload: []byte("Hello !"),
	})
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}
	t.Logf("Response from server: %v", resp)
}
