package repository

import (
	"context"
	"errors"
	"os"

	"github.com/redis/go-redis/v9"
	constants "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/repository"
)

// ======================================================================================================
// Redis cache manager connection repository implementation
// ======================================================================================================
type RedisManagerConnectionRepository struct {
	client *redis.Client
}

// CreateConnection implements repository.IManagerConnectionRepository.
func (r *RedisManagerConnectionRepository) CreateConnection(ctx context.Context, input *model.CreateConnectionInput) (bool, error) {
	// Load Lua script from file
	scriptBytes, err := os.ReadFile(constants.LUA_SCRIPT_CREATE_CONNECTION_PATH)
	if err != nil {
		return false, errors.New("failed to read Lua script file: " + err.Error())
	}
	script := string(scriptBytes)
	// Create key and arguments for the Lua script
	keys := []string{
		input.UserConnectionsKey,
		input.ServiceConnectionsKey,
		input.ConnectionKey,
	}
	args := []interface{}{
		input.ConnectionId,
		input.UserId,
		input.ServiceId,
		input.IPAddress,
		input.ConnectedAt,
		input.UserAgent,
		input.MaxConnsPerUser,
	}
	// Send script to Redis
	res, err := r.client.Eval(ctx, script, keys, args...).Result()
	if err != nil {
		return false, errors.New("failed to execute Lua script: " + err.Error())
	}
	intRes, ok := res.(int64)
	if !ok {
		return false, errors.New("unexpected result type from Redis Eval")
	}
	return intRes == 1, nil
}

// RemoveConnection implements repository.IManagerConnectionRepository.
func (r *RedisManagerConnectionRepository) RemoveConnection(ctx context.Context, input *model.RemoveConnectionInput) (bool, error) {
	// Load Lua script from file
	scriptBytes, err := os.ReadFile(constants.LUA_SCRIPT_REMOVE_CONNECTION_PATH)
	if err != nil {
		return false, errors.New("failed to read Lua script file: " + err.Error())
	}
	script := string(scriptBytes)
	// Create key and arguments for the Lua script
	keys := []string{
		input.UserConnectionsKey,
		input.ServiceConnectionsKey,
		input.ConnectionKey,
	}
	args := []interface{}{
		input.ConnectionId,
	}
	// Send script to Redis
	res, err := r.client.Eval(ctx, script, keys, args...).Result()
	if err != nil {
		return false, errors.New("failed to execute Lua script: " + err.Error())
	}
	intRes, ok := res.(int64)
	if !ok {
		return false, errors.New("unexpected result type from Redis Eval")
	}
	return intRes == 1, nil
}

/**
 * NewRedisManagerConnectionRepository creates a new instance of RedisManagerConnectionRepository
 * implementation Domain ManagerConnectionRepository
 */
func NewRedisManagerConnectionRepository(client *redis.Client) domainRepository.IManagerConnectionRepository {
	return &RedisManagerConnectionRepository{
		client: client,
	}
}

// ======================================================================================================
