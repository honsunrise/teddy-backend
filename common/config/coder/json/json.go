package json

import (
	"encoding/json"
	"teddy-backend/common/config"
)

type jsonEncoder struct{}

func (j jsonEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonEncoder) Decode(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

func NewCoder() config.Coder {
	return jsonEncoder{}
}
