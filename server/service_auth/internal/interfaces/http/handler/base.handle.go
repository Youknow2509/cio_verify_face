package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	constants "github.com/youknow2509/cio_verify_face/server/service_auth/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/interfaces/dto"
	interfaceResponse "github.com/youknow2509/cio_verify_face/server/service_auth/internal/interfaces/response"
	utilsContext "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/context"
	utilsUuid "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/uuid"
)

/**
 * Auth handler
 */
type AuthBaseHandler struct {
}

/**
 * GetAuthBaseHandler creates a Get instance of AuthBaseHandler
 */
func GetAuthBaseHandler() *AuthBaseHandler {
	return &AuthBaseHandler{}
}

// Handle user login for user
// @Summary      User login
// @Description  User login for user
// @Tags         Core Auth
// @Accept       json
// @Produce      json
// @Param        request   body dto.LoginRequest  true  "Request body login"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/auth/login [post]
func (h *AuthBaseHandler) Login(c *gin.Context) {
	// Bind the request to the LoginRequest DTO
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid request parameters",
		)
		return
	}
	// Validate the request
	validate := c.MustGet(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME).(*validator.Validate)
	if err := validate.Struct(request); err != nil {
		var fieldErrors []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, fieldError.Field())
		}
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid request parameters: "+strings.Join(fieldErrors, ", "),
		)
		return
	}
	// Call handle to service
	response, err := applicationService.GetCoreAuthService().Login(
		c,
		&applicationModel.LoginInput{
			UserName:  request.UserName,
			Password:  request.Password,
			ClientIp:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		},
	)
	if err != nil {
		interfaceResponse.ErrorResponse(
			c,
			err.Code,
			err.Message,
		)
		return
	}
	interfaceResponse.SuccessResponse(
		c,
		interfaceResponse.ErrCodeSuccess,
		response,
	)
}

// Handle user login for admin
// @Summary      Admin login
// @Description  User login for admin
// @Tags         Core Auth
// @Accept       json
// @Produce      json
// @Param        request   body dto.LoginRequest  true  "Request body login"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/auth/login/admin [post]
func (h *AuthBaseHandler) LoginAdmin(c *gin.Context) {
	// Bind the request to the LoginRequest DTO
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid request parameters",
		)
		return
	}
	// Validate the request
	validate := c.MustGet(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME).(*validator.Validate)
	if err := validate.Struct(request); err != nil {
		var fieldErrors []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, fieldError.Field())
		}
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid request parameters: "+strings.Join(fieldErrors, ", "),
		)
		return
	}
	// Call handle to service
	response, err := applicationService.GetCoreAuthService().LoginAdmin(
		c,
		&applicationModel.LoginInputAdmin{
			UserName: request.UserName,
			Password: request.Password,
			ClientIp: c.ClientIP(),
		},
	)
	if err != nil {
		interfaceResponse.ErrorResponse(
			c,
			err.Code,
			err.Message,
		)
		return
	}
	interfaceResponse.SuccessResponse(
		c,
		interfaceResponse.ErrCodeSuccess,
		response,
	)
}

// Handle user logout
// @Summary      User logout
// @Description  User logout
// @Tags         Core Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authorization Bearer token"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/auth/logout [post]
func (h *AuthBaseHandler) Logout(c *gin.Context) {
	// Get data auth from context
	userIdStr, sessionIdStr, _, exists := utilsContext.GetSessionFromContext(c)
	if !exists {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid request parameters",
		)
		return
	}
	// Validate id str to uuid
	userId, err := utilsUuid.ParseUUID(userIdStr)
	if err != nil {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid data session",
		)
		return
	}
	sessionId, err := utilsUuid.ParseUUID(sessionIdStr)
	if err != nil {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid data session",
		)
		return
	}
	// Call handle to service
	if err := applicationService.GetCoreAuthService().Logout(
		c,
		&applicationModel.LogoutInput{
			UserId:    userId,
			SessionId: sessionId,
		},
	); err != nil {
		interfaceResponse.ErrorResponse(
			c,
			err.Code,
			err.Message,
		)
		return
	}
	interfaceResponse.SuccessResponse(
		c,
		interfaceResponse.ErrCodeSuccess,
		nil,
	)
}

