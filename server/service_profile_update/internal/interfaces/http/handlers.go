package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/interfaces/http/dto"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/interfaces/http/mapper"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/interfaces/http/middleware"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/interfaces/response"
)

// FaceProfileUpdateHandler handles face profile update requests
type FaceProfileUpdateHandler struct{}

// NewFaceProfileUpdateHandler creates a new handler
func NewFaceProfileUpdateHandler() *FaceProfileUpdateHandler {
	return &FaceProfileUpdateHandler{}
}

// CreateRequest handles POST /api/v1/profile-update/requests
// @Summary Create a face profile update request
// @Description Employee creates a request to update their face profile
// @Tags Face Profile Update
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateFaceProfileUpdateRequestDTO true "Request body"
// @Success 200 {object} response.Response{data=model.CreateFaceProfileUpdateRequestOutput}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 429 {object} response.Response
// @Router /api/v1/profile-update/requests [post]
func (h *FaceProfileUpdateHandler) CreateRequest(c *gin.Context) {
	// Get session from context (set by auth middleware)
	session := getSessionFromContext(c)
	if session == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	// Bind and validate DTO
	var requestDTO dto.CreateFaceProfileUpdateRequestDTO
	if !middleware.BindAndValidate(c, &requestDTO, "json") {
		return
	}

	svc := service.GetFaceProfileUpdateService()
	if svc == nil {
		response.InternalError(c, "Service unavailable")
		return
	}

	// Map DTO to application input
	input := mapper.ToCreateFaceProfileUpdateRequestInput(&requestDTO, session)
	result, appErr := svc.CreateUpdateRequest(c.Request.Context(), input)

	if appErr != nil {
		response.FromAppError(c, appErr)
		return
	}

	response.Success(c, result)
}

// GetMyRequests handles GET /api/v1/profile-update/requests/me
// @Summary Get my face profile update requests
// @Description Employee gets their own update requests
// @Tags Face Profile Update
// @Produce json
// @Security BearerAuth
// @Param month query string false "Month filter (YYYY-MM)"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response{data=model.GetMyUpdateRequestsOutput}
// @Failure 401 {object} response.Response
// @Router /api/v1/profile-update/requests/me [get]
func (h *FaceProfileUpdateHandler) GetMyRequests(c *gin.Context) {
	session := getSessionFromContext(c)
	if session == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	// Bind and validate query parameters
	var queryDTO dto.GetMyRequestsQueryDTO
	if !middleware.BindAndValidate(c, &queryDTO, "query") {
		return
	}

	// Set defaults
	queryDTO.SetDefaults()

	svc := service.GetFaceProfileUpdateService()
	if svc == nil {
		response.InternalError(c, "Service unavailable")
		return
	}

	// Map DTO to application input
	input := mapper.ToGetMyUpdateRequestsInput(&queryDTO, session)
	result, appErr := svc.GetMyUpdateRequests(c.Request.Context(), input)
	if appErr != nil {
		response.FromAppError(c, appErr)
		return
	}

	response.Success(c, result)
}

// GetPendingRequests handles GET /api/v1/profile-update/requests/pending
// @Summary Get pending face profile update requests
// @Description Manager gets pending requests for their company
// @Tags Face Profile Update
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response{data=model.GetPendingRequestsOutput}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/profile-update/requests/pending [get]
func (h *FaceProfileUpdateHandler) GetPendingRequests(c *gin.Context) {
	session := getSessionFromContext(c)
	if session == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	// Bind and validate query parameters
	var queryDTO dto.GetPendingRequestsQueryDTO
	if !middleware.BindAndValidate(c, &queryDTO, "query") {
		return
	}

	// Set defaults
	queryDTO.SetDefaults()

	svc := service.GetFaceProfileUpdateService()
	if svc == nil {
		response.InternalError(c, "Service unavailable")
		return
	}

	// Map DTO to application input
	input := mapper.ToGetPendingRequestsInput(&queryDTO, session)
	result, appErr := svc.GetPendingRequests(c.Request.Context(), input)
	if appErr != nil {
		response.FromAppError(c, appErr)
		return
	}

	response.Success(c, result)
}

