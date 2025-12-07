package mail

import "errors"

/**
 * Interface create html mail content
 */
type IHtmlMailContent interface {
	ForgotPassword(
		to string,
		url_auth string,
		new_password string,
		expired int64,
	) (string, error)
	PasswordResetNotification(
		to string,
		fullName string,
		resetURL string,
		expiresIn int,
	) (string, error)
	ReportAttentionNotification(
		email string,
		downloadURL string,
		reportType string,
		format string,
		startDate string,
		endDate string,
	) (string, error)
}

/**
 * Manage instance of HtmlMailContent
 */
var _vIHtmlMailContent IHtmlMailContent

func SetHtmlMailContent(v IHtmlMailContent) error {
	if v == nil {
		return errors.New("HtmlMailContent is nil")
	}
	if _vIHtmlMailContent != nil {
		return errors.New("HtmlMailContent already set")
	}
	_vIHtmlMailContent = v
	return nil
}

func GetHtmlMailContent() IHtmlMailContent {
	return _vIHtmlMailContent
}
