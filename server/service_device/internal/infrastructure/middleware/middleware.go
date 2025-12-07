package middleware

/**
 * Get auth device access token JWT middleware instance
 */
func GetAuthDeviceAccessTokenJwtMiddleware() *AuthDeviceTokenJwtMiddleware {
	return &AuthDeviceTokenJwtMiddleware{}
}

/**
 * Get auth admin access token JWT middleware instance
 */
func GetAuthAdminAccessTokenJwtMiddleware() *AuthAdminAccessTokenJwtMiddleware {
	return &AuthAdminAccessTokenJwtMiddleware{}
}

/**
 * Get validate middleware instance
 */
func GetValidateMiddleware() *ValidateMiddleware {
	return &ValidateMiddleware{}
}

/**
 * Get auth access token JWT middleware instance
 */
func GetAuthAccessTokenJwtMiddleware() *AuthAccessTokenJwtMiddleware {
	return &AuthAccessTokenJwtMiddleware{}
}
