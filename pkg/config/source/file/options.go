package file

import (
	"teddy-backend/pkg/config"
)

type filePathKey struct{}
type fileFormatKey struct{}

func WithPath(p string) config.Option {
	return func(o config.Options) {
		o[filePathKey{}] = p
	}
}

func WithFormat(f config.SourceFormat) config.Option {
	return func(o config.Options) {
		o[fileFormatKey{}] = f
	}
}
