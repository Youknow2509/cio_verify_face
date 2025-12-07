package mail

import (
	"html"
	"strconv"
	"time"

	domainMail "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mail"
)

// Struct for HTML content mail
type HtmlMailContent struct {
}

// ForgotPassword implements mail.IHtmlMailContent.
func (h *HtmlMailContent) ForgotPassword(to string, url_auth string, new_password string, expired int64) (string, error) {
	minuteExpired := time.Unix(expired, 0).Minute()
	minuteExpiredStr := strconv.Itoa(minuteExpired)
	return `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Forgot Password</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 0;
				}
				.container {
					max-width: 600px;
					margin: 50px auto;
					background-color: #ffffff;
					padding: 20px;
					border-radius: 5px;
					box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
				}
				h1 {
					color: #333333;
				}
				p {
					color: #666666;
					line-height: 1.6;
				}
				a {
					color: #1a73e8;
					text-decoration: none;
				}
				.button {
					display: inline-block;
					padding: 10px 20px;
					margin-top: 20px;
					background-color: #1a73e8;
					color: #ffffff;
					border-radius: 5px;
					text-decoration: none;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Password Reset Request</h1>
				<p>Dear User,</p>
				<p>We received a request to reset your password. Please click the button below to reset your password:</p>
				<a href="` + url_auth + `" class="button">Reset Password</a>
				<p>Your new password is: <strong>` + new_password + `</strong></p>
				<p>This link will expire in ` + minuteExpiredStr + ` minutes.</p>
				<p>If you did not request a password reset, please ignore this email.</p>
				<p>Best regards,<br>Your Company Team</p>
			</div>
		</body>
		</html>
	`, nil
}

// PasswordResetNotification implements mail.IHtmlMailContent.
func (h *HtmlMailContent) PasswordResetNotification(to string, fullName string, resetURL string, expiresIn int) (string, error) {
	// Escape HTML to prevent XSS
	escapedFullName := html.EscapeString(fullName)
	escapedResetURL := html.EscapeString(resetURL)
	
	return `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Password Reset</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 0;
				}
				.container {
					max-width: 600px;
					margin: 50px auto;
					background-color: #ffffff;
					padding: 20px;
					border-radius: 5px;
					box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
				}
				h1 {
					color: #333333;
				}
				p {
					color: #666666;
					line-height: 1.6;
				}
				.button {
					display: inline-block;
					padding: 10px 20px;
					margin-top: 20px;
					background-color: #1a73e8;
					color: #ffffff;
					border-radius: 5px;
					text-decoration: none;
				}
				.warning {
					color: #d93025;
					font-weight: bold;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Password Reset Request</h1>
				<p>Dear ` + escapedFullName + `,</p>
				<p>We received a request to reset your password. Click the button below to reset it:</p>
				<a href="` + escapedResetURL + `" class="button">Reset Password</a>
				<p>This link will expire in <strong>` + strconv.Itoa(expiresIn) + ` hours</strong>.</p>
				<p class="warning">If you did not request this password reset, please ignore this email or contact support if you have concerns.</p>
				<p>Best regards,<br>Your CIO Verify Face Team</p>
			</div>
		</body>
		</html>
	`, nil
}

// ReportAttentionNotification implements mail.IHtmlMailContent.
func (h *HtmlMailContent) ReportAttentionNotification(email string, downloadURL string, reportType string, format string, startDate string, endDate string) (string, error) {
	// Escape HTML to prevent XSS
	escapedDownloadURL := html.EscapeString(downloadURL)
	escapedReportType := html.EscapeString(reportType)
	escapedFormat := html.EscapeString(format)
	escapedStartDate := html.EscapeString(startDate)
	escapedEndDate := html.EscapeString(endDate)
	
	return `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Report Ready</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					margin: 0;
					padding: 0;
				}
				.container {
					max-width: 600px;
					margin: 50px auto;
					background-color: #ffffff;
					padding: 20px;
					border-radius: 5px;
					box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
				}
				h1 {
					color: #333333;
				}
				p {
					color: #666666;
					line-height: 1.6;
				}
				.button {
					display: inline-block;
					padding: 10px 20px;
					margin-top: 20px;
					background-color: #34a853;
					color: #ffffff;
					border-radius: 5px;
					text-decoration: none;
				}
				.info {
					background-color: #f0f0f0;
					padding: 15px;
					border-radius: 5px;
					margin: 15px 0;
				}
				.info-row {
					margin: 5px 0;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Your Report is Ready</h1>
				<p>Dear User,</p>
				<p>Your requested report has been generated and is ready for download.</p>
				<div class="info">
					<div class="info-row"><strong>Report Type:</strong> ` + escapedReportType + `</div>
					<div class="info-row"><strong>Format:</strong> ` + escapedFormat + `</div>
					<div class="info-row"><strong>Period:</strong> ` + escapedStartDate + ` to ` + escapedEndDate + `</div>
				</div>
				<a href="` + escapedDownloadURL + `" class="button">Download Report</a>
				<p>This download link will expire in 7 days.</p>
				<p>Best regards,<br>Your CIO Verify Face Team</p>
			</div>
		</body>
		</html>
	`, nil
}

// New HTMLContentMail and impl IHtmlMailContent
func NewHTMLContentMail() domainMail.IHtmlMailContent {
	return &HtmlMailContent{}
}