// ApproveRequest handles POST /api/v1/profile-update/requests/:id/approve
// @Summary Approve a face profile update request
// @Description Manager approves a pending request
// @Tags Face Profile Update
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} response.Response{data=model.ApproveRequestOutput}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/profile-update/requests/{id}/approve [post]
func (h *FaceProfileUpdateHandler) ApproveRequest(c *gin.Context) {
	session := getSessionFromContext(c)
	if session == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	// Bind and validate URI parameter
	var paramDTO dto.ApproveRequestParamDTO
	if !middleware.BindAndValidate(c, &paramDTO, "uri") {
		return
	}

	svc := service.GetFaceProfileUpdateService()
	if svc == nil {
		response.InternalError(c, "Service unavailable")
		return
	}

	// Map DTO to application input
	input := mapper.ToApproveRequestInput(&paramDTO, session)
	result, appErr := svc.ApproveRequest(c.Request.Context(), input)

	if appErr != nil {
		response.FromAppError(c, appErr)
		return
	}

	response.Success(c, result)
}

// RejectRequest handles POST /api/v1/profile-update/requests/:id/reject
// @Summary Reject a face profile update request
// @Description Manager rejects a pending request
// @Tags Face Profile Update
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Param request body dto.RejectRequestBodyDTO true "Rejection reason"
// @Success 200 {object} response.Response{data=model.RejectRequestOutput}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/profile-update/requests/{id}/reject [post]
func (h *FaceProfileUpdateHandler) RejectRequest(c *gin.Context) {
	session := getSessionFromContext(c)
	if session == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	// Bind and validate URI parameter
	var paramDTO dto.RejectRequestParamDTO
	if !middleware.BindAndValidate(c, &paramDTO, "uri") {
		return
	}

	// Bind and validate request body (reason is optional)
	var bodyDTO dto.RejectRequestBodyDTO
	_ = c.ShouldBindJSON(&bodyDTO) // Ignore error as reason is optional

	svc := service.GetFaceProfileUpdateService()
	if svc == nil {
		response.InternalError(c, "Service unavailable")
		return
	}

	// Map DTO to application input
	input := mapper.ToRejectRequestInput(&paramDTO, &bodyDTO, session)
	result, appErr := svc.RejectRequest(c.Request.Context(), input)

	if appErr != nil {
		response.FromAppError(c, appErr)
		return
	}

	response.Success(c, result)
}

// ValidateToken handles GET /api/v1/profile-update/token/validate
// @Summary Validate an update token
// @Description Validate if an update token is valid
// @Tags Face Profile Update
// @Produce json
// @Param token query string true "Update token"
// @Success 200 {object} response.Response{data=model.ValidateUpdateTokenOutput}
// @Failure 400 {object} response.Response
// @Router /api/v1/profile-update/token/validate [get]
func (h *FaceProfileUpdateHandler) ValidateToken(c *gin.Context) {
	// Bind and validate query parameter
	var queryDTO dto.ValidateTokenQueryDTO
	if !middleware.BindAndValidate(c, &queryDTO, "query") {
		return
	}

	svc := service.GetFaceProfileUpdateService()
	if svc == nil {
		response.InternalError(c, "Service unavailable")
		return
	}

	// Map DTO to application input
	input := mapper.ToValidateUpdateTokenInput(&queryDTO)
	result, appErr := svc.ValidateUpdateToken(c.Request.Context(), input)

	if appErr != nil {
		response.FromAppError(c, appErr)
		return
	}

	response.Success(c, result)
}

