package middleware

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