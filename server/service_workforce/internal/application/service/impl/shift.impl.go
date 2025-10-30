package service

import (
	// "context"
	"context"
	"encoding/json"
	"fmt"

	applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
	service "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/cache"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/logger"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/repository"
)

const (
	// Cache key prefixes
	shiftCachePrefix = "shift:"
	shiftCacheTTL    = 3600 // 1 hour in seconds
)

// =================================================
// Shift service implementation interface
// =================================================
type ShiftService struct {
	shiftRepo        repository.IShiftRepository
	logger           logger.ILogger
	distributedCache cache.IDistributedCache
	localCache       cache.ILocalCache
}

// CreateShift implements service.IShiftService.
func (s *ShiftService) CreateShift(ctx context.Context, input *applicationModel.CreateShiftInput) (*applicationModel.CreateShiftOutput, *applicationError.Error) {
	// Validate input
	if input == nil {
		s.logger.Error("CreateShift - Input is nil")
		return nil, &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("CreateShift - Start", "user_id", input.UserId, "company_id", input.CompanyId)

	// Convert work days from []int to []int32
	workDays := make([]int32, len(input.WorkDays))
	for i, day := range input.WorkDays {
		workDays[i] = int32(day)
	}

	// Create domain input
	domainInput := &domainModel.CreateShiftInput{
		CompanyID:             input.CompanyId,
		Name:                  input.Name,
		Description:           input.Description,
		StartTime:             input.StartTime,
		EndTime:               input.EndTime,
		BreakDurationMinutes:  int32(input.BreakDurationMinutes),
		GracePeriodMinutes:    int32(input.GracePeriodMinutes),
		EarlyDepartureMinutes: int32(input.EarlyDepartureMinutes),
		WorkDays:              workDays,
	}

	// Call repository
	shiftID, err := s.shiftRepo.CreateShift(ctx, domainInput)
	if err != nil {
		s.logger.Error("CreateShift - Failed to create shift", "error", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to create shift",
		}
	}

	s.logger.Info("CreateShift - Success", "shift_id", shiftID.String())

	return &applicationModel.CreateShiftOutput{
		ShiftId: shiftID.String(),
	}, nil
}

// DeleteShift implements service.IShiftService.
func (s *ShiftService) DeleteShift(ctx context.Context, input *applicationModel.DeleteShiftInput) *applicationError.Error {
	// Validate input
	if input == nil {
		s.logger.Error("DeleteShift - Input is nil")
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("DeleteShift - Start", "user_id", input.UserId, "shift_id", input.ShiftId)

	// Call repository
	err := s.shiftRepo.DeleteShift(ctx, input.ShiftId)
	if err != nil {
		s.logger.Error("DeleteShift - Failed to delete shift", "error", err, "shift_id", input.ShiftId)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to delete shift",
		}
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("%s%s", shiftCachePrefix, input.ShiftId.String())

	// Delete from distributed cache
	if delErr := s.distributedCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("DeleteShift - Failed to delete from distributed cache", "error", delErr)
	}

	// Delete from local cache
	if delErr := s.localCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("DeleteShift - Failed to delete from local cache", "error", delErr)
	}

	s.logger.Info("DeleteShift - Success", "shift_id", input.ShiftId)

	return nil
}

// EditShift implements service.IShiftService.
func (s *ShiftService) EditShift(ctx context.Context, input *applicationModel.EditShiftInput) *applicationError.Error {
	// Validate input
	if input == nil {
		s.logger.Error("EditShift - Input is nil")
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("EditShift - Start", "user_id", input.UserId, "shift_id", input.ShiftId)

	// Convert work days from []int to []int32
	workDays := make([]int32, len(input.WorkDays))
	for i, day := range input.WorkDays {
		workDays[i] = int32(day)
	}

	// Create domain input for update
	domainInput := &domainModel.UpdateTimeShiftInput{
		ShiftID:               input.ShiftId,
		StartTime:             input.StartTime,
		EndTime:               input.EndTime,
		BreakDurationMinutes:  int32(input.BreakDurationMinutes),
		GracePeriodMinutes:    int32(input.GracePeriodMinutes),
		EarlyDepartureMinutes: int32(input.EarlyDepartureMinutes),
		WorkDays:              workDays,
	}

	// Call repository
	err := s.shiftRepo.UpdateTimeShift(ctx, domainInput)
	if err != nil {
		s.logger.Error("EditShift - Failed to update shift", "error", err, "shift_id", input.ShiftId)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to update shift",
		}
	}

	// Invalidate cache after update
	cacheKey := fmt.Sprintf("%s%s", shiftCachePrefix, input.ShiftId.String())

	// Delete from distributed cache
	if delErr := s.distributedCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("EditShift - Failed to delete from distributed cache", "error", delErr)
	}

	// Delete from local cache
	if delErr := s.localCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("EditShift - Failed to delete from local cache", "error", delErr)
	}

	s.logger.Info("EditShift - Success", "shift_id", input.ShiftId)

	return nil
}

