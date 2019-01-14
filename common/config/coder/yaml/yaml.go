package yaml

import (
	"github.com/ghodss/yaml"
	"teddy-backend/common/config"
)

type yamlEncoder struct{}

func (y yamlEncoder) Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y yamlEncoder) Decode(d []byte, v interface{}) error {
	return yaml.Unmarshal(d, v)
}

func NewCoder() config.Coder {
	return yamlEncoder{}
}
