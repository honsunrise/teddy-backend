package config

import "errors"

var WatchNotSupport = errors.New("not support watch")

type WatcherStopper interface {
	Stop() error
}

type WatchResult struct {
	WatcherStopper
	Results chan<- map[string]interface{}
	Errors  chan<- error
}
