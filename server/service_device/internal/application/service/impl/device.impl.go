package service

import (
	"context"
	json "encoding/json"

	"github.com/google/uuid"
	applicationError "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/error"
	model "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/model"
	service "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/service"
	constants "github.com/youknow2509/cio_verify_face/server/service_device/internal/constants"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/cache"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/repository"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/token"
	global "github.com/youknow2509/cio_verify_face/server/service_device/internal/global"
	sharedCache "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/cache"
	sharedCrypto "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/crypto"
	sharedRandom "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/random"
)

// =================================================
// Device application service
// =================================================
type DeviceService struct{}

// UpdateStatusDevice implements service.IDeviceService.
func (d *DeviceService) UpdateStatusDevice(ctx context.Context, input *model.UpdateStatusDeviceInput) *applicationError.Error {
	// Check permission
	if input.Role > 1 {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to update device status.",
		}
	}
	// Check device exist
	deviceRepo, _ := domainRepo.GetDeviceRepository()
	ok, err := deviceRepo.DeviceExist(ctx, &domainModel.DeviceExistInput{DeviceId: input.DeviceId})
	if err != nil {
		global.Logger.Error("Error when get device by id", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !ok {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "Device not found.",
		}
	}
	// Update status device
	if input.Status == 1 {
		if err := deviceRepo.EnableDevice(ctx, &domainModel.EnableDeviceInput{DeviceId: input.DeviceId}); err != nil {
			global.Logger.Error("Error when enable device", "err", err)
			return &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "System is busy now. Please try again later.",
			}
		}
	} else {
		if err := deviceRepo.DisableDevice(ctx, &domainModel.DisableDeviceInput{DeviceId: input.DeviceId}); err != nil {
			global.Logger.Error("Error when disable device", "err", err)
			return &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "System is busy now. Please try again later.",
			}
		}
	}
	// Rm cache of device info
	userRepo, _ := domainRepo.GetUserRepository()
	companyInfo, err := userRepo.GetCompanyIdOfUser(
		ctx,
		&domainModel.GetCompanyIdOfUserInput{
			UserID: input.UserId,
		},
	)
	if err != nil {
		return nil
	}
	if companyInfo == nil {
		return nil
	}
	key := []string{
		sharedCache.GetKeyDeviceBase(sharedCrypto.GetHash(input.DeviceId.String())),
		sharedCache.GetKeyListDeviceInCompany(
			sharedCrypto.GetHash(companyInfo.CompanyID.String()),
			20,
			1,
		),
	}
	go func() {
		cacheService, _ := domainCache.GetDistributedCache()
		for _, k := range key {
			if err := cacheService.Delete(context.Background(), k); err != nil {
				global.Logger.Error("Error when delete device info cache", "err", err)
			}
		}
	}()
	return nil
}

