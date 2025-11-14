package constants

// ========================
//
//	For ttl cache
//
// ========================

const (
	// Shift
	TTL_Shift_Cache                 = 60 * 30 // 30 minutes
	TTL_Shift_Employee_Cache        = 60 * 30 // 30 minutes
	TTL_Info_Base_Employee_In_Shift = 60 * 5  // 5 minutes
	TTL_List_Employee_Shift_Cache   = 60 * 10 // 10 minutes
)

const (
	// Save second
	TTL_OTP_REGISTER                = 60 * 8            // 8 minutes
	TTL_LOCAL_OTP_REGISTER          = 3                 // 3 seconds
	TTL_TOKEN_VERIFY_REGISTER       = 60 * 30           // 30 minutes
	TTL_ACCESS_TOKEN                = 60 * 60 * 2       // 2 hour
	TTL_REFRESH_TOKEN               = 60 * 60 * 24 * 7  // 7 days
	TTL_RESET_PASSWORD_TOKEN        = 60 * 15           // 15 minutes
	TTL_FRIEND_REQUEST              = 60 * 5            // 5 minutes
	TTL_USER_INFO_VIEW              = 60 * 10           // 10 minutes
	TTL_LIST_FRIEND_REQUEST_TO_USER = 60 * 3            // 3 minutes
	TTL_LIST_FRIENDS_OF_USER        = 60 * 10           // 10 minutes
	TTL_TOKEN_DEVICE                = 60 * 60 * 24 * 30 // 30 days
	TTL_DEVICE_INFO                 = 60 * 30           // 30 minutes

	// Local cache TTLs (shorter for faster invalidation)
	TTL_LOCAL_USER_INFO_VIEW  = 60 * 2 // 2 minutes
	TTL_LOCAL_ACCESS_TOKEN    = 60 * 5 // 5 minutes
	TTL_LOCAL_USER_PERMISSION = 60 * 1 // 1 minute
	TTL_LOCAL_DEVICE_CHECK    = 60 * 1 // 1 minute

	// Additional cache TTLs
	TTL_USER_PERMISSION = 60 * 10 // 10 minutes
	TTL_DEVICE_CHECK    = 60 * 5  // 5 minutes
)

// For spam and count spam
const (
	// Save second
	TTL_BLOCK_SPAM_REGISTER       = 5 * 60 // 5 minutes
	TTL_LOCAL_BLOCK_SPAM_REGISTER = 3      // 3 seconds

	TTL_COUNT_SPAM_REGISTER       = 60 * 60 // 1 hour
	TTL_LOCAL_COUNT_SPAM_REGISTER = 2       // 2 seconds

	TTL_COUNT_VERIFY_OTP_REGISTER_FAIL       = 60 * 60 // 1 hour
	TTL_LOCAL_COUNT_VERIFY_OTP_REGISTER_FAIL = 2       // 2 seconds

	TTL_BLOCK_SPAM_REFRESH_TOKEN       = 60 * 60 // 1 hour
	TTL_LOCAL_BLOCK_SPAM_REFRESH_TOKEN = 5       // 5 seconds

	TTL_BLOCK_SPAM_RESET_PW       = 60 * 30 // 30 minutes
	TTL_LOCAL_BLOCK_SPAM_RESET_PW = 3       // 3 seconds

	TTL_COUNT_SPAM_RESET_PW       = 60 * 30 // 30 minutes
	TTL_LOCAL_COUNT_SPAM_RESET_PW = 2       // 2 seconds

	TTL_BLOCK_SPAM_VERIFY_RESET_PW_TK       = 60 * 60 * 12 // 12 hours
	TTL_LOCAL_BLOCK_SPAM_VERIFY_RESET_PW_TK = 60 * 10      // 10 minutes

	// v.v
)
