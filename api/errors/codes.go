package errors

const (
	ErrCodeUnknown = 10000 + iota
	ErrCodeInternal
	ErrCodeUnauthorized
	ErrCodeForbidden
	ErrCodeCaptchaIDNotFound
	ErrCodeCaptchaExtNotSupport
)
