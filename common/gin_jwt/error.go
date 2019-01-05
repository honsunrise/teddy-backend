package gin_jwt

import (
	"github.com/zhsyourai/teddy-backend/common/nice_error"
	"net/http"
)

var (
	ErrMissingRealm = nice_error.DefineNiceError(http.StatusInternalServerError, "realm is missing")

	ErrMissingKeyFunction = nice_error.DefineNiceError(http.StatusInternalServerError, "key function is missing")

	ErrMissingSigningAlgorithm = nice_error.DefineNiceError(http.StatusInternalServerError, "signing algorithm is missing")

	ErrInvalidKey = nice_error.DefineNiceError(http.StatusInternalServerError, "key invalid")

	ErrContextNotHaveToken = nice_error.DefineNiceError(http.StatusInternalServerError, "context not have token")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = nice_error.DefineNiceError(http.StatusForbidden, "you don't have permission to access this resource")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrTokenInvalid = nice_error.DefineNiceError(http.StatusUnauthorized, "token is invalid")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = nice_error.DefineNiceError(http.StatusUnauthorized, "auth header is invalid")
)
