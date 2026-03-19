package apperr

type Code string

const (
	Unknown           Code = "UNKNOWN"
	Unexpected        Code = "UNEXPECTED"
	Timeout           Code = "TIMEOUT"
	NetworkConnection Code = "NETWORK_CONNECTION"
	BadRequest        Code = "BAD_REQUEST"
	InvalidData       Code = "INVALID_DATA"
	Inconsistency     Code = "INCONSISTENCY"
	Unauthenticated   Code = "UNAUTHENTICATED"
	NotAllowed        Code = "NOT_ALLOWED"
)

type Kind string

const (
	Validation        Kind = "Validation"
	Request           Kind = "Request"
	Unauthorized      Kind = "Unauthorized"
	Forbidden         Kind = "Forbidden"
	Conflict          Kind = "Conflict"
	Internal          Kind = "Internal"
	External          Kind = "External"
)
