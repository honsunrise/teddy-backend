package config

import (
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"sync"
	"time"
)

type Config struct {
	exit       chan bool
	sources    []Source
	resultMap  map[string]interface{}
	modTimeMap map[interface{}]time.Time
	rwmutex    sync.RWMutex
}

func NewConfig(sources ...Source) (*Config, error) {
	config := &Config{
		exit:       make(chan bool),
		resultMap:  make(map[string]interface{}),
		modTimeMap: make(map[interface{}]time.Time),
		sources:    sources,
	}

	if err := config.Fetch(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Fetch() error {
	c.rwmutex.Lock()
	defer c.rwmutex.Unlock()

	for _, s := range c.sources {
		if modTime, ok := c.modTimeMap[s]; ok {
			tmpModTime, err := s.LastModifyTime()
			if err != nil {
				return err
			}
			tmpMap, err := s.Read()
			if err != nil {
				return err
			}
			if modTime.Before(tmpModTime) {
				if err := mergo.Map(&c.resultMap, tmpMap, mergo.WithOverride); err != nil {
					return err
				}
			}
		} else {
			tmpModTime, err := s.LastModifyTime()
			if err != nil {
				return err
			}
			tmpMap, err := s.Read()
			if err != nil {
				return err
			}
			c.modTimeMap[s] = tmpModTime
			if err := mergo.Map(&c.resultMap, tmpMap, mergo.WithOverride); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Config) Scan(v interface{}) error {
	c.rwmutex.RLock()
	defer c.rwmutex.RUnlock()

	return mapstructure.Decode(c.resultMap, v)
}

func (c *Config) Watch() (*WatchResult, error) {
	return nil, nil
}

func (c *Config) Close() error {
	c.rwmutex.Lock()
	defer c.rwmutex.Unlock()

	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
	}
	return nil
}
