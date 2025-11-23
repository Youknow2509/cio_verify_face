package cache

import "fmt"

// =================================
//
//	Define key cache
//
// =================================

// Key get status token device
func GetKeyStatusTokenDevice(tokenHash string) string {
	return fmt.Sprintf("device:token:status:%s", tokenHash)
}

// Key get status token user
func GetKeyStatusTokenUser(tokenHash string) string {
	return fmt.Sprintf("user:token:status:%s", tokenHash)
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
// 			Define value cache
// =================================
