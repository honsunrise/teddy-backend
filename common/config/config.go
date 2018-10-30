package config

import (
	"flag"
)

var config Config

var consulAddress string
var consulPrefix string

func init() {
	flag.StringVar(&consulAddress, "consul_addr", "127.0.0.1:8500", "the consul address default 127.0.0.1:8500")
	flag.StringVar(&consulPrefix, "consul_prefix", "teddy/config", "the consul config prefix default /teddy/config")
}

func Init() error {
	return nil
}

func GetConfig() *Config {
	return &config
}

func Watch() <-chan *Config {
	return nil
}
