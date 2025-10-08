package errors

/**
 * Define errors return by service
 */
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


/**
 * Create a new error
 */
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

/**
 * Get error with code
 */
func GetError(code int) *Error {
	// merge map data
	data := map[int]string{}
	for k, v := range mapUserErrors {
		data[k] = v
	}
	for k, v := range mapSystemErrors {
		data[k] = v
	}
	for k, v := range mapAuthErrors {
		data[k] = v
	}
	for k, v := range mapCompanyErrors {
		data[k] = v
	}
	// v.v
	if msg, ok := data[code]; ok {
		return NewError(code, msg)
	}
	return nil
}

