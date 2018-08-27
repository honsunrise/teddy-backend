package config

import (
	"flag"
	mconfig "github.com/micro/go-config"
	"github.com/micro/go-config/encoder/yaml"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/consul"
	"github.com/micro/go-config/source/env"
	mflag "github.com/micro/go-config/source/flag"
)

var config Config

var consulAddress string
var consulPrefix string

func init() {
	flag.StringVar(&consulAddress, "consul_addr", "127.0.0.1:8500", "the consul address default 127.0.0.1:8500")
	flag.StringVar(&consulPrefix, "consul_prefix", "teddy/config", "the consul config prefix default /teddy/config")
}

func Init() error {
	err := mconfig.Load(
		env.NewSource(
			env.WithStrippedPrefix("TEDDY"),
		),
		mflag.NewSource(),
		consul.NewSource(
			consul.WithAddress(consulAddress),
			consul.WithPrefix(consulPrefix),
			consul.StripPrefix(true),
			source.WithEncoder(yaml.NewEncoder()),
		),
	)
	if err != nil {
		return err
	}

	err = mconfig.Scan(&config)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	return &config
}

func Watch() <-chan *Config {
	return nil
}
