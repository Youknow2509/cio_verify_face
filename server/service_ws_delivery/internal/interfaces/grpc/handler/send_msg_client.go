package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/handler"
	pb "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/proto"
)

/**
 * Impl grpc gen dispatcher server
 */
type SendMsgClient struct{}

// SendMessage implements gen.DispatcherServer.
func (s *SendMsgClient) SendMessage(ctx context.Context, data *pb.MessageRequest) (*pb.SendMessageResponse, error) {
	// parse uuid
	connId, err := uuid.Parse(data.ConnId)
	if err != nil {
		return &pb.SendMessageResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}
	if err := handler.GetWorkerHandler().SendDataToClient(
		ctx,
		connId,
		data.Payload,
	); err != nil {
		return &pb.SendMessageResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}
	return &pb.SendMessageResponse{
		Success: true,
		Message: "Message sent successfully",
	}, nil
}

/**
 * New SendMsgClient
 */
func NewSendMsgClient() *SendMsgClient {
	return &SendMsgClient{}
}