// GetDetailShift implements service.IShiftService.
func (s *ShiftService) GetDetailShift(ctx context.Context, input *applicationModel.GetDetailShiftInput) (*applicationModel.GetDetailShiftOutput, *applicationError.Error) {
	// Validate input
	if input == nil {
		s.logger.Error("GetDetailShift - Input is nil")
		return nil, &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("GetDetailShift - Start", "user_id", input.UserId, "shift_id", input.ShiftId)

	cacheKey := fmt.Sprintf("%s%s", shiftCachePrefix, input.ShiftId.String())

	// Try to get from local cache first
	if cachedData, err := s.localCache.Get(ctx, cacheKey); err == nil && cachedData != "" {
		s.logger.Info("GetDetailShift - Cache hit (local)", "shift_id", input.ShiftId)
		var output applicationModel.GetDetailShiftOutput
		if unmarshalErr := json.Unmarshal([]byte(cachedData), &output); unmarshalErr == nil {
			return &output, nil
		} else {
			s.logger.Warn("GetDetailShift - Failed to unmarshal local cache", "error", unmarshalErr)
		}
	}

	// Try to get from distributed cache
	if cachedData, err := s.distributedCache.Get(ctx, cacheKey); err == nil && cachedData != "" {
		s.logger.Info("GetDetailShift - Cache hit (distributed)", "shift_id", input.ShiftId)
		var output applicationModel.GetDetailShiftOutput
		if unmarshalErr := json.Unmarshal([]byte(cachedData), &output); unmarshalErr == nil {
			// Store in local cache for faster access next time
			if jsonData, _ := json.Marshal(output); len(jsonData) > 0 {
				_ = s.localCache.SetTTL(ctx, cacheKey, string(jsonData), shiftCacheTTL)
			}
			return &output, nil
		} else {
			s.logger.Warn("GetDetailShift - Failed to unmarshal distributed cache", "error", unmarshalErr)
		}
	}

	s.logger.Info("GetDetailShift - Cache miss, fetching from database", "shift_id", input.ShiftId)

	// Call repository
	shift, err := s.shiftRepo.GetShiftByID(ctx, input.ShiftId)
	if err != nil {
		s.logger.Error("GetDetailShift - Failed to get shift details", "error", err, "shift_id", input.ShiftId)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to get shift details",
		}
	}

	// Convert work days from []int32 to []int
	workDays := make([]int, len(shift.WorkDays))
	for i, day := range shift.WorkDays {
		workDays[i] = int(day)
	}

	output := &applicationModel.GetDetailShiftOutput{
		ShiftId:               shift.ShiftID.String(),
		CompanyId:             shift.CompanyID.String(),
		Name:                  shift.Name,
		Description:           shift.Description,
		StartTime:             shift.StartTime,
		EndTime:               shift.EndTime,
		BreakDurationMinutes:  int(shift.BreakDurationMinutes),
		GracePeriodMinutes:    int(shift.GracePeriodMinutes),
		EarlyDepartureMinutes: int(shift.EarlyDepartureMinutes),
		WorkDays:              workDays,
		IsActive:              shift.IsActive,
	}

	// Cache the result
	if jsonData, jsonErr := json.Marshal(output); jsonErr == nil {
		// Store in distributed cache
		if setErr := s.distributedCache.SetTTL(ctx, cacheKey, string(jsonData), shiftCacheTTL); setErr != nil {
			s.logger.Warn("GetDetailShift - Failed to set distributed cache", "error", setErr)
		}
		// Store in local cache
		if setErr := s.localCache.SetTTL(ctx, cacheKey, string(jsonData), shiftCacheTTL); setErr != nil {
			s.logger.Warn("GetDetailShift - Failed to set local cache", "error", setErr)
		}
	}

	s.logger.Info("GetDetailShift - Success", "shift_id", input.ShiftId)

	return output, nil
}

// New instance
func NewShiftService() service.IShiftService {
	shiftRepo, err := repository.GetShiftRepository()
	if err != nil {
		panic(fmt.Sprintf("Failed to get shift repository: %v", err))
	}

	log := logger.GetLogger()
	if log == nil {
		panic("Failed to get logger instance")
	}

	distributedCache, err := cache.GetDistributedCache()
	if err != nil {
		panic(fmt.Sprintf("Failed to get distributed cache: %v", err))
	}

	localCache, err := cache.GetLocalCache()
	if err != nil {
		panic(fmt.Sprintf("Failed to get local cache: %v", err))
	}

	return &ShiftService{
		shiftRepo:        shiftRepo,
		logger:           log,
		distributedCache: distributedCache,
		localCache:       localCache,
	}
}
