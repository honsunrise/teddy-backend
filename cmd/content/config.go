package main

import "teddy-backend/internal/types"

type ObjectStore struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type Config struct {
	Server      types.Server            `mapstructure:"server"`
	Databases   map[string]string       `mapstructure:"databases"`
	ObjectStore map[string]*ObjectStore `mapstructure:"object_store"`
}