// RefreshDeviceToken implements service.IDeviceService.
func (d *DeviceService) RefreshDeviceToken(ctx context.Context, input *model.RefreshDeviceTokenInput) (*model.RefreshDeviceTokenOutput, *applicationError.Error) {
	// Check user have permission to get device token
	if input.Role > 1 {
		return nil, &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to get device token.",
		}
	}
	if input.Role == domainModel.RoleManager {
		// Check user in company
		userRepo, _ := domainRepo.GetUserRepository()
		userInfo, err := userRepo.UserPermissionDevice(ctx, &domainModel.UserPermissionDeviceInput{
			UserID:   input.UserId,
			DeviceID: input.DeviceId,
		})
		if err != nil {
			global.Logger.Error("Error when check user permission device", "err", err)
			return nil, &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "System is busy now. Please try again later.",
			}
		}
		if !userInfo {
			return nil, &applicationError.Error{
				ErrorSystem: nil,
				ErrorClient: "You don't have permission to get device token.",
			}
		}
	}
	// Call to grpc service to refresh token and create new token
	domainToken := domainToken.GetTokenService()
	newToken, err := domainToken.CreateDeviceToken(ctx, &domainModel.TokenDeviceJwtInput{
		DeviceId: input.DeviceId.String(),
		TokenId:  uuid.New().String(),
	})
	if err != nil {
		global.Logger.Error("Error when create device token", "err", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if newToken == "" {
		return nil, &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "Failed to create new device token.",
		}
	}
	// save new token to cache
	distributedCacheService, _ := domainCache.GetDistributedCache()
	key := sharedCache.GetKeyDeviceToken(sharedCrypto.GetHash(input.DeviceId.String()))
	distributedCacheService.SetTTL(
		ctx,
		key,
		newToken,
		constants.TTL_DEVICE_TOKEN,
	)
	// Rm cache of device info
	userRepo, _ := domainRepo.GetUserRepository()
	companyInfo, err := userRepo.GetCompanyIdOfUser(
		ctx,
		&domainModel.GetCompanyIdOfUserInput{
			UserID: input.UserId,
		},
	)
	if err != nil {
		return nil, nil
	}
	if companyInfo == nil {
		return nil, nil
	}
	keyRm := []string{
		sharedCache.GetKeyDeviceBase(sharedCrypto.GetHash(input.DeviceId.String())),
		sharedCache.GetKeyListDeviceInCompany(
			sharedCrypto.GetHash(companyInfo.CompanyID.String()),
			20,
			1,
		),
	}
	go func() {
		cacheService, _ := domainCache.GetDistributedCache()
		for _, k := range keyRm {
			if err := cacheService.Delete(context.Background(), k); err != nil {
				global.Logger.Error("Error when delete device info cache", "err", err)
			}
		}
	}()
	return &model.RefreshDeviceTokenOutput{
		DeviceId:    input.DeviceId.String(),
		DeviceToken: newToken,
	}, nil
}

