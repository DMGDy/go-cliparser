package main

import (
	"fmt"
	"flag"
	"os"
)

// treat as C union (change later)
type Value struct {
	val_type string
	val interface{}
}

// maybe better way ?
type OverRange struct {
	lower int
	upper int
}

type Command struct {
	options []Option
	name string // ie 'ble', 'ping', 'ls', 'camera'
	min_opts int // minimum options to provide
	max_opts int //
}

type Option struct {
	val_range OverRange
	name string
	usage string
	val Value
	// 0 for the following if such doesnt require bounding values
	minv int // Assumption: float was not found to be parsed in originial program (if it was it was parsed as a string, see 'ping.py')
	maxv int // This program will keep [minv, maxv) as integer values
}

func BoolVal(b bool) Value {
	return Value {
		val_type: "bool",
		val: b,
	}
}

func IntVal(i int) Value {
	return Value {
		val_type: "int",
		val: i,
	}
}

func StringVal(s string) Value {
	return Value {
		val_type: "string",
		val: s,
	}
}

func Float64Val(f float64) Value {
	return Value {
		val_type: "float64",
		val: f,
	}
}

func ParseArgs(c Command) *flag.FlagSet{
	set := flag.NewFlagSet(c.name, flag.ExitOnError)


	// check if valid values
	fmt.Println(len(c.options))
	for _, o := range c.options {
		fmt.Println(o.val.val_type)
		switch o.val.val_type {
		case "int":
			set.Int(o.name, o.val.val.(int), o.usage)
		case "string":
			set.String(o.name, o.val.val.(string), o.usage)
		case "float64":
			set.Float64(o.name, o.val.val.(float64), o.usage)
		case "bool":
			set.Bool(o.name, o.val.val.(bool), o.usage)
		default:
			fmt.Fprintf(os.Stderr, "Error parsing arguments!\n")
			os.Exit(1)
		}
	}
	set.Parse(os.Args[2:])

	return set
}

func CreateOpts(name string, usage string, val Value, minv int, maxv int) Option {

	return Option {
		name: name,
		usage: usage,
		val: val,
		minv: minv,
		maxv: maxv,
	}
}

func CreateCommand(name string, min_opts int, max_opts int, opts []Option) Command {

	return Command {
		name: name,
		options: opts,
		min_opts: min_opts,
		max_opts: max_opts,
	}
}

func main() {
	fmt.Println("Hello, World!")
	var opts []Option

	bvalue := BoolVal(false)
	opts = append(opts, CreateOpts("scanon", "Turns BLE scan on.", bvalue, 0, 0))

	bvalue = BoolVal(true)
	opts = append(opts, CreateOpts("scanoff", "Turns BLE scan on.", bvalue, 0, 0))

	ivalue := IntVal(1000800)
	opts = append(opts, CreateOpts("uuid", "This indicates which BLE process to send this message to.", ivalue, 0, 99999999))

	command := CreateCommand("ble", 1, 1, opts)

	flag_set := ParseArgs(command)
	scanon := flag_set.Lookup("scanon")
	
	if scanon == nil {
		fmt.Fprintf(os.Stderr,"Error getting flag value\n")
		os.Exit(1)
	}

	fmt.Println(scanon.Value)

	uuid := flag_set.Lookup("uuid")
	
	if scanon == nil {
		fmt.Fprintf(os.Stderr,"Error getting flag value\n")
		os.Exit(1)
	}

	fmt.Println(uuid.Value)
}
