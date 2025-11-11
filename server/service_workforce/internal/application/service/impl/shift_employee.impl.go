package service

import (
	// "context"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
	service "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/cache"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/logger"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/repository"
	utilsCache "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/shared/utils/cache"
	utilsCrypto "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/shared/utils/crypto"
)

// =================================================
// ShiftEmployee service implementation interface
// =================================================
type ShiftEmployeeService struct {
	shiftUserRepo    repository.IShiftUserRepository
	shiftRepo        repository.IShiftRepository
	userRepo         repository.IUserRepository
	logger           logger.ILogger
	distributedCache cache.IDistributedCache
	localCache       cache.ILocalCache
}

// GetInfoEmployeeInShift implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) GetInfoEmployeeInShift(ctx context.Context, input *applicationModel.GetInfoEmployeeInShiftInput) ([]*applicationModel.EmployeeInfoInShiftBase, *applicationError.Error) {
	if input == nil {
		s.logger.Error("GetInfoEmployeeInShift - Input is nil")
		return nil, &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}
	// Check permission
	if input.Role == domainModel.RoleManager {
		ok, err := s.shiftUserRepo.IsUserManagetShift(
			ctx,
			&domainModel.IsUserManagetShiftInput{
				ShiftID:       input.ShiftId,
				CompanyUserID: input.CompanyId,
			},
		)
		if err != nil {
			s.logger.Error("GetInfoEmployeeInShift - Failed to check user permission", "error", err)
			return nil, &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "Failed to check user permission",
			}
		}
		if !ok {
			s.logger.Error("GetInfoEmployeeInShift - User does not have permission to view shift employee info", "user_id", input.UserId, "shift_id", input.ShiftId)
			return nil, &applicationError.Error{
				ErrorSystem: fmt.Errorf("user does not have permission to view shift employee info"),
				ErrorClient: "You do not have permission to view this shift employee info",
			}
		}
	}
	// Get data in cache
	key := utilsCache.GetKeyInfoEmployeeInShift(
		utilsCrypto.GetHash(input.ShiftId.String()),
	)

	// Try local cache first
	if cachedData, err := s.localCache.Get(ctx, key); err == nil && cachedData != "" {
		var employeeInfo []*applicationModel.EmployeeInfoInShiftBase
		if err := json.Unmarshal([]byte(cachedData), &employeeInfo); err == nil {
			s.logger.Info("GetInfoEmployeeInShift - Cache hit (local)", "shift_id", input.ShiftId)
			return employeeInfo, nil
		}
		s.logger.Warn("GetInfoEmployeeInShift - Failed to unmarshal local cache", "error", err)
	}

	// Try distributed cache next
	if cachedData, err := s.distributedCache.Get(ctx, key); err == nil && cachedData != "" {
		var employeeInfo []*applicationModel.EmployeeInfoInShiftBase
		if err := json.Unmarshal([]byte(cachedData), &employeeInfo); err == nil {
			s.logger.Info("GetInfoEmployeeInShift - Cache hit (distributed)", "shift_id", input.ShiftId)
			return employeeInfo, nil
		}
		s.logger.Warn("GetInfoEmployeeInShift - Failed to unmarshal distributed cache", "error", err)
	}

	// Cache miss — proceed to fetch from source
	s.logger.Info("GetInfoEmployeeInShift - Cache miss", "shift_id", input.ShiftId)

	// Call repository
	reps, err := s.shiftRepo.GetInfoEmployeeInShift(
		ctx,
		&domainModel.GetInfoEmployeeInShiftInput{
			CompanyId: input.CompanyId,
			CurShidId: input.ShiftId,
		},
	)
	if err != nil {
		s.logger.Error("GetInfoEmployeeInShift - Failed to get employee info in shift", "error", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to get employee info in shift",
		}
	}
	if len(reps) == 0 {
		s.logger.Info("GetInfoEmployeeInShift - No employees found in shift", "shift_id", input.ShiftId)
		return []*applicationModel.EmployeeInfoInShiftBase{}, nil
	}

	// Map domain model to application model
	output := make([]*applicationModel.EmployeeInfoInShiftBase, 0)
	for _, r := range reps {
		output = append(output, &applicationModel.EmployeeInfoInShiftBase{
			UserId:         r.UserId.String(),
			Name:           r.Name,
			NumberEmployee: r.NumberEmployee,
			CurrentShift:   r.CurrentShift,
			ShiftActive:    r.ShiftActive,
		})
	}

	// Cache the result for future requests
	if data, err := json.Marshal(output); err == nil {
		// Store in local cache (ignore errors - cache is non-critical)
		_ = s.localCache.SetTTL(ctx, key, string(data), 2) // adjust TTL as needed
		// Store in distributed cache
		_ = s.distributedCache.SetTTL(ctx, key, string(data), constants.TTL_Info_Base_Employee_In_Shift)
	}

	return output, nil
}

