package gin_jwt

import "errors"

var (
	ErrMissingRealm = errors.New("realm is missing")

	ErrMissingKeyFunction = errors.New("key function is missing")

	ErrMissingSigningAlgorithm = errors.New("signing algorithm is missing")

	ErrContextNotHaveToken = errors.New("context not have token")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrTokenInvalid = errors.New("token is invalid")

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cokie is empty
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrInvalidKey indicates the the given public key is invalid
	ErrInvalidKey = errors.New("key invalid")
)
