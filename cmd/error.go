package cmd

type ErrorKind string

const (
	// ErrInvalidHashLen indicates that the input hash to sign or verify is not
	// the required length = 32-bytes.
	ErrInvalidHashLen = ErrorKind("ErrInvalidHashLen")

	// ErrRelayConnection indicates that an attempt to connect the relay was failed.
	ErrRelayConnection = ErrorKind("ErrRelayConnection")
)

// Error satisfies the error interface and prints human-readable errors.
func (e ErrorKind) Error() string {
	return string(e)
}

// Error identifies an error related to a nostr-cli cmd. It has full
// support for errors.Is and errors.As, so the caller can ascertain the
// specific reason for the error by checking the underlying error.
type Error struct {
	Err         error
	Description string
}

// Error satisfies the error interface and prints human-readable errors.
func (e Error) Error() string {
	return e.Description
}

// Unwrap returns the underlying wrapped error.
func (e Error) Unwrap() error {
	return e.Err
}

// cmdError creates an Error given a set of arguments.
func cmdError(kind ErrorKind, desc string) Error {
	return Error{Err: kind, Description: desc}
}
