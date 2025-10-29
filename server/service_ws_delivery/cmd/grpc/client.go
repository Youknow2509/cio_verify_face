package main

import (
	"context"
	"fmt"

	pb "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/proto"
	"google.golang.org/grpc"
)

/**
 * Test send data to server grpc
 */
func ServerGrpcSendMsg() {
	// Connect to the gRPC server
	conn, err := grpc.Dial(
		"127.0.0.1:50051",
		grpc.WithInsecure(),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gRPC server: %v", err)
		return
	}
	// Create a new gRPC client
	client := pb.NewDispatcherClient(conn)
	// Send a message
	resp, err := client.SendMessage(context.Background(), &pb.MessageRequest{
		ConnId:  "31f09f0f-5f76-45f0-aa66-e8e607855343",
		Payload: []byte("Hello v!"),
	})
	if err != nil {
		fmt.Printf("Failed to send message: %v", err)
		return
	}
	fmt.Printf("Response from server: %v", resp)
}

func main() {
	ServerGrpcSendMsg()
}
