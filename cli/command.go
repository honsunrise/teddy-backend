package commands

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"teddy-backend/cli/version"
)

var app = kingpin.New("teddy-cli", "Teddy backend command line tools")

var registry = make(map[string]func() error)

func register(processors map[string]func() error) {
	for command, processor := range processors {
		if registry[command] != nil {
			panic(fmt.Errorf("command %q is already registered", command))
		}
		registry[command] = processor
	}
}

func init() {
	register(version.Prepare(app))
}

func Run() int {
	selected := kingpin.MustParse(app.Parse(os.Args[1:]))
	if selected == "" {
		app.Usage(os.Args[1:])
		return -1
	}
	processor := registry[selected]
	if processor == nil {
		panic(fmt.Errorf("command %q not found", selected))
	}
	if processor() != nil {
		return -1
	}
	return 0
}
