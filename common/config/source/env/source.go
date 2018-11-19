package env

import (
	"github.com/imdario/mergo"
	"github.com/zhsyourai/teddy-backend/common/config"
	"os"
	"strings"
	"time"
)

type env struct {
	prefixes         []string
	strippedPrefixes []string
	readTime         time.Time
}

func (e *env) Read() (map[string]interface{}, error) {
	var result map[string]interface{}

	for _, env := range os.Environ() {

		if len(e.prefixes) > 0 || len(e.strippedPrefixes) > 0 {
			notFound := true

			if _, ok := matchPrefix(e.prefixes, env); ok {
				notFound = false
			}

			if match, ok := matchPrefix(e.strippedPrefixes, env); ok {
				env = strings.TrimPrefix(env, match)
				notFound = false
			}

			if notFound {
				continue
			}
		}

		pair := strings.SplitN(env, "=", 2)
		value := pair[1]
		keys := strings.Split(strings.ToLower(pair[0]), "_")
		reverse(keys)

		tmp := make(map[string]interface{})
		for i, k := range keys {
			if i == 0 {
				tmp[k] = value
				continue
			}

			tmp = map[string]interface{}{k: tmp}
		}

		if err := mergo.Map(&result, tmp); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func matchPrefix(pre []string, s string) (string, bool) {
	for _, p := range pre {
		if strings.HasPrefix(s, p) {
			return p, true
		}
	}

	return "", false
}

func reverse(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}

func (e *env) Watch(path ...string) (*config.WatchResult, error) {
	return nil, config.WatchNotSupport
}

func (e *env) LastModifyTime() (time.Time, error) {
	return e.readTime, nil
}

func NewSource(opts ...config.Option) config.Source {
	options := config.BuildOptions(opts...)

	var spkey []string
	var pkey []string
	if p, ok := options[strippedPrefixKey{}].([]string); ok {
		spkey = p
	}

	if p, ok := options[prefixKey{}].([]string); ok {
		pkey = p
	}

	return &env{prefixes: pkey, strippedPrefixes: spkey, readTime: time.Now()}
}
