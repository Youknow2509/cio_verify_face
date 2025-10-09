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
