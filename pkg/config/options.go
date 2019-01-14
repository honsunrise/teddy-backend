package config

type Options map[interface{}]interface{}

type Option func(o Options)

func BuildOptions(opts ...Option) Options {
	options := make(Options)
	for _, o := range opts {
		o(options)
	}

	return options
}
