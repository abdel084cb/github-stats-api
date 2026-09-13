package apperrors

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUpstreamUnavailable = errors.New("upstream unavailable")
	ErrInvalidUsername     = errors.New("invalid username")
)
