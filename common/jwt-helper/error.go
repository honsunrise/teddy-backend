package jwt_helper

import "errors"

var (
	// ErrMissingRealm indicates Realm name is required
	ErrMissingRealm = errors.New("realm is missing")

	ErrMissingSigningAlgorithm = errors.New("signing algorithm is missing")

	ErrMissingKeyFunction = errors.New("key function is missing")

	ErrNotSupportSigningAlgorithm = errors.New("this signing algorithm NOT support")

	ErrTokenInvalid = errors.New("token is invalid")

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	// ErrInvalidKey indicates the the given public key is invalid
	ErrInvalidKey = errors.New("key invalid")
)
