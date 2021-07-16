package algorithms

import "errors"

var (
	ErrUnexpectedSigningMethod         = errors.New("unexpected signing method")
	ErrUnsupportedAuthenticationMethod = errors.New("unsupported authentication method")
)
