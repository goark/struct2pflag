//go:build run

package main

import (
	"fmt"

	"github.com/goark/struct2pflag"
	"github.com/spf13/pflag"
)

type Env struct {
	B bool   `pflag:"boolean,b,This is a boolean flag"`
	N int    `pflag:"integer,n,This is an integer flag"`
	S string `pflag:"string,s,this is a string flag"`
}

func (e Env) Run() {
	fmt.Printf("B=%#v\n", e.B)
	fmt.Printf("N=%#v\n", e.N)
	fmt.Printf("S=%#v\n", e.S)
}

func main() {
	var env Env
	struct2pflag.BindDefault(&env)
	pflag.Parse()
	env.Run()
}
