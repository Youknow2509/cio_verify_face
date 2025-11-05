package middleware

/**
 * Get auth admin middleware instance
 */
func GetAuthAdminMiddleware() *AuthAdminAccessTokenJwtMiddleware {
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
