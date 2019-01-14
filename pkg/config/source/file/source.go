// Package file is a file source. Expected format is json
package file

import (
	"io/ioutil"
	"os"
	"path"
	"teddy-backend/pkg/config"
	"teddy-backend/pkg/config/coder/json"
	"teddy-backend/pkg/config/coder/yaml"
	"time"
)

type file struct {
	path  string
	coder config.Coder
}

var (
	DefaultPath = "config.json"
)

func (l *file) Read() (map[string]interface{}, error) {
	f, err := os.Open(l.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	err = l.coder.Decode(b, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (l *file) Watch(path ...string) (*config.WatchResult, error) {
	if _, err := os.Stat(l.path); err != nil {
		return nil, err
	}
	return newWatcher(l, path)
}

func (l *file) LastModifyTime() (time.Time, error) {
	stat, err := os.Stat(l.path)
	if err != nil {
		return time.Time{}, err
	}
	return stat.ModTime(), nil
}

func NewSource(opts ...config.Option) config.Source {
	options := config.BuildOptions(opts...)
	tmpPath := DefaultPath
	var tmpFormat config.SourceFormat
	p, ok := options[filePathKey{}].(string)
	if ok {
		tmpPath = p
	}

	var coder config.Coder = nil
	f, ok := options[fileFormatKey{}].(config.SourceFormat)
	if ok {
		tmpFormat = f
	} else {
		ext := path.Ext(tmpPath)
		switch ext {
		case ".json":
			tmpFormat = config.Json
		case ".yaml":
			tmpFormat = config.Yaml
		}
	}
	switch tmpFormat {
	case config.Json:
		coder = json.NewCoder()
	case config.Yaml:
		coder = yaml.NewCoder()
	}

	return &file{path: tmpPath, coder: coder}
}
