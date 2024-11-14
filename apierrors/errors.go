package apierrors

import (
	"errors"
)

// A list of error messages
var (
	ErrInternalServer         = errors.New("internal server error")
	ErrInvalidOffsetParameter = errors.New("invalid offset query parameter")
	ErrInvalidLimitParameter  = errors.New("invalid limit query parameter")
	ErrLimitOverMax           = errors.New("limit query parameter is larger than the maximum allowed")
)
