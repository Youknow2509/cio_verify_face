package errors

// =======================================================
//
//	Define device errors returned by the service
//
// =======================================================
const (
	DeviceNotFoundErrorCode       = 30001
)

var mapDeviceErrors = map[int]string{
	DeviceNotFoundErrorCode: "device not found",
}
