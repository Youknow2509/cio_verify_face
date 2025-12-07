package dto

// ResetPasswordRequestDTO represents the request body for resetting employee password
type ResetPasswordRequestDTO struct {
	EmployeeID string `json:"employee_id" binding:"required,min=1,max=100" example:"EMP001"`
}

// ConfirmPasswordResetDTO represents the request body for confirming password reset
type ConfirmPasswordResetDTO struct {
	Token string `json:"token" binding:"required" example:"abc123..."`
}