// GetDeviceToken implements service.IDeviceService.
func (d *DeviceService) GetDeviceToken(ctx context.Context, input *model.GetDeviceTokenInput) (*model.GetDeviceTokenOutput, *applicationError.Error) {
	// Check user have permission to get device token
	if input.Role > 1 {
		return nil, &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to get device token.",
		}
	}
	if input.Role == domainModel.RoleManager {
		// Check user in company
		userRepo, _ := domainRepo.GetUserRepository()
		userInfo, err := userRepo.UserPermissionDevice(ctx, &domainModel.UserPermissionDeviceInput{
			UserID:   input.UserId,
			DeviceID: input.DeviceId,
		})
		if err != nil {
			global.Logger.Error("Error when check user permission device", "err", err)
			return nil, &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "System is busy now. Please try again later.",
			}
		}
		if !userInfo {
			return nil, &applicationError.Error{
				ErrorSystem: nil,
				ErrorClient: "You don't have permission to get device token.",
			}
		}
	}
	// Check cache
	key := sharedCache.GetKeyDeviceToken(sharedCrypto.GetHash(input.DeviceId.String()))
	distributedCacheService, _ := domainCache.GetDistributedCache()
	deviceTokenCacheStr, err := distributedCacheService.Get(ctx, key)
	if err != nil {
		global.Logger.Error("Error when get device token from cache", "err", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if deviceTokenCacheStr != "" {
		return &model.GetDeviceTokenOutput{
			DeviceId:    input.DeviceId.String(),
			DeviceToken: deviceTokenCacheStr,
		}, nil
	}
	// Get device token
	deviceRepo, _ := domainRepo.GetDeviceRepository()
	deviceToken, err := deviceRepo.GetDeviceToken(
		ctx,
		&domainModel.GetDeviceTokenInput{
			DeviceId: input.DeviceId,
		},
	)
	if err != nil {
		global.Logger.Error("Error when get device token", "err", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	// Save in cache
	ttl := constants.TTL_DEVICE_TOKEN
	if err := distributedCacheService.SetTTL(
		ctx,
		key,
		deviceToken.Token,
		int64(ttl),
	); err != nil {
		global.Logger.Error("Error when set device token in cache", "err", err)
		// Not return error if cache error
	}
	return &model.GetDeviceTokenOutput{
		DeviceId:    input.DeviceId.String(),
		DeviceToken: deviceToken.Token,
	}, nil
}

// UpdateInfoDevice implements service.IDeviceService.
func (d *DeviceService) UpdateInfoDevice(ctx context.Context, input *model.UpdateInfoDeviceInput) *applicationError.Error {
	// Check permission
	domainUser, err := domainRepo.GetUserRepository()
	if err != nil {
		global.Logger.Error("Error when get user repository", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if input.Role > 1 {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to update device info.",
		}
	}
	ok, err := domainUser.UserPermissionDevice(ctx, &domainModel.UserPermissionDeviceInput{
		UserID:   input.UserId,
		DeviceID: input.DeviceId,
	})
	if err != nil {
		global.Logger.Error("Error when check user permission device", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !ok && input.Role != 0 {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to update device info.",
		}
	}
	// Check device exist
	deviceRepo, _ := domainRepo.GetDeviceRepository()
	exist, err := deviceRepo.DeviceExist(ctx, &domainModel.DeviceExistInput{
		DeviceId: input.DeviceId,
	})
	if err != nil {
		global.Logger.Error("Error when check device exist", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !exist {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "Device not found.",
		}
	}
	// Update device info
	if err := deviceRepo.UpdateDeviceInfo(
		ctx,
		&domainModel.UpdateDeviceInfoInput{},
	); err != nil {
		global.Logger.Error("Error when update device info", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	// Rm cache of device info
	userRepo, _ := domainRepo.GetUserRepository()
	companyInfo, err := userRepo.GetCompanyIdOfUser(
		ctx,
		&domainModel.GetCompanyIdOfUserInput{
			UserID: input.UserId,
		},
	)
	if err != nil {
		return nil
	}
	if companyInfo == nil {
		return nil
	}
	key := []string{
		sharedCache.GetKeyDeviceBase(sharedCrypto.GetHash(input.DeviceId.String())),
		sharedCache.GetKeyListDeviceInCompany(
			sharedCrypto.GetHash(companyInfo.CompanyID.String()),
			20,
			1,
		),
	}
	go func() {
		cacheService, _ := domainCache.GetDistributedCache()
		for _, k := range key {
			if err := cacheService.Delete(context.Background(), k); err != nil {
				global.Logger.Error("Error when delete device info cache", "err", err)
			}
		}
	}()
	return nil
}

// UpdateLocationDevice implements service.IDeviceService.
func (d *DeviceService) UpdateLocationDevice(ctx context.Context, input *model.UpdateLocationDeviceInput) *applicationError.Error {
	// Check permission
	domainUser, err := domainRepo.GetUserRepository()
	if err != nil {
		global.Logger.Error("Error when get user repository", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if input.Role > 1 {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to update device info.",
		}
	}
	ok, err := domainUser.UserPermissionDevice(ctx, &domainModel.UserPermissionDeviceInput{
		UserID:   input.UserId,
		DeviceID: input.DeviceId,
	})
	if err != nil {
		global.Logger.Error("Error when check user permission device", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !ok && input.Role != 0 {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to update device info.",
		}
	}
	// Check device exist
	deviceRepo, _ := domainRepo.GetDeviceRepository()
	exist, err := deviceRepo.DeviceExist(ctx, &domainModel.DeviceExistInput{
		DeviceId: input.DeviceId,
	})
	if err != nil {
		global.Logger.Error("Error when check device exist", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !exist {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "Device not found.",
		}
	}
	// Update device location
	if err := deviceRepo.UpdateDeviceLocation(
		ctx,
		&domainModel.UpdateDeviceLocationInput{
			DeviceId:   input.DeviceId,
			LocationId: input.NewLocationId,
			Address:    input.NewAddress,
		},
	); err != nil {
		global.Logger.Error("Error when update device location", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	// Rm cache of device info
	userRepo, _ := domainRepo.GetUserRepository()
	companyInfo, err := userRepo.GetCompanyIdOfUser(
		ctx,
		&domainModel.GetCompanyIdOfUserInput{
			UserID: input.UserId,
		},
	)
	if err != nil {
		return nil
	}
	if companyInfo == nil {
		return nil
	}
	key := []string{
		sharedCache.GetKeyDeviceBase(sharedCrypto.GetHash(input.DeviceId.String())),
		sharedCache.GetKeyListDeviceInCompany(
			sharedCrypto.GetHash(companyInfo.CompanyID.String()),
			20,
			1,
		),
	}
	go func() {
		cacheService, _ := domainCache.GetDistributedCache()
		for _, k := range key {
			if err := cacheService.Delete(context.Background(), k); err != nil {
				global.Logger.Error("Error when delete device info cache", "err", err)
			}
		}
	}()
	return nil
}

// UpdateNameDevice implements service.IDeviceService.
func (d *DeviceService) UpdateNameDevice(ctx context.Context, input *model.UpdateNameDeviceInput) *applicationError.Error {
	// Check permission
	domainUser, err := domainRepo.GetUserRepository()
	if err != nil {
		global.Logger.Error("Error when get user repository", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if input.Role > 1 {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to update device info.",
		}
	}
	ok, err := domainUser.UserPermissionDevice(ctx, &domainModel.UserPermissionDeviceInput{
		UserID:   input.UserId,
		DeviceID: input.DeviceId,
	})
	if err != nil {
		global.Logger.Error("Error when check user permission device", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !ok && input.Role != 0 {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to update device info.",
		}
	}
	// Check device exist
	deviceRepo, _ := domainRepo.GetDeviceRepository()
	exist, err := deviceRepo.DeviceExist(ctx, &domainModel.DeviceExistInput{
		DeviceId: input.DeviceId,
	})
	if err != nil {
		global.Logger.Error("Error when check device exist", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !exist {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "Device not found.",
		}
	}
	// Update device name
	if err := deviceRepo.UpdateDeviceName(
		ctx,
		&domainModel.UpdateDeviceNameInput{
			DeviceId: input.DeviceId,
			Name:     input.NewName,
		},
	); err != nil {
		global.Logger.Error("Error when update device name", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	// Rm cache of device info
	userRepo, _ := domainRepo.GetUserRepository()
	companyInfo, err := userRepo.GetCompanyIdOfUser(
		ctx,
		&domainModel.GetCompanyIdOfUserInput{
			UserID: input.UserId,
		},
	)
	if err != nil {
		return nil
	}
	if companyInfo == nil {
		return nil
	}
	key := []string{
		sharedCache.GetKeyDeviceBase(sharedCrypto.GetHash(input.DeviceId.String())),
		sharedCache.GetKeyListDeviceInCompany(
			sharedCrypto.GetHash(companyInfo.CompanyID.String()),
			20,
			1,
		),
	}
	go func() {
		cacheService, _ := domainCache.GetDistributedCache()
		for _, k := range key {
			if err := cacheService.Delete(context.Background(), k); err != nil {
				global.Logger.Error("Error when delete device info cache", "err", err)
			}
		}
	}()
	return nil
}

// CreateNewDevice implements service.IDeviceService.
func (d *DeviceService) CreateNewDevice(ctx context.Context, input *model.CreateNewDeviceInput) (*model.CreateNewDeviceOutput, *applicationError.Error) {
	// Check permission
	domainUser, err := domainRepo.GetUserRepository()
	if err != nil {
		global.Logger.Error("Error when get user repository", "err", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if input.Role > 1 {
		return nil, &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to update device info.",
		}
	}
	// Get company id
	companyInfo, err := domainUser.GetCompanyIdOfUser(
		ctx,
		&domainModel.GetCompanyIdOfUserInput{
			UserID: input.UserId,
		},
	)
	if err != nil {
		global.Logger.Error("Error when get company id of user", "err", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if companyInfo == nil {
		return nil, &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "you don't have permission to create device.",
		}
	}
	// Create new device
	deviceUuid := uuid.New()
	deviceRepo, _ := domainRepo.GetDeviceRepository()
	deviceModel := &domainModel.NewDevice{
		DeviceId:     deviceUuid,
		CompanyId:    companyInfo.CompanyID,
		Name:         input.DeviceName,
		Address:      input.Address,
		SerialNumber: input.SerialNumber,
		MacAddress:   input.MacAddress,
		Token:        sharedRandom.RandomString(32),
	}
	if err := deviceRepo.CreateNewDevice(
		ctx,
		deviceModel,
	); err != nil {
		global.Logger.Error("Error when create new device", "err", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	// Save in cache
	cacheServie, _ := domainCache.GetDistributedCache()
	key := sharedCache.GetKeyDeviceBase(sharedCrypto.GetHash(deviceUuid.String()))
	ttl := constants.TTL_DEVICE_INFO
	if err := cacheServie.SetTTL(
		ctx,
		key,
		deviceModel,
		int64(ttl),
	); err != nil {
		global.Logger.Error("Error when set device info in cache", "err", err)
		// Not return error if cache error
	}
	return &model.CreateNewDeviceOutput{
		DeviceId:     deviceUuid.String(),
		CompanyId:    companyInfo.CompanyID.String(),
		Name:         input.DeviceName,
		Address:      input.Address,
		SerialNumber: input.SerialNumber,
		MacAddress:   input.MacAddress,
	}, nil
}

// DeleteDeviceById implements service.IDeviceService.
func (d *DeviceService) DeleteDeviceById(ctx context.Context, input *model.DeleteDeviceInput) *applicationError.Error {
	// Check permission
	domainUser, err := domainRepo.GetUserRepository()
	if err != nil {
		global.Logger.Error("Error when get user repository", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if input.Role > 1 {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to update device info.",
		}
	}
	ok, err := domainUser.UserPermissionDevice(ctx, &domainModel.UserPermissionDeviceInput{
		UserID:   input.UserId,
		DeviceID: input.DeviceId,
	})
	if err != nil {
		global.Logger.Error("Error when check user permission device", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !ok && input.Role != 0 {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to update device info.",
		}
	}
	// Check device exist
	deviceRepo, _ := domainRepo.GetDeviceRepository()
	exist, err := deviceRepo.DeviceExist(ctx, &domainModel.DeviceExistInput{
		DeviceId: input.DeviceId,
	})
	if err != nil {
		global.Logger.Error("Error when check device exist", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !exist {
		return &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "Device not found.",
		}
	}
	// Delete device
	if err := deviceRepo.DeleteDevice(
		ctx,
		&domainModel.DeleteDeviceInput{
			DeviceId: input.DeviceId,
		},
	); err != nil {
		global.Logger.Error("Error when delete device", "err", err)
		return &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	// Rm cache of device info
	userRepo, _ := domainRepo.GetUserRepository()
	companyInfo, err := userRepo.GetCompanyIdOfUser(
		ctx,
		&domainModel.GetCompanyIdOfUserInput{
			UserID: input.UserId,
		},
	)
	if err != nil {
		return nil
	}
	if companyInfo == nil {
		return nil
	}
	key := []string{
		sharedCache.GetKeyDeviceBase(sharedCrypto.GetHash(input.DeviceId.String())),
		sharedCache.GetKeyListDeviceInCompany(
			sharedCrypto.GetHash(companyInfo.CompanyID.String()),
			20,
			1,
		),
	}
	go func() {
		cacheService, _ := domainCache.GetDistributedCache()
		for _, k := range key {
			if err := cacheService.Delete(context.Background(), k); err != nil {
				global.Logger.Error("Error when delete device info cache", "err", err)
			}
		}
	}()
	return nil
}

// GetDeviceById implements service.IDeviceService.
func (d *DeviceService) GetDeviceById(ctx context.Context, input *model.GetDeviceByIdInput) (*model.GetDeviceByIdOutput, *applicationError.Error) {
	// Check user have permission to get device info
	if input.Role > 1 {
		return nil, &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to get device info.",
		}
	}
	// Get device info
	deviceRepo, _ := domainRepo.GetDeviceRepository()
	cacheService, _ := domainCache.GetDistributedCache()
	key := sharedCache.GetKeyDeviceBase(sharedCrypto.GetHash(input.DeviceId.String()))
	var deviceInfoCache domainModel.DeviceInfoBaseOutput
	deviceInfoCacheStr, err := cacheService.Get(ctx, key)
	if err != nil {
		global.Logger.Error("Error when get device info from cache", "err", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	var unMarshal bool = false
	if err := json.Unmarshal([]byte(deviceInfoCacheStr), &deviceInfoCache); err != nil {
		global.Logger.Error("Error when unmarshal device info from cache", "err", err)
		unMarshal = true
	}
	if unMarshal {
		deviceInfo, err := deviceRepo.DeviceInfoBase(
			ctx,
			&domainModel.DeviceInfoBaseInput{
				DeviceId: input.DeviceId,
			},
		)
		if err != nil {
			global.Logger.Error("Error when get device info", "err", err)
			return nil, &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "System is busy now. Please try again later.",
			}
		}
		deviceInfoCache = *deviceInfo
		// Save in cache
		ttl := constants.TTL_DEVICE_INFO
		if err := cacheService.SetTTL(
			ctx,
			key,
			deviceInfoCache,
			int64(ttl),
		); err != nil {
			global.Logger.Error("Error when set device info in cache", "err", err)
			// Not return error if cache error
		}
	}
	// Check user in company
	userRepo, _ := domainRepo.GetUserRepository()
	userInfo, err := userRepo.UserExistsInCompany(
		ctx,
		&domainModel.UserExistsInCompanyInput{
			UserID:    input.UserId,
			CompanyID: deviceInfoCache.CompanyId,
		},
	)
	if err != nil {
		global.Logger.Error("Error when check user in company", "err", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !userInfo && input.Role == domainModel.RoleManager {
		return nil, &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to get device info.",
		}
	}
	return &model.GetDeviceByIdOutput{
		DeviceId:     deviceInfoCache.DeviceId.String(),
		CompanyId:    deviceInfoCache.CompanyId.String(),
		Name:         deviceInfoCache.Name,
		Address:      deviceInfoCache.Address,
		SerialNumber: deviceInfoCache.SerialNumber,
		MacAddress:   deviceInfoCache.MacAddress,
		CreateAt:     deviceInfoCache.CreateAt,
		Token:        deviceInfoCache.Token,
	}, nil
}

// GetListDevices implements service.IDeviceService.
func (d *DeviceService) GetListDevices(ctx context.Context, input *model.ListDevicesInput) (*model.ListDevicesOutput, *applicationError.Error) {
	// Check user have permission to get device info
	if input.Role > 1 {
		return nil, &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to get device info.",
		}
	}
	// Get company if not input
	if input.CompanyId == uuid.Nil || input.CompanyId.String() == "" {
		userRepo, _ := domainRepo.GetUserRepository()
		companyInfo, err := userRepo.GetCompanyIdOfUser(
			ctx,
			&domainModel.GetCompanyIdOfUserInput{
				UserID: input.UserId,
			},
		)
		if err != nil {
			global.Logger.Error("Error when get company id of user", "err", err)
			return nil, &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "System is busy now. Please try again later.",
			}
		}
		if companyInfo == nil {
			return nil, &applicationError.Error{
				ErrorSystem: nil,
				ErrorClient: "you don't have permission to get device info.",
			}
		}
		input.CompanyId = companyInfo.CompanyID
	}
	// Get device info
	deviceRepo, _ := domainRepo.GetDeviceRepository()
	cacheService, _ := domainCache.GetDistributedCache()
	limit := input.Size
	offset := (input.Page - 1) * input.Size
	key := sharedCache.GetKeyListDeviceInCompany(sharedCrypto.GetHash(input.CompanyId.String()), limit, offset)
	var deviceInfoCache domainModel.ListDeviceInCompanyOutput
	deviceInfoCacheStr, err := cacheService.Get(ctx, key)
	if err != nil {
		global.Logger.Error("Error when get device info from cache", "err", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	var unMarshal bool = false
	if err := json.Unmarshal([]byte(deviceInfoCacheStr), &deviceInfoCache); err != nil {
		global.Logger.Error("Error when unmarshal device info from cache", "err", err)
		unMarshal = true
	}
	if len(deviceInfoCache.Devices) == 0 || unMarshal {
		deviceInfo, err := deviceRepo.ListDeviceInCompany(
			ctx,
			&domainModel.ListDeviceInCompanyInput{
				CompanyId: input.CompanyId,
				Limit:     limit,
				Offset:    offset,
			},
		)
		if err != nil {
			global.Logger.Error("Error when get device info", "err", err)
			return nil, &applicationError.Error{
				ErrorSystem: err,
				ErrorClient: "System is busy now. Please try again later.",
			}
		}
		if deviceInfo.Devices == nil {
			return &model.ListDevicesOutput{
				Devices: []*model.GetDeviceByIdOutput{},
			}, nil
		}
		deviceInfoCache = *deviceInfo
		// Save in cache
		ttl := constants.TTL_DEVICE_INFO
		if err := cacheService.SetTTL(
			ctx,
			key,
			deviceInfoCache,
			int64(ttl),
		); err != nil {
			global.Logger.Error("Error when set list device info in cache", "err", err)
			// Not return error if cache error
		}
	}
	// Check user in company
	userRepo, _ := domainRepo.GetUserRepository()
	userInfo, err := userRepo.UserExistsInCompany(
		ctx,
		&domainModel.UserExistsInCompanyInput{
			UserID:    input.UserId,
			CompanyID: input.CompanyId,
		},
	)
	if err != nil {
		global.Logger.Error("Error when check user in company", "err", err)
		return nil, &applicationError.Error{
			ErrorSystem: err,
			ErrorClient: "System is busy now. Please try again later.",
		}
	}
	if !userInfo && input.Role == domainModel.RoleManager {
		return nil, &applicationError.Error{
			ErrorSystem: nil,
			ErrorClient: "You don't have permission to get device info.",
		}
	}
	out := make([]*model.GetDeviceByIdOutput, 0)
	for _, device := range deviceInfoCache.Devices {
		out = append(out, &model.GetDeviceByIdOutput{
			DeviceId:     device.DeviceId.String(),
			CompanyId:    device.CompanyId.String(),
			Name:         device.Name,
			Address:      device.Address,
			SerialNumber: device.SerialNumber,
			MacAddress:   device.MacAddress,
			CreateAt:     device.CreateAt,
			Token:        device.Token,
			Status:       device.Status,
		})
	}
	return &model.ListDevicesOutput{
		Devices: out,
	}, nil
}

// UpdateDeviceById implements service.IDeviceService.
func (d *DeviceService) UpdateDeviceById(ctx context.Context, input *model.UpdateDeviceInput) (*model.UpdateDeviceOutput, *applicationError.Error) {
	panic("unimplemented")
}

// NewDeviceService create new instance and implement IDeviceService
func NewDeviceService() service.IDeviceService {
	return &DeviceService{}
}
