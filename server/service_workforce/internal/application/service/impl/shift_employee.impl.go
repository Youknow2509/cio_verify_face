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
	// Cache key prefixes for shift employee
	shiftEmployeeCachePrefix = "shift_employee:"
	shiftEmployeeCacheTTL    = 1800 // 30 minutes in seconds
)

// =================================================
// ShiftEmployee service implementation interface
// =================================================
type ShiftEmployeeService struct {
	shiftUserRepo    repository.IShiftUserRepository
	logger           logger.ILogger
	distributedCache cache.IDistributedCache
	localCache       cache.ILocalCache
}

// AddShiftEmployee implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) AddShiftEmployee(ctx context.Context, input *applicationModel.AddShiftEmployeeInput) (**applicationModel.AddShiftEmployeeOutput, *applicationError.Error) {
	// Validate input
	if input == nil {
		s.logger.Error("AddShiftEmployee - Input is nil")
		return nil, &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("AddShiftEmployee - Start", "user_id", input.UserId, "employee_id", input.EmployeeId, "shift_id", input.ShiftId)

	// Check if user already has a shift in the time range
	checkInput := &domainModel.CheckUserExistShiftInput{
		EmployeeID:    input.EmployeeId,
		EffectiveFrom: input.EffectiveFrom,
		EffectiveTo:   input.EffectiveTo,
		Limit:         1,
		Offset:        0,
	}

	exists, err := s.shiftUserRepo.CheckUserExistShift(ctx, checkInput)
	if err != nil {
		s.logger.Error("AddShiftEmployee - Failed to check existing shift", "error", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to check existing shift",
		}
	}

	if exists {
		s.logger.Warn("AddShiftEmployee - Employee already has a shift in this time range", "employee_id", input.EmployeeId)
		return nil, &applicationError.Error{
			ErrorSystem: fmt.Errorf("employee already has a shift in this time range"),
			ErrorClient: "Employee already has a shift in this time range",
		}
	}

	// Create domain input
	domainInput := &domainModel.AddShiftForEmployeeInput{
		EmployeeID:    input.EmployeeId,
		ShiftID:       input.ShiftId,
		EffectiveFrom: input.EffectiveFrom,
		EffectiveTo:   input.EffectiveTo,
	}

	// Call repository
	err = s.shiftUserRepo.AddShiftForEmployee(ctx, domainInput)
	if err != nil {
		s.logger.Error("AddShiftEmployee - Failed to add shift to employee", "error", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to add shift to employee",
		}
	}

	// Invalidate cache for this employee
	cacheKey := fmt.Sprintf("%s%s", shiftEmployeeCachePrefix, input.EmployeeId.String())
	if delErr := s.distributedCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("AddShiftEmployee - Failed to delete from distributed cache", "error", delErr)
	}
	if delErr := s.localCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("AddShiftEmployee - Failed to delete from local cache", "error", delErr)
	}

	s.logger.Info("AddShiftEmployee - Success", "employee_id", input.EmployeeId)

	// Note: The repository doesn't return the ID, so we can't populate ShiftUserId
	// You may need to modify the repository interface to return the ID
	output := &applicationModel.AddShiftEmployeeOutput{
		ShiftUserId: "created", // Placeholder - need to modify repository to return ID
	}

	return &output, nil
}

// DeleteShiftUser implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) DeleteShiftUser(ctx context.Context, input *applicationModel.DeleteShiftUserInput) *applicationError.Error {
	// Validate input
	if input == nil {
		s.logger.Error("DeleteShiftUser - Input is nil")
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("DeleteShiftUser - Start", "user_id", input.UserId, "shift_user_id", input.ShiftUserId)

	// Call repository
	err := s.shiftUserRepo.DeleteEmployeeShift(ctx, input.ShiftUserId)
	if err != nil {
		s.logger.Error("DeleteShiftUser - Failed to delete shift assignment", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to delete shift assignment",
		}
	}

	s.logger.Info("DeleteShiftUser - Success", "shift_user_id", input.ShiftUserId)

	return nil
}

// DisableShiftUser implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) DisableShiftUser(ctx context.Context, input *applicationModel.DisableShiftUserInput) *applicationError.Error {
	// Validate input
	if input == nil {
		s.logger.Error("DisableShiftUser - Input is nil")
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("DisableShiftUser - Start", "user_id", input.UserId, "shift_user_id", input.ShiftUserId)

	// Call repository
	err := s.shiftUserRepo.DisableEmployeeShift(ctx, input.ShiftUserId)
	if err != nil {
		s.logger.Error("DisableShiftUser - Failed to disable shift assignment", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to disable shift assignment",
		}
	}

	s.logger.Info("DisableShiftUser - Success", "shift_user_id", input.ShiftUserId)

	return nil
}

// EditShiftForUserWithEffectiveDate implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) EditShiftForUserWithEffectiveDate(ctx context.Context, input *applicationModel.EditShiftForUserWithEffectiveDateInput) *applicationError.Error {
	// Validate input
	if input == nil {
		s.logger.Error("EditShiftForUserWithEffectiveDate - Input is nil")
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("EditShiftForUserWithEffectiveDate - Start", "user_id", input.UserId, "shift_user_id", input.ShiftUserId)

	// Create domain input
	domainInput := &domainModel.EditEffectiveShiftForEmployeeInput{
		EmployeeShiftID: input.ShiftUserId,
		EffectiveFrom:   input.NewEffectiveFrom,
		EffectiveTo:     input.NewEffectiveTo,
	}

	// Call repository
	err := s.shiftUserRepo.EditEffectiveShiftForEmployee(ctx, domainInput)
	if err != nil {
		s.logger.Error("EditShiftForUserWithEffectiveDate - Failed to edit shift effective date", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to edit shift effective date",
		}
	}

	s.logger.Info("EditShiftForUserWithEffectiveDate - Success", "shift_user_id", input.ShiftUserId)

	return nil
}

