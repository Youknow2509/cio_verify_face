package error

/**
 * Error definitions for application errors
 */
type Error struct {
	ErrorSystem error  `json:"error_system"`
	ErrorClient string `json:"error_client"`
}
