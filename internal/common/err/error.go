package err

type Error struct {
	Error string `json:"error"`
}

type Errors struct {
	Errors []string `json:"errors"`
}

var (
	RespDBDataInsertFailure = "db data insert failure"
	RespDBDataAccessFailure = "db data access failure"
	RespDBDataUpdateFailure = "db data update failure"
	RespDBDataRemoveFailure = "db data remove failure"

	RespJSONEncodeFailure = "json encode failure"
	RespJSONDecodeFailure = "json decode failure"

	RespInvalidURLParamID = "invalid url param-id"
)
