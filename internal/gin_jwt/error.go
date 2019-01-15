package gin_jwt

import (
	"errors"
)

var (
	ErrMissingRealm = errors.New("realm is missing")

	ErrMissingKeyFunction = errors.New("key function is missing")

	ErrMissingSigningAlgorithm = errors.New("signing algorithm is missing")

	ErrInvalidKey = errors.New("key invalid")

	ErrContextNotHaveToken = errors.New("context not have token")

	ErrForbidden = errors.New("you don't have permission to access this resource")

	ErrTokenInvalid = errors.New("token is invalid")

	ErrInvalidAuthHeader = errors.New("auth header is invalid")
)
