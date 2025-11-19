package constants

// Cache TTLs in seconds
const (
	CacheTTLShortSeconds = 60   // 1 minute
	CacheTTLMidSeconds   = 300  // 5 minutes
	CacheTTLLongSeconds  = 1800 // 30 minutes
)

// Cache key prefixes/templates
const (
	// Daily summaries by date (companyID, date YYYY-MM-DD)
	CacheKeyDailyByDate = "analytics:daily:%s:%s"
	// Daily summaries by month (companyID, month YYYY-MM)
	CacheKeyDailyByMonth = "analytics:month:%s:%s"
	// Total employees per company
	CacheKeyTotalEmployees = "analytics:total_employees:%s"
	// Export report cache key (companyID, startDate, endDate, format)
	CacheKeyExportReport = "analytics:export:%s:%s:%s:%s"
)
