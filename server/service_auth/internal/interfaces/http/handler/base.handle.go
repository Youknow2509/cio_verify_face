package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/interfaces/dto"
	constants "github.com/youknow2509/cio_verify_face/server/service_auth/internal/constants"
	interfaceResponse "github.com/youknow2509/cio_verify_face/server/service_auth/internal/interfaces/response"
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
// @Param        request   body dto.LogoutRequest  true  "Request body logout"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/auth/logout [post]
func (h *AuthBaseHandler) Logout(c *gin.Context) {
	// TODO: Implement logout handler
}

// Handle refresh token
// @Summary      User refresh token
// @Description  User refresh token
// @Tags         Core Auth
// @Accept       json
// @Produce      json
// @Param        request   body dto.RefreshTokenRequest  true  "Request body refresh token"
// @Success      200  {object}  dto.ResponseData
// @Failure      400  {object}  dto.ErrResponseData
// @Router       /v1/auth/refresh-token [post]
func (h *AuthBaseHandler) RefreshToken(c *gin.Context) {
	// TODO: Implement refresh token handler
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
// @Router       /v1/auth/refresh-token [post]
func (h *AuthBaseHandler) GetMyInfo(c *gin.Context) {
	// TODO: Implement get my info handler
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