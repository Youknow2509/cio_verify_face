package errors

/**
 * Define errors return by service
 */
type Error struct {
	ErrorSystem error  `json:"error_system"`
	ErrorClient string `json:"error_client"`
}
