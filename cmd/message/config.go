package main

import "teddy-backend/internal/types"

type Config struct {
	Server    types.Server      `mapstructure:"server"`
	Databases map[string]string `mapstructure:"databases"`
	Mail      types.Mail        `mapstructure:"mail"`
}
