package main

import (
	"fmt"
	"flag"
	"os"
)

type Command struct {
	options *[]Option
	command string // ie 'ble', 'ping', 'ls', 'camera'
	min_opts int // minimum options to provide
	max_opts int //
}

type Option[T any] struct {
	val_range OverRange
	name string
	usage string
	val T // Subcommand values in single, can support comma seperated field (?)
	// 0 for the following if such doesnt require bounding values
	minv int // Assumption: float was not found to be parsed in originial program (if it was it was parsed as a string, see 'ping.py')
	maxv int // This program will keep [minv, maxv) as integer values
}

// maybe better way ?
type OverRange struct {
	lower int
	upper int
}

func ParseArgs(c Command) {
	set := flag.NewFlagSet(c.command, flag.ExitOnError)

	// check if valid values
	for i, o := range c.options {
		set.Var(o.val, o.name, o.usage)
	}

	set.Parse(os.Args[2:])

}

func CreateOpt[T any](name string, usage string, vals T, minv int, maxv int) *Option {

	opt := &Option {
		name: name,
		usage: usage,
		vals: vals,
		minv: minv,
		maxv: maxv,
	}

	return opt
}

func CreateCommand(name string, min_opts int, max_opts int, opts []Option) *Command {
	command := &Command {
		name: name,
		options: opts,
		min_opts: min_opts,
		max_opts: max_opts,
	}

	return command
}

func main() {
	fmt.Println("Hello, World!")
	opts := make([]*Option, 3)

	opts = append(opts, CreateOpts("scanon", "Turns BLE scan on", false, 0, 0))

	command := CreateCommand("ble", 1, 1, opts)

}
