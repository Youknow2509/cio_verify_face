package cache

import "fmt"

// =================================
//
//	Define key cache
//
// =================================

// Key get attendance records employee
func GetKeyAttendanceRecordsEmployee(employeeIdHash string, yearMonth string, pageSize int, pageStage []byte) string {
	if pageStage == nil {
		return fmt.Sprintf("employee:attendance:records:%s:%s", employeeIdHash, yearMonth)
	} else {
		return fmt.Sprintf("employee:attendance:records:%s:%s:%d:%s", employeeIdHash, yearMonth, pageSize, string(pageStage))
	}
}

// Key get attendance records company
func GetKeyAttendanceRecordsCompany(companyIdHash string, yearMonth string, pageSize int, pageStage []byte) string {
	if pageStage == nil {
		return fmt.Sprintf("company:attendance:records:%s:%s", companyIdHash, yearMonth)
	} else {
		return fmt.Sprintf("company:attendance:records:%s:%s:%d:%s", companyIdHash, yearMonth, pageSize, string(pageStage))
	}
}

// Key get list shift time employee
func GetKeyListShiftTimeEmployee(employeeIdHash string) string {
	return fmt.Sprintf("employee:shift:time:list:%s", employeeIdHash)
}

// //////////////////////
// Key GetKeyUserIsManagerCompanyForUser
func GetKeyUserIsManagerCompanyForUser(userIdReqHash string, userIdHash string) string {
	return fmt.Sprintf("user:is:manager:company:%s:%s", userIdReqHash, userIdHash)
}

// Key GetKeyUserIsManagerCompany
func GetKeyUserIsManagerCompany(userIdHash string, companyIdHash string) string {
	return fmt.Sprintf("user:is:manager:company:%s:%s", userIdHash, companyIdHash)
}

// Key user register OTP value
func GetKeyUserRegisterOTP(mailHash string) string {
	return fmt.Sprintf("user:register:otp:%s", mailHash)
}

// Key user register OTP value count
func GetKeyUserRegisterOTPCount(mailHash string) string {
	return fmt.Sprintf("user:register:otp:count:%s", mailHash)
}

// Key block spam register
func GetKeyBlockSpamRegister(mailHash string) string {
	return fmt.Sprintf("user:register:spam:block:%s", mailHash)
}

// Key token user verify register
func GetKeyUserVerifyRegisterToken(mailHash string) string {
	return fmt.Sprintf("user:register:verify:token:%s", mailHash)
}

// Key user access token is active
func GetKeyUserAccessTokenIsActive(sessionIdHash string) string {
	return fmt.Sprintf("user:access:token:is:active:%s", sessionIdHash)
}

// Key check spam reset password
func GetKeyCheckSpamResetPassword(mailHash string) string {
	return fmt.Sprintf("user:reset:password:spam:%s", mailHash)
}

// key count spam reset password
func GetKeyCountSpamResetPassword(mailHash string) string {
	return fmt.Sprintf("user:reset:password:spam:count:%s", mailHash)
}

// Key token user reset password
func GetKeyUserResetPasswordToken(userIdHash string) string {
	return fmt.Sprintf("user:reset:password:token:%s", userIdHash)
}

// Get Key count verify register fail
func GetKeyCountVerifyRegisterFail(mailHash string) string {
	return fmt.Sprintf("user:register:verify:fail:count:%s", mailHash)
}

// Get Key block spam refresh token
func GetKeyBlockSpamRefreshToken(userIdHash string) string {
	return fmt.Sprintf("user:refresh:token:spam:block:%s", userIdHash)
}

// Key block spam validate reset password token
func GetKeyBlockSpamValidateResetPasswordToken(userIdHash string) string {
	return fmt.Sprintf("user:reset:password:token:spam:block:%s", userIdHash)
}

// Key password user after verify reset password token
func GetKeyPasswordUserHashAfterVerifyResetPasswordToken(userIdHash string) string {
	return fmt.Sprintf("user:reset:password:token:password:%s", userIdHash)
}

// Key friend request cache
func GetKeyCacheFriendRequest(friendRequestIdHash string) string {
	return fmt.Sprintf("user:friends:request:cache:%s", friendRequestIdHash)
}

// Key user info view cache
func GetKeyCacheUserInfoView(userIdHash string) string {
	return fmt.Sprintf("user:info:view:cache:%s", userIdHash)
}

// Key user info view cache with email hash
func GetKeyCacheUserInfoViewWithEmailHash(emailHash string) string {
	return fmt.Sprintf("user:info:view:cache:email:%s", emailHash)
}

// Key cache list friend request to user
func GetKeyCacheListFriendRequestToUser(userIdHash string, page int) string {
	return fmt.Sprintf("user:friends:request:list:to:user:%s:%d", userIdHash, page)
}

// Key cache list friend request from user
func GetKeyCacheListFriendsOfUser(userIdHash string, page int) string {
	return fmt.Sprintf("user:friends:list:of:user:%s:%d", userIdHash, page)
}

// =================================
// Attendance related keys
// =================================

// Key for user's last check-in on a given date. Useful to prevent duplicate check-ins.
func GetKeyAttendanceUserLastCheckIn(userIdHash string, date string) string {
	return fmt.Sprintf("attendance:user:last_checkin:%s:%s", userIdHash, date)
}

// Key for device attendance records cache per company and date
func GetKeyAttendanceDeviceRecords(companyIdHash string, deviceIdHash string, page int, size int, date string) string {
	return fmt.Sprintf("attendance:device:records:%s:%s:%d:%d:%s", companyIdHash, deviceIdHash, page, size, date)
}

// Key get daily attendance summary employee
func GetKeyDailyAttendanceSummaryEmployee(employeeIdHash string, summaryMonth string, pageSize int, pageStage []byte) string {
	if pageStage == nil {
		return fmt.Sprintf("employee:attendance:summary:%s:%s", employeeIdHash, summaryMonth)
	} else {
		return fmt.Sprintf("employee:attendance:summary:%s:%s:%d:%s", employeeIdHash, summaryMonth, pageSize, string(pageStage))
	}
}

// Key get daily attendance summary for company
func GetKeyDailyAttendanceSummary(companyIdHash string, summaryMonth string, workDate int64, pageSize int, pageStage []byte) string {
	if pageStage == nil {
		return fmt.Sprintf("company:attendance:summary:%s:%s:%d", companyIdHash, summaryMonth, workDate)
	} else {
		return fmt.Sprintf("company:attendance:summary:%s:%s:%d:%d:%s", companyIdHash, summaryMonth, workDate, pageSize, string(pageStage))
	}
}