// EnableShiftUser implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) EnableShiftUser(ctx context.Context, input *applicationModel.EnableShiftUserInput) *applicationError.Error {
	// Validate input
	if input == nil {
		s.logger.Error("EnableShiftUser - Input is nil")
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("EnableShiftUser - Start", "user_id", input.UserId, "shift_user_id", input.ShiftUserId)

	// Call repository
	err := s.shiftUserRepo.EnableEmployeeShift(ctx, input.ShiftUserId)
	if err != nil {
		s.logger.Error("EnableShiftUser - Failed to enable shift assignment", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to enable shift assignment",
		}
	}

	s.logger.Info("EnableShiftUser - Success", "shift_user_id", input.ShiftUserId)

	return nil
}

// GetShiftForUserWithEffectiveDate implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) GetShiftForUserWithEffectiveDate(ctx context.Context, input *applicationModel.GetShiftForUserWithEffectiveDateInput) (*applicationModel.GetShiftForUserWithEffectiveDateOutput, *applicationError.Error) {
	// Validate input
	if input == nil {
		s.logger.Error("GetShiftForUserWithEffectiveDate - Input is nil")
		return nil, &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}

	s.logger.Info("GetShiftForUserWithEffectiveDate - Start", "user_id", input.UserId, "page", input.Page, "size", input.Size)

	// Create cache key based on user and date range
	cacheKey := fmt.Sprintf("%s%s:%s:%s:%d:%d",
		shiftEmployeeCachePrefix,
		input.UserId.String(),
		input.EffectiveFrom.Format("2006-01-02"),
		input.EffectiveTo.Format("2006-01-02"),
		input.Page,
		input.Size,
	)

	// Try to get from local cache first
	if cachedData, err := s.localCache.Get(ctx, cacheKey); err == nil && cachedData != "" {
		s.logger.Info("GetShiftForUserWithEffectiveDate - Cache hit (local)", "user_id", input.UserId)
		var output applicationModel.GetShiftForUserWithEffectiveDateOutput
		if unmarshalErr := json.Unmarshal([]byte(cachedData), &output); unmarshalErr == nil {
			return &output, nil
		} else {
			s.logger.Warn("GetShiftForUserWithEffectiveDate - Failed to unmarshal local cache", "error", unmarshalErr)
		}
	}

	// Try to get from distributed cache
	if cachedData, err := s.distributedCache.Get(ctx, cacheKey); err == nil && cachedData != "" {
		s.logger.Info("GetShiftForUserWithEffectiveDate - Cache hit (distributed)", "user_id", input.UserId)
		var output applicationModel.GetShiftForUserWithEffectiveDateOutput
		if unmarshalErr := json.Unmarshal([]byte(cachedData), &output); unmarshalErr == nil {
			// Store in local cache for faster access next time
			if jsonData, _ := json.Marshal(output); len(jsonData) > 0 {
				_ = s.localCache.SetTTL(ctx, cacheKey, string(jsonData), shiftEmployeeCacheTTL)
			}
			return &output, nil
		} else {
			s.logger.Warn("GetShiftForUserWithEffectiveDate - Failed to unmarshal distributed cache", "error", unmarshalErr)
		}
	}

	s.logger.Info("GetShiftForUserWithEffectiveDate - Cache miss, fetching from database", "user_id", input.UserId)

	// Calculate offset from page and size
	offset := int32((input.Page - 1) * input.Size)
	limit := int32(input.Size)

	// Create domain input
	domainInput := &domainModel.GetShiftEmployeeWithEffectiveDateInput{
		EmployeeID:    input.UserId, // Using UserId as EmployeeID
		EffectiveFrom: input.EffectiveFrom,
		Limit:         limit,
		Offset:        offset,
	}

	// Call repository
	shifts, err := s.shiftUserRepo.GetShiftEmployeeWithEffectiveDate(ctx, domainInput)
	if err != nil {
		s.logger.Error("GetShiftForUserWithEffectiveDate - Failed to get shifts for employee", "error", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to get shifts for employee",
		}
	}

	s.logger.Info("GetShiftForUserWithEffectiveDate - Found shifts", "count", len(shifts))

	// For now, returning empty output as the structure is not fully defined
	// You may need to populate this based on the shifts returned
	_ = shifts // Use the shifts variable to avoid unused warning

	output := &applicationModel.GetShiftForUserWithEffectiveDateOutput{}

	// Cache the result
	if jsonData, jsonErr := json.Marshal(output); jsonErr == nil {
		// Store in distributed cache
		if setErr := s.distributedCache.SetTTL(ctx, cacheKey, string(jsonData), shiftEmployeeCacheTTL); setErr != nil {
			s.logger.Warn("GetShiftForUserWithEffectiveDate - Failed to set distributed cache", "error", setErr)
		}
		// Store in local cache
		if setErr := s.localCache.SetTTL(ctx, cacheKey, string(jsonData), shiftEmployeeCacheTTL); setErr != nil {
			s.logger.Warn("GetShiftForUserWithEffectiveDate - Failed to set local cache", "error", setErr)
		}
	}

	s.logger.Info("GetShiftForUserWithEffectiveDate - Success", "user_id", input.UserId)

	return output, nil
}

// New instance
func NewShiftEmployeeService() service.IShiftEmployeeService {
	shiftUserRepo, err := repository.GetShiftUserRepository()
	if err != nil {
		panic(fmt.Sprintf("Failed to get shift user repository: %v", err))
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

	return &ShiftEmployeeService{
		shiftUserRepo:    shiftUserRepo,
		logger:           log,
		distributedCache: distributedCache,
		localCache:       localCache,
	}
}