// AddListShiftEmployee implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) AddListShiftEmployee(ctx context.Context, input *applicationModel.AddShiftEmployeeListInput) *applicationError.Error {
	// Validate input
	if input == nil {
		s.logger.Error("AddListShiftEmployee - Input is nil")
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("input is nil"),
			ErrorClient: "Invalid input data",
		}
	}
	// Get company req and check permission
	var companyId uuid.UUID
	if input.CompanyRequestId == input.CompanyId {
		companyId = input.CompanyRequestId
	} else if input.Role == domainModel.RoleAdmin {
		companyId = input.CompanyRequestId
	} else {
		s.logger.Error("AddListShiftEmployee - User does not have permission to add employees to this company", "user_id", input.UserId, "company_request_id", input.CompanyRequestId, "company_id", input.CompanyId)
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("user does not have permission to add employees to this company"),
			ErrorClient: "You do not have permission to add employees to this company",
		}
	}
	s.logger.Info("AddListShiftEmployee - Start", "user_id", input.UserId, "number_of_employees", len(input.EmployeeIDs))
	// Check user exist shift in time range
	for _, empId := range input.EmployeeIDs {
		checkInput := &domainModel.CheckUserExistShiftInput{
			EmployeeID:    empId,
			EffectiveFrom: input.EffectiveFrom,
			EffectiveTo:   input.EffectiveTo,
			Limit:         1,
			Offset:        0,
		}

		exists, err := s.shiftUserRepo.CheckUserExistShift(ctx, checkInput)
		if err != nil {
			s.logger.Error("AddListShiftEmployee - Failed to check existing shift for employee", "employee_id", empId, "error", err)
			return &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "Failed to check existing shift for employee",
			}
		}

		if exists {
			s.logger.Warn("AddListShiftEmployee - Employee already has a shift in this time range", "employee_id", empId)
			return &applicationError.Error{
				ErrorSystem: fmt.Errorf("employee already has a shift in this time range"),
				ErrorClient: fmt.Sprintf("Employee with ID %s already has a shift in this time range", empId.String()),
			}
		}
	}
	
	// Call repository
	reqRepo := &domainModel.AddListShiftForEmployeesInput{
		CompanyID:     companyId,
		ShiftID:       input.ShiftId,
		EmployeeIDs:   input.EmployeeIDs,
		EffectiveFrom: input.EffectiveFrom,
		EffectiveTo:   input.EffectiveTo,
	}
	err := s.shiftUserRepo.AddListShiftForEmployees(
		ctx,
		reqRepo,
	)
	if err != nil {
		s.logger.Error("AddListShiftEmployee - Failed to add shifts to employees", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to add shifts to employees",
		}
	}
	// rm cache list workforce company
	key := utilsCache.GetKeyListShiftInCompanyPrefix(
		utilsCrypto.GetHash(input.CompanyId.String()),
	)
	if err := s.distributedCache.DeleteByPrefix(ctx, key); err != nil {
		s.logger.Warn("AddListShiftEmployee - Failed to delete list shift cache in company", "error", err)
	}

	return nil
}

