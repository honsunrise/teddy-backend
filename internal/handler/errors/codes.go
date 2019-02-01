package errors

const (
	ErrCodeUnknown = 10000 + iota
	ErrCodeInternal
	ErrCodeUnauthorized
	ErrCodeForbidden
	ErrCodeBadRequest
	ErrCodeUsernameOrPasswordNotCorrect
	ErrCodeCaptchaIDNotFound
	ErrCodeCaptchaExtNotSupport
	ErrCodeCaptchaNotCorrect
	ErrCodeRegisterTypeNotSupport
	ErrCodeAccountExists
)