// UpdateFaceProfile handles POST /api/v1/profile-update/face
// @Summary Update face profile
// @Description Update face profile using a valid update token
// @Tags Face Profile Update
// @Accept multipart/form-data
// @Produce json
// @Param token formData string true "Update token"
// @Param image formData file true "Face image"
// @Success 200 {object} response.Response{data=model.UpdateFaceProfileOutput}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/profile-update/face [post]
func (h *FaceProfileUpdateHandler) UpdateFaceProfile(c *gin.Context) {
	// Bind and validate form data
	var formDTO dto.UpdateFaceProfileFormDTO
	if !middleware.BindAndValidate(c, &formDTO, "form") {
		return
	}

	// Get and validate image file
	file, err := c.FormFile("image")
	if err != nil {
		response.BadRequest(c, "Image file is required")
		return
	}

	// Validate file size (e.g., max 10MB)
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if file.Size > maxFileSize {
		response.BadRequest(c, "Image file too large (max 10MB)")
		return
	}

	// Read file content
	f, err := file.Open()
	if err != nil {
		response.BadRequest(c, "Failed to read image file")
		return
	}
	defer f.Close()

	imageData, err := io.ReadAll(f)
	if err != nil {
		response.BadRequest(c, "Failed to read image data")
		return
	}

	svc := service.GetFaceProfileUpdateService()
	if svc == nil {
		response.InternalError(c, "Service unavailable")
		return
	}

	// Map DTO to application input
	input := mapper.ToUpdateFaceProfileInput(&formDTO, imageData, file.Filename)
	result, appErr := svc.UpdateFaceProfile(c.Request.Context(), input)

	if appErr != nil {
		response.FromAppError(c, appErr)
		return
	}

	response.Success(c, result)
}

// PasswordResetHandler handles password reset requests
type PasswordResetHandler struct{}

// NewPasswordResetHandler creates a new handler
func NewPasswordResetHandler() *PasswordResetHandler {
	return &PasswordResetHandler{}
}

// ResetEmployeePassword handles POST /api/v1/password/reset
// @Summary Reset employee password
// @Description Manager resets an employee's password
// @Tags Password Reset
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ResetPasswordRequestDTO true "Request body"
// @Success 200 {object} response.Response{data=model.ResetEmployeePasswordOutput}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 429 {object} response.Response
// @Router /api/v1/password/reset [post]
func (h *PasswordResetHandler) ResetEmployeePassword(c *gin.Context) {
	session := getSessionFromContext(c)
	if session == nil {
		response.Unauthorized(c, "Authentication required")
		return
	}

	// Bind and validate DTO
	var requestDTO dto.ResetPasswordRequestDTO
	if !middleware.BindAndValidate(c, &requestDTO, "json") {
		return
	}

	svc := service.GetPasswordResetService()
	if svc == nil {
		response.InternalError(c, "Service unavailable")
		return
	}

	// Map DTO to application input
	input := mapper.ToResetEmployeePasswordInput(&requestDTO, session)
	result, appErr := svc.ResetEmployeePassword(c.Request.Context(), input)

	if appErr != nil {
		response.FromAppError(c, appErr)
		return
	}

	response.Success(c, result)
}

// Helper functions
func getSessionFromContext(c *gin.Context) *model.SessionInfo {
	session, exists := c.Get("session")
	if !exists {
		return nil
	}
	if s, ok := session.(*model.SessionInfo); ok {
		return s
	}
	return nil
}

// SetupRoutes configures all HTTP routes
func SetupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "profile_update",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Face profile update routes
		faceHandler := NewFaceProfileUpdateHandler()
		profileUpdate := v1.Group("/profile-update")
		{
			// Employee endpoints (require authentication)
			profileUpdate.POST("/requests", middleware.AuthMiddleware(), faceHandler.CreateRequest)
			profileUpdate.GET("/requests/me", middleware.AuthMiddleware(), faceHandler.GetMyRequests)

			// Manager endpoints (require authentication)
			profileUpdate.GET("/requests/pending", middleware.AuthMiddleware(), faceHandler.GetPendingRequests)
			profileUpdate.POST("/requests/:id/approve", middleware.AuthMiddleware(), faceHandler.ApproveRequest)
			profileUpdate.POST("/requests/:id/reject", middleware.AuthMiddleware(), faceHandler.RejectRequest)

			// Token validation (public - no auth required)
			profileUpdate.GET("/token/validate", faceHandler.ValidateToken)

			// Face update (requires valid update token, not user token)
			profileUpdate.POST("/face", faceHandler.UpdateFaceProfile)
		}

		// Password reset routes
		passwordHandler := NewPasswordResetHandler()
		password := v1.Group("/password")
		{
			// Require authentication for password reset
			password.POST("/reset", middleware.AuthMiddleware(), passwordHandler.ResetEmployeePassword)
		}
	}
}
