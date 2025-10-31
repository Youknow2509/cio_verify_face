package errors

// =======================================================
//
//	Define system errors returned by the service
//
// =======================================================
const (
	SystemTemporaryUnavailableErrorCode = 1001
)

var mapSystemErrors = map[int]string{
	SystemTemporaryUnavailableErrorCode: "System temporary unavailable",
}
