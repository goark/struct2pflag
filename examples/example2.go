//go:build run

package main

import (
	"fmt"

	"github.com/goark/struct2pflag"
	"github.com/spf13/pflag"
)

type SubConfig1 struct {
	Debug bool `pflag:"debug,d,Enable debug mode"`
}

type SubConfig2 struct {
	Verbose bool `pflag:"verbose,v,Enable verbose mode"`
}

type Config struct {
	Name string `pflag:"name,Set your name"`
	Sub1 SubConfig1
	Sub2 SubConfig2 `pflag:""`
}

func (c *Config) Run() {
	fmt.Printf("Name=%#v\n", c.Name)
	fmt.Printf("Sub1.Debug=%#v\n", c.Sub1.Debug)
	fmt.Printf("Sub2.Verbose=%#v\n", c.Sub2.Verbose)
}

func main() {
	var cfg Config
	struct2pflag.BindDefault(&cfg)
	pflag.Parse()

	cfg.Run()
}
