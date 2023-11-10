package utils

const (
	REQ_WRONG_BODY_FORMAT       = 1
	REQ_FIELD_ERROR             = 2
	DB_QUERY_ERROR              = 3
	AUTH_MISSING_JWT            = 4
	AUTH_WRONG_JWT              = 5
	AUTH_TOKEN_EXPIRED          = 6
	AUTH_EMAIL_NOT_FOUND        = 7
	AUTH_WRONG_PASSWORD         = 8
	AUTH_USER_NOT_LOGGEDIN      = 9
	AUTH_DIFFERENT_IP_OR_DEVICE = 10
	RESETPW_INVALID_TOKEN       = 11
	RESETPW_EXPIRED_TOKEN       = 12
	RESETPW_MAX_LIMIT_REQUEST   = 13
	FAILED_GENERATE_JWT         = 14
	FAILED_ENCRYPT              = 15
	VERIFY_INVALID_TOKEN        = 16
	VERIFY_MAX_LIMIT_REQUEST    = 17
	VERIFY_ALREADY_VERIFIED     = 18
)
