package errors

// Static error messages
var (
	// General error(s)
	SomethingWentWrong    = "something went wrong"
	InternalServerError   = "internal server error"
	PasswordHashingFailed = "password hashing failed"

	// DB errors
	DBDataInsertFailure = "db data insert failure"
	DBDataAccessFailure = "db data access failure"
	DBDataUpdateFailure = "db data update failure"
	DBDataRemoveFailure = "db data remove failure"

	// Json errors
	JSONEncodeFailure = "json encode failure"
	JSONDecodeFailure = "json decode failure"

	// Request errors
	InvalidURLParamID = "invalid url param- id"
)

// Formatted error messages
var (
	// Resource erros
	ResourceNotFound = "resource '%s' not found"
	UserNotFound     = "user with id: '%s' not found"
)
