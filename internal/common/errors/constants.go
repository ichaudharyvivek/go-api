package errors

// All possible errors that can happen in an application.
// Most of them are not used in this API but it is good for reference.
//
// NOTE: These are error codes which is used in 'ApiError' struct
const (
	// Database errors
	ErrDBInsertFailure     = "DB_INSERT_FAILURE"
	ErrDBUpdateFailure     = "DB_UPDATE_FAILURE"
	ErrDBDeleteFailure     = "DB_DELETE_FAILURE"
	ErrDBQueryFailure      = "DB_QUERY_FAILURE"
	ErrDBTxBeginFailure    = "DB_TX_BEGIN_FAILURE"
	ErrDBTxCommitFailure   = "DB_TX_COMMIT_FAILURE"
	ErrDBTxRollbackFailure = "DB_TX_ROLLBACK_FAILURE"
	ErrDBNoRows            = "DB_NO_ROWS"
	ErrDBDuplicateEntry    = "DB_DUPLICATE_ENTRY"

	// Validation errors
	ErrValidationFailed     = "VALIDATION_FAILED"
	ErrInvalidRequestBody   = "INVALID_REQUEST_BODY"
	ErrInvalidURLParam      = "INVALID_URL_PARAM"
	ErrMissingRequiredField = "MISSING_REQUIRED_FIELD"
	ErrInvalidEnumValue     = "INVALID_ENUM_VALUE"

	// Authentication & Authorization
	ErrAuthTokenMissing  = "AUTH_TOKEN_MISSING"
	ErrAuthTokenInvalid  = "AUTH_TOKEN_INVALID"
	ErrAuthTokenExpired  = "AUTH_TOKEN_EXPIRED"
	ErrUserNotAuthorized = "USER_NOT_AUTHORIZED"

	// Business logic / domain errors
	ErrUserNotFound     = "USER_NOT_FOUND"
	ErrEmailAlreadyUsed = "EMAIL_ALREADY_USED"
	ErrUsernameTaken    = "USERNAME_TAKEN"
	ErrPasswordMismatch = "PASSWORD_MISMATCH"
	ErrUserBlocked      = "USER_BLOCKED"

	// External service failures
	ErrExternalAPIFailure     = "EXTERNAL_API_FAILURE"
	ErrPaymentGatewayFailure  = "PAYMENT_GATEWAY_FAILURE"
	ErrThirdPartyTimeout      = "THIRD_PARTY_TIMEOUT"
	ErrThirdPartyBadResponse  = "THIRD_PARTY_BAD_RESPONSE"
	ErrFileStorageUnavailable = "FILE_STORAGE_UNAVAILABLE"

	// System / infra
	ErrConfigLoadFailure    = "CONFIG_LOAD_FAILURE"
	ErrEnvMissing           = "ENV_MISSING"
	ErrInternalServer       = "INTERNAL_SERVER_ERROR"
	ErrJSONEncodeFailure    = "JSON_ENCODE_FAILURE"
	ErrJSONDecodeFailure    = "JSON_DECODE_FAILURE"
	ErrFileReadFailure      = "FILE_READ_FAILURE"
	ErrFileWriteFailure     = "FILE_WRITE_FAILURE"
	ErrUnexpectedNilPointer = "UNEXPECTED_NIL_POINTER"

	// Rate limiting / throttling
	ErrRateLimitExceeded = "RATE_LIMIT_EXCEEDED"

	// Feature flags / rollout
	ErrFeatureDisabled = "FEATURE_DISABLED"
)
