package model

/**
 * Mail model
 */
type MailForgotPassword struct {
	To          string `json:"to" validate:"required,email"`
	UrlAuth     string `json:"url_auth" validate:"required,url"`
	NewPassword string `json:"new_password" validate:"required"`
	Expired     int64  `json:"expired" validate:"required,gt=0"`
}

/**
 * Password Reset Notification model
 */
type PasswordResetNotification struct {
	To        string `json:"to" validate:"required,email"`
	FullName  string `json:"full_name" validate:"required"`
	ResetURL  string `json:"reset_url" validate:"required,url"`
	ExpiresIn int    `json:"expires_in" validate:"required,gt=0"`
}

/**
 * Report Attention Notification model
 */
type ReportAttentionNotification struct {
	Email       string `json:"email" validate:"required,email"`
	CompanyID   string `json:"company_id" validate:"required"`
	DownloadURL string `json:"download_url" validate:"required,url"`
	Type        string `json:"type" validate:"required"`
	Format      string `json:"format" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date" validate:"required"`
	CreatedAt   string `json:"created_at" validate:"required"`
}
