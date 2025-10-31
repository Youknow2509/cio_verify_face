package routes

import (
	"context"

	handler "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/grpc/handler"
	pb "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/proto"
	"google.golang.org/grpc"
)

// ==================================================
//
//	Grpc routes
//
// ==================================================
type grpcRoutes struct {
	pb.UnimplementedDispatcherServer
}

// NewGrpcRoutesBase creates a new GrpcRoutes
func NewGrpcRoutesBase() *grpcRoutes {
	return &grpcRoutes{}
}

// InitGrpcRoutes creates a new GrpcRoutes
func InitGrpcRoutes(server *grpc.Server) {
	pb.RegisterDispatcherServer(
		server,
		NewGrpcRoutesBase(),
	)
}

/**
 * Routes send gRPC requests to the dispatcher service
 */
func (gr *grpcRoutes) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.SendMessageResponse, error) {
	return handler.NewSendMsgClient().SendMessage(ctx, req)
}