// Handle refresh token
// @Summary      User refresh token
// @Description  User refresh token
// @Tags         Core Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authorization Bearer token"
// @Param        request   body dto.RefreshTokenRequest  true  "Request body refresh token"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/auth/refresh [post]
func (h *AuthBaseHandler) RefreshToken(c *gin.Context) {
	// Get data auth from context
	userIdStr, sessionIdStr, userRole, exists := utilsContext.GetSessionFromContext(c)
	if !exists {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid request parameters",
		)
		return
	}
	// Validate id str to uuid
	userId, err := utilsUuid.ParseUUID(userIdStr)
	if err != nil {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid data session",
		)
		return
	}
	sessionId, err := utilsUuid.ParseUUID(sessionIdStr)
	if err != nil {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid data session",
		)
		return
	}
	// Bind the request to the RefreshTokenRequest DTO
	var request dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid request parameters",
		)
		return
	}
	// Validate the request
	validate := c.MustGet(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME).(*validator.Validate)
	if err := validate.Struct(request); err != nil {
		var fieldErrors []string
		for _, fieldError := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, fieldError.Field())
		}
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid request parameters: "+strings.Join(fieldErrors, ", "),
		)
		return
	}
	// Call handle to service
	response, serviceErr := applicationService.GetCoreAuthService().RefreshToken(
		c,
		&applicationModel.RefreshTokenInput{
			UserId:       userId,
			SessionId:    sessionId,
			RefreshToken: request.RefreshToken,
			ClientIp:     c.ClientIP(),
			UserRole:     userRole,
		},
	)
	if serviceErr != nil {
		interfaceResponse.ErrorResponse(
			c,
			serviceErr.Code,
			serviceErr.Message,
		)
		return
	}
	interfaceResponse.SuccessResponse(
		c,
		interfaceResponse.ErrCodeSuccess,
		response,
	)
}

// User get info
// @Summary      User get base info
// @Description  User get base info
// @Tags         Core Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authorization Bearer token"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/auth/me [post]
func (h *AuthBaseHandler) GetMyInfo(c *gin.Context) {
	// Get data auth from context
	userIdStr, _, role, exists := utilsContext.GetSessionFromContext(c)
	if !exists {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid request parameters",
		)
		return
	}
	// Validate id str to uuid
	userId, err := utilsUuid.ParseUUID(userIdStr)
	if err != nil {
		interfaceResponse.BadRequestResponse(
			c,
			interfaceResponse.ErrCodeParamInvalid,
			"Invalid data session",
		)
		return
	}
	// Call handle to service
	response, err_r := applicationService.GetCoreAuthService().GetMyInfo(
		c,
		&applicationModel.GetMyInfoInput{
			ClientIp: c.ClientIP(),
			UserId:   userId,
			Role:     role,
		},
	)
	if err_r != nil {
		interfaceResponse.ErrorResponse(
			c,
			err_r.Code,
			err_r.Message,
		)
		return
	}
	interfaceResponse.SuccessResponse(
		c,
		interfaceResponse.ErrCodeSuccess,
		response,
	)
}

// CreateDevice create device
// @Summary      Create session device
// @Description  Create session device
// @Tags         Core Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authorization Bearer token"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/auth/device [post]
func (h *AuthBaseHandler) CreateDevice(c *gin.Context) {
	// TODO: Implement create device handler
}

// Delete device by id
// @Summary      Delete session device
// @Description  Delete session device
// @Tags         Core Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authorization Bearer token"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/auth/device/{id} [delete]
func (h *AuthBaseHandler) DeleteDevice(c *gin.Context) {
	// TODO: Implement delete device handler
}

// RefreshTokenDevice refresh device token
// @Summary      Refresh session device
// @Description  Refresh session device
// @Tags         Core Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authorization Bearer token"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/device/{id}/refresh [post]
func (h *AuthBaseHandler) RefreshTokenDevice(c *gin.Context) {
	// TODO: Implement RefreshTokenDevice device handler
}
