package response

const (
	//
	ErrorCodeBindRequest     = 1001
	ErrorCodeValidateRequest = 1002
	//
	ErrorCodeSystemTemporary = 2001
	//
	ErrorCodeDeleteDevice = 3001
)

// message
var msg = map[int]string{
	ErrorCodeDeleteDevice:    "Failed to delete device",
	ErrorCodeValidateRequest: "Validation failed",
	ErrorCodeBindRequest:     "Failed to bind request",
	ErrorCodeSystemTemporary: "System is temporary unavailable",
}
