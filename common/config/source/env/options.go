package env

import (
	"github.com/zhsyourai/teddy-backend/common/config"
	"strings"
)

type strippedPrefixKey struct{}
type prefixKey struct{}

func WithStrippedPrefix(p ...string) config.Option {
	return func(o config.Options) {
		o[strippedPrefixKey{}] = normal(p)
	}
}

func WithPrefix(p ...string) config.Option {
	return func(o config.Options) {
		o[prefixKey{}] = normal(p)
	}
}

func normal(prefixes []string) []string {
	var result []string
	for _, p := range prefixes {
		if !strings.HasSuffix(p, "_") {
			result = append(result, p+"_")
		} else {
			result = append(result, p)
		}
	}

	return result
}
