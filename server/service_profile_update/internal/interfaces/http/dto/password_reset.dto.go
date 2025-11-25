package dto

// ResetPasswordRequestDTO represents the request body for resetting employee password
type ResetPasswordRequestDTO struct {
	EmployeeID string `json:"employee_id" binding:"required,min=1,max=100" example:"EMP001"`
}
