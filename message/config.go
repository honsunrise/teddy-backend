package main

import "github.com/zhsyourai/teddy-backend/common/types"

type Config struct {
	Server    types.Server      `mapstructure:"server"`
	Databases map[string]string `mapstructure:"databases"`
	Mail      types.Mail        `mapstructure:"mail"`
}
