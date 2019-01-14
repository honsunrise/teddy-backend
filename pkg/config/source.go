package config

import (
	"errors"
	"time"
)

type SourceFormat uint32

const (
	Json SourceFormat = 1 << iota
	Yaml
)

var FormatNotSupport = errors.New("not support watch")

type Source interface {
	LastModifyTime() (time.Time, error)
	Read() (map[string]interface{}, error)
	Watch(path ...string) (*WatchResult, error)
}
