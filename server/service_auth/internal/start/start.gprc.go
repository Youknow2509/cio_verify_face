package start

// import (
// 	"fmt"
// 	"net"
// 	"os"

// 	global "github.com/youknow2509/cio_verify_face/server/service_auth/internal/global"
// 	routes "github.com/youknow2509/cio_verify_face/server/service_auth/internal/interfaces/grpc/routes"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials"
// )

// // init server grpc
// func initServerGrpc() error {
// 	config := global.SettingServer.GrpcServer
// 	// Initialize the gRPC server
// 	lis, err := net.Listen(
// 		"tcp", fmt.Sprintf("localhost:%d", config.Port),
// 	)
// 	if err != nil {
// 		global.Logger.Error("failed to listen: %v", err)
// 		return err
// 	}
// 	// ServerOption
// 	var opts []grpc.ServerOption
// 	// TLS
// 	if config.Tls.Enabled {
// 		// check file existence
// 		if _, err := os.Stat(config.Tls.CertFile); os.IsNotExist(err) {
// 			global.Logger.Error("TLS cert file does not exist: %v", err)
// 			return err
// 		}
// 		if _, err := os.Stat(config.Tls.KeyFile); os.IsNotExist(err) {
// 			global.Logger.Error("TLS key file does not exist: %v", err)
// 			return err
// 		}
// 		// create TLS credentials
// 		creds, err := credentials.NewServerTLSFromFile(
// 			config.Tls.CertFile,
// 			config.Tls.KeyFile,
// 		)
// 		if err != nil {
// 			global.Logger.Error("failed to generate credentials: %v", err)
// 			return err
// 		}
// 		opts = []grpc.ServerOption{grpc.Creds(creds)}
// 	}
// 	// init server
// 	grpcServer := grpc.NewServer(opts...)
// 	routes.InitGrpcRoutes(grpcServer)
// 	// start server
// 	fmt.Printf("gRPC server listening on %s\n", lis.Addr().String())
// 	go func() {
// 		if err := grpcServer.Serve(lis); err != nil {
// 			global.Logger.Error("failed to start gRPC server: %v", err)
// 			return
// 		}
// 	}()
// 	return nil
// }
