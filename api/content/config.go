package main

import "github.com/zhsyourai/teddy-backend/common/types"

type ObjectStore struct {
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type Config struct {
	Server      types.Server            `mapstructure:"server"`
	ObjectStore map[string]*ObjectStore `mapstructure:"object_store"`
}
