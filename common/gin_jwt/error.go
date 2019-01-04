package gin_jwt

import (
	"errors"
	"github.com/zhsyourai/teddy-backend/common/nice_error"
	"net/http"
)

var (
	ErrMissingRealm = nice_error.DefineNiceError(http.StatusInternalServerError, "realm is missing", "please set realm")

	ErrMissingKeyFunction = errors.New("key function is missing")

	ErrMissingSigningAlgorithm = errors.New("signing algorithm is missing")

	ErrContextNotHaveToken = errors.New("context not have token")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrTokenInvalid = errors.New("token is invalid")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	// ErrInvalidKey indicates the the given public key is invalid
	ErrInvalidKey = errors.New("key invalid")
)
