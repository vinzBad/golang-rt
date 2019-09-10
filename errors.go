package rt

import "errors"

var (
	// ErrInvalidAPIURL is returned by the New method of rt
	ErrInvalidAPIURL = errors.New("invalid api url")
	// ErrParseRTMessageError is returned when rt is unable to parse a message by RequestTracker
	ErrParseRTMessageError = errors.New("failed to parse RequestTracker message")
	// ErrCredentialsNeeded is returned when rt receives a 401 from RequestTracker
	ErrCredentialsNeeded = errors.New("credentials needed")
)