// AddShiftEmployee implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) AddShiftEmployee(ctx context.Context, input *applicationModel.AddShiftEmployeeInput) *applicationError.Error {
	// Validate input
	if input == nil {
		s.logger.Error("AddShiftEmployee - Input is nil")
		return &applicationError.Error{
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
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to check existing shift",
		}
	}

	if exists {
		s.logger.Warn("AddShiftEmployee - Employee already has a shift in this time range", "employee_id", input.EmployeeId)
		return &applicationError.Error{
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
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to add shift to employee",
		}
	}

	// Invalidate cache for this employee
	cacheKey := utilsCache.GetKeyShiftEmployee(
		utilsCrypto.GetHash(input.ShiftId.String()),
		utilsCrypto.GetHash(input.EmployeeId.String()),
	)
	if delErr := s.distributedCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("AddShiftEmployee - Failed to delete from distributed cache", "error", delErr)
	}
	if delErr := s.localCache.Delete(ctx, cacheKey); delErr != nil {
		s.logger.Warn("AddShiftEmployee - Failed to delete from local cache", "error", delErr)
	}

	s.logger.Info("AddShiftEmployee - Success", "employee_id", input.EmployeeId)

	return nil
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
	// Get company user and check permission
	isUserInCompany, err := s.userRepo.UserExistsInCompany(
		ctx,
		&domainModel.UserExistsInCompanyInput{
			CompanyID: input.CompanyId,
			UserID:    input.UserIdReq,
		},
	)
	if err != nil {
		s.logger.Error("DeleteShiftUser - Failed to check user in company", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to check user in company",
		}
	}
	if !isUserInCompany && input.Role != domainModel.RoleAdmin {
		s.logger.Error("DeleteShiftUser - User does not have permission to delete shift assignment", "user_id", input.UserId, "user_id_req", input.UserIdReq, "company_id", input.CompanyId)
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("user does not have permission to delete shift assignment"),
			ErrorClient: "You do not have permission to delete this shift assignment",
		}
	}
	s.logger.Info("DeleteShiftUser - Start", "user_id", input.UserId, "shift_user_id", input.ShiftId)
	// Call repository
	if err := s.shiftUserRepo.DeleteEmployeeShift(
		ctx,
		&domainModel.DeleteEmployeeShiftInput{
			ShiftId:    input.ShiftId,
			EmployeeID: input.UserIdReq,
		},
	); err != nil {
		s.logger.Error("DeleteShiftUser - Failed to delete shift assignment", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to delete shift assignment",
		}
	}

	s.logger.Info("DeleteShiftUser - Success", "shift_id", input.ShiftId, "user_id_req", input.UserIdReq)
	// rm cache list workforce company
	keyInList := utilsCache.GetKeyListShiftInCompanyPrefix(
		utilsCrypto.GetHash(input.CompanyId.String()),
	)
	keyShiftUserInfo := utilsCache.GetKeyShiftEmployee(
		utilsCrypto.GetHash(input.ShiftId.String()),
		utilsCrypto.GetHash(input.UserIdReq.String()),
	)
	if err := s.distributedCache.DeleteByPrefix(ctx, keyInList); err != nil {
		s.logger.Warn("DeleteShiftUser - Failed to delete list shift cache in company", "error", err)
	}
	if err := s.distributedCache.Delete(ctx, keyShiftUserInfo); err != nil {
		s.logger.Warn("DeleteShiftUser - Failed to delete shift user info cache", "error", err)
	}

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
	// Get company user and check permission
	isUserInCompany, err := s.userRepo.UserExistsInCompany(
		ctx,
		&domainModel.UserExistsInCompanyInput{
			CompanyID: input.CompanyId,
			UserID:    input.UserIdReq,
		},
	)
	if err != nil {
		s.logger.Error("DeleteShiftUser - Failed to check user in company", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to check user in company",
		}
	}
	if !isUserInCompany && input.Role != domainModel.RoleAdmin {
		s.logger.Error("DeleteShiftUser - User does not have permission to delete shift assignment", "user_id", input.UserId, "user_id_req", input.UserIdReq, "company_id", input.CompanyId)
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("user does not have permission to delete shift assignment"),
			ErrorClient: "You do not have permission to delete this shift assignment",
		}
	}
	s.logger.Info("DisableShiftUser - Start", "user_id", input.UserId, "shift_id", input.ShiftId)

	// Call repository
	if err := s.shiftUserRepo.DisableEmployeeShift(ctx, &domainModel.DisableEmployeeShiftInput{
		ShiftID:    input.ShiftId,
		EmployeeID: input.UserIdReq,
	}); err != nil {
		s.logger.Error("DisableShiftUser - Failed to disable shift assignment", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to disable shift assignment",
		}
	}
	s.logger.Info("DisableShiftUser - Success", "shift_id", input.ShiftId, "user_id_req", input.UserIdReq)
	// Rm cache
	key := utilsCache.GetKeyShiftEmployee(
		utilsCrypto.GetHash(input.ShiftId.String()),
		utilsCrypto.GetHash(input.UserIdReq.String()),
	)
	if err := s.distributedCache.Delete(ctx, key); err != nil {
		s.logger.Warn("DisableShiftUser - Failed to delete from distributed cache", "error", err)
	}
	if err := s.localCache.Delete(ctx, key); err != nil {
		s.logger.Warn("DisableShiftUser - Failed to delete from local cache", "error", err)
	}
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
	// Get company user and check permission
	isUserInCompany, err := s.userRepo.UserExistsInCompany(
		ctx,
		&domainModel.UserExistsInCompanyInput{
			CompanyID: input.CompanyId,
			UserID:    input.UserIdReq,
		},
	)
	if err != nil {
		s.logger.Error("DeleteShiftUser - Failed to check user in company", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to check user in company",
		}
	}
	if !isUserInCompany && input.Role != domainModel.RoleAdmin {
		s.logger.Error("DeleteShiftUser - User does not have permission to delete shift assignment", "user_id", input.UserId, "user_id_req", input.UserIdReq, "company_id", input.CompanyId)
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("user does not have permission to delete shift assignment"),
			ErrorClient: "You do not have permission to delete this shift assignment",
		}
	}
	// Create domain input
	domainInput := &domainModel.EditEffectiveShiftForEmployeeInput{
		EmployeeID:    input.UserIdReq,
		ShiftID:       input.ShiftId,
		EffectiveFrom: input.NewEffectiveFrom,
		EffectiveTo:   input.NewEffectiveTo,
	}

	// Call repository
	if err := s.shiftUserRepo.EditEffectiveShiftForEmployee(ctx, domainInput); err != nil {
		s.logger.Error("EditShiftForUserWithEffectiveDate - Failed to edit shift effective date", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to edit shift effective date",
		}
	}
	// Rm cache
	key := utilsCache.GetKeyShiftEmployee(
		utilsCrypto.GetHash(input.ShiftId.String()),
		utilsCrypto.GetHash(input.UserIdReq.String()),
	)
	if err := s.distributedCache.Delete(ctx, key); err != nil {
		s.logger.Warn("EditShiftForUserWithEffectiveDate - Failed to delete from distributed cache", "error", err)
	}
	if err := s.localCache.Delete(ctx, key); err != nil {
		s.logger.Warn("EditShiftForUserWithEffectiveDate - Failed to delete from local cache", "error", err)
	}

	s.logger.Info("EditShiftForUserWithEffectiveDate - Success", "shift_id", input.ShiftId, "user_id_req", input.UserIdReq)

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
	// Get company user and check permission
	isUserInCompany, err := s.userRepo.UserExistsInCompany(
		ctx,
		&domainModel.UserExistsInCompanyInput{
			CompanyID: input.CompanyId,
			UserID:    input.UserIdReq,
		},
	)
	if err != nil {
		s.logger.Error("DeleteShiftUser - Failed to check user in company", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to check user in company",
		}
	}
	if !isUserInCompany && input.Role != domainModel.RoleAdmin {
		s.logger.Error("DeleteShiftUser - User does not have permission to delete shift assignment", "user_id", input.UserId, "user_id_req", input.UserIdReq, "company_id", input.CompanyId)
		return &applicationError.Error{
			ErrorSystem: fmt.Errorf("user does not have permission to delete shift assignment"),
			ErrorClient: "You do not have permission to delete this shift assignment",
		}
	}
	// Call repository
	if err := s.shiftUserRepo.EnableEmployeeShift(ctx,
		&domainModel.EnableEmployeeShiftIInput{
			ShiftID:    input.ShiftId,
			EmployeeID: input.UserIdReq,
		}); err != nil {
		s.logger.Error("EnableShiftUser - Failed to enable shift assignment", "error", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "Failed to enable shift assignment",
		}
	}
	// Rm cache
	key := utilsCache.GetKeyShiftEmployee(
		utilsCrypto.GetHash(input.ShiftId.String()),
		utilsCrypto.GetHash(input.UserIdReq.String()),
	)
	if err := s.distributedCache.Delete(ctx, key); err != nil {
		s.logger.Warn("EnableShiftUser - Failed to delete from distributed cache", "error", err)
	}
	if err := s.localCache.Delete(ctx, key); err != nil {
		s.logger.Warn("EnableShiftUser - Failed to delete from local cache", "error", err)
	}

	s.logger.Info("EnableShiftUser - Success", "shift_id", input.ShiftId, "user_id_req", input.UserIdReq)

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
	cacheKey := utilsCache.GetKeyShiftEmployeeWithEffectiveDate(
		utilsCrypto.GetHash(input.UserId.String()),
		fmt.Sprintf("%d", input.EffectiveFrom.Unix()),
		fmt.Sprintf("%d", input.EffectiveTo.Unix()),
		input.Page,
		input.Size,
	)
	shiftEmployeeCacheTTL := int64(constants.TTL_Shift_Employee_Cache)
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

	shiftRepo, err := repository.GetShiftRepository()
	if err != nil {
		panic(fmt.Sprintf("Failed to get shift repository: %v", err))
	}

	userRepo, err := repository.GetUserRepository()
	if err != nil {
		panic(fmt.Sprintf("Failed to get user repository: %v", err))
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
		shiftRepo:        shiftRepo,
		userRepo:         userRepo,
		logger:           log,
		distributedCache: distributedCache,
		localCache:       localCache,
	}
}
