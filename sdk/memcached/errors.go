package memcached

import "errors"

var (
	ErrNotStored          = errors.New("not stored")
	ErrInvalidValueHeader = errors.New("invalid value header")
	ErrNotFound           = notFoundError{errors.New("not found")}
	ErrUnknownResponse    = unknownError{errors.New("unknown response")}
)

type notFoundError struct{ error }
type unknownError struct{ error }

func (notFoundError) NotFoundErrorMarker() {}
func (unknownError) UnknownErrorMarker()   {}
