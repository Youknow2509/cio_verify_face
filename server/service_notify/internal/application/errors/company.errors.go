package errors

// =======================================================
//
//	Define company errors returned by the service
//
// =======================================================
const (
	DeviceNotFoundErrorCode = 30001
)

var mapCompanyErrors = map[int]string{
	DeviceNotFoundErrorCode: "Device not found",
}
