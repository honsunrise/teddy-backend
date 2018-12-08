package main

import "github.com/zhsyourai/teddy-backend/common/types"

type Config struct {
	Server    types.Server              `json:"server"`
	Databases map[string]types.Database `json:"databases"`
}
