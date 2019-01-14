package main

import "teddy-backend/common/types"

type Config struct {
	Server    types.Server      `json:"server"`
	Databases map[string]string `json:"databases"`
}
