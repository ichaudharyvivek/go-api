package errors

// Ideally this will include the formatted error strings
var (
	RespDBDataInsertFailure = "db data insert failure"
	RespDBDataAccessFailure = "db data access failure"
	RespDBDataUpdateFailure = "db data update failure"
	RespDBDataRemoveFailure = "db data remove failure"

	RespJSONEncodeFailure = "json encode failure"
	RespJSONDecodeFailure = "json decode failure"

	RespInvalidURLParamID = "invalid url param-id"
)

var (
	ResourceNotFound = "Resource '%s' not found"
	UserNotFound     = "user with id: '%s' not found"
)
