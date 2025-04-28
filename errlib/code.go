package errlib

// Code is an application specific error code that can be mapped to ports layers status codes (HTTP, gRPC, etc)
type Code string

const (
	InternalCode      Code = "INTERNAL"
	NotFoundCode      Code = "NOT_FOUND"
	InvalidInputCode  Code = "INVALID_INPUT"
	ConflictCode      Code = "CONFLICT"
	UnprocessableCode Code = "UNPROCESSABLE"
	UnauthorizedCode  Code = "UNAUTHORIZED"
)
