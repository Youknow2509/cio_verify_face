package tests

// import (
// 	"context"
// 	"testing"

// 	libsClients "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/conn"
// 	libsConfig "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/config"
// 	domainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
// 	domainRepository "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/repository"
// 	infrastructureRepository "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/repository"
// )

// // Create interface for manager connection repository
// func createManagerConnectionRepository(t *testing.T) domainRepository.IManagerConnectionRepository {
// 	// Initialize Redis client
// 	redisSetting := &libsConfig.RedisSetting{
// 		Type:           1,
// 		UseTLS:         false,
// 		CertPath:       "",
// 		KeyPath:        "",
// 		Password:       "",
// 		DB:             0,
// 		Host:           "127.0.0.1",
// 		Port:           6379,
// 		MasterName:     "mymaster",
// 		SentinelAddrs:  []string{"127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381"},
// 		Address:        []string{"127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381"},
// 		RouteByLatency: false,
// 		RouteRandomly:  false,
// 		PoolSize:       10,
// 		MinIdleConns:   2,
// 		MaxRetries:     3,
// 	}
// 	err := libsClients.InitRedisClient(redisSetting)
// 	if err != nil {
// 		t.Fatalf("Failed to initialize Redis client: %v", err)
// 		return nil
// 	}
// 	redisClient, err := libsClients.GetRedisClient()
// 	if err != nil {
// 		t.Fatalf("Failed to get Redis client: %v", err)
// 		return nil
// 	}
// 	// Initialize infrastructure repository impl
// 	infrastructureRepo := infrastructureRepository.NewRedisManagerConnectionRepository(redisClient)
// 	err = domainRepository.SetManagerConnectionRepository(infrastructureRepo)
// 	if err != nil {
// 		t.Fatalf("Failed to set manager connection repository: %v", err)
// 		return nil
// 	}
// 	managerConnectionRepository := domainRepository.GetManagerConnectionRepository()
// 	return managerConnectionRepository
// }

// // Test create connection Lua script
// func TestCreateConnectionWithLua(t *testing.T) {
// 	// Get interface manager
// 	managerConnectionRepository := createManagerConnectionRepository(t)
// 	if managerConnectionRepository == nil {
// 		t.Fatal("Manager connection repository is nil")
// 		return
// 	}
// 	// Handle Lua script for creating connection
// 	ctx := context.Background()
// 	UserId := "123"
// 	ServiceId := "sv1"
// 	ConnectionId := "conn7"
// 	ok, err := managerConnectionRepository.CreateConnection(
// 		ctx,
// 		&domainModel.CreateConnectionInput{
// 			// Keys for Lua script
// 			UserConnectionsKey:    "user:connections:" + UserId,
// 			ServiceConnectionsKey: "service:connections:" + ServiceId,
// 			ConnectionKey:         "connection:" + ConnectionId,
// 			// Args for Lua script
// 			ConnectionId:    ConnectionId,
// 			UserId:          UserId,
// 			ServiceId:       ServiceId,
// 			IPAddress:       "127.0.0.1",
// 			ConnectedAt:     "1633072800", // Example timestamp
// 			UserAgent:       "Mozilla/5.0",
// 			MaxConnsPerUser: 5,
// 		},
// 	)
// 	if err != nil {
// 		t.Fatalf("Failed to create connection: %v", err)
// 		return
// 	}
// 	if !ok {
// 		t.Fatal("Failed to create connection")
// 		return
// 	}
// }
