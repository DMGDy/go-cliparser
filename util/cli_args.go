package util

// usage format to add
// ./cli [COMMAND] --[ARGUMENTS] | --[ARGUMENTS] [VAL] | -[ARGUMENTS] [VAL] | --[ARGUMENTS]=[VAL] ( --verbose | -verbose )
import (
	"flag"
	"os"
	"fmt"
	"strconv"
	"strings"
)

// treat as C union (change later)
type Value struct {
	val_type string
	// empty interface for accepting any time
	val interface{}
}

// maybe better way ?
type Range struct {
	lower int
	upper int
}

type Command struct {
	subcommands []Subcommand
	Name string // ie 'ble', 'ping', 'ls', 'camera'
	// rename to args
	min_opts int // minimum options to provide
	max_opts int //
}

// rename
type Subcommand struct {
	// if wanted to do something in a range of 
	val_range Range
	Name string
	usage string
	defval Value
	// This program will keep [minv, maxv) as integer values wrapped in Range type
	minmaxv Range
}

// Actual CLI arguments are given by Flag
type ArgFlag struct {
	FlagSet *flag.FlagSet
	command *Command
}

// to be used and exported outside this package
// string being the command (ie. 'ble', 'camera', 'ls')
var ArgMap = make(map[string]*ArgFlag)

var SubCmdVal = make(map[string]string)

// Satisfy flag.Value interface
func (r *Range) Set(s string) error {
	val := RangeVal(s).val.(Range)
	r.lower = val.lower
	r.upper = val.upper
	// no error returned, program will exit with error
	return nil
}

func (r *Range) String() string {
	return fmt.Sprintf("%d-%d", r.lower, r.upper)
}

func intInRange(n int, r Range) bool {
	in_range := false

	if r.upper == r.lower {
		in_range = true
	}

	if n >= r.lower && n < r.upper {
		in_range = true
	}

	return in_range
}

// validate the passed values for the command
// return map of subcommands to values (if any)
func ValidateValues(af *ArgFlag) {
	// read and assign the flags with values provided or assign default values if available
	fmt.Println(os.Args)
	af.FlagSet.Parse(os.Args[2:])
	for _, subcommand := range af.command.subcommands{
		name := subcommand.name
		flag := af.FlagSet.Lookup(name)
		// if 'nil' subcommand is invalid
		//	- if time permits, can implement levenshtein distance to see closest subcommand
		//	- same check can be used on main Command
		if flag == nil {
			fmt.Errorf("Unrecognized subcommand: %s\n", name)
			os.Exit(1)
		}

		val := flag.Value.String()

		val_type := subcommand.defval.val_type
		

		// need to cast as returned value is just string
		switch val_type {
		case "int":
			int_val, err := strconv.Atoi(val)
			if err != nil {
				fmt.Errorf("Could not parse as an integer: %s\n", val)
			}
			// add to map if it is in range
			if intInRange(int_val , subcommand.val_range) {
				// can provide just string, since it will become string in sending the command
				SubCmdVal[subcommand.name] = val
			}

		case "string":
				SubCmdVal[subcommand.name] = val
		case "bool":
				SubCmdVal[subcommand.name] = val
		case "float64":
			// implement range value testing
			SubCmdVal[subcommand.name] = val
		// TODO: implement Range type parsing
		case "range":
			// string -> Range
			split := strings.Split(val, "-")
			if len(split) < 2 {
				// should never get here
				fmt.Errorf("Range in incorrect form: %s\n", val)
				os.Exit(1)
			}
		default:
			fmt.Errorf("Unrecognized type: %s\n", val_type)
			os.Exit(1)
		}
		
	}
}

// On command implementer to use these properly otherwise parser will fail
func RangeVal(s string) Value {
	split := strings.Split(s, "-")

	if len(split) < 2 {
		fmt.Errorf("Not in form of a propper Range: %s\n", s)
		fmt.Errorf("Hint: Range should be `1-10`, `5-100`, '4-5`\n")
		os.Exit(1)
	}

	lower,err := strconv.Atoi(split[0])
	upper, err:= strconv.Atoi(split[1])

	if err != nil {
		fmt.Errorf("Not in form of a propper Range: %s\n", s)
		fmt.Errorf("Hint: Range should be `1-10`, `5-100`, '4-5`\n")
		os.Exit(1)
	}

	return Value {
		val_type: "range",
		val: Range {
			lower: lower,
			upper: upper,
		},
	}
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

func ParseCommand(c Command) {
	set := flag.NewFlagSet(c.name, flag.ExitOnError)

	for _, o := range c.subcommands{
		switch o.defval.val_type {
		case "int":
			n := o.defval.val.(int)
			set.Int(o.name, n , o.usage)
		case "string":
			s := o.defval.val.(string)
			set.String(o.name, s, o.usage)
		case "float64":
			f := o.defval.val.(float64)
			set.Float64(o.name, f, o.usage)
		case "bool":
			b := o.defval.val.(bool)
			set.Bool(o.name, b, o.usage)
		case "range":
			r := o.defval.val.(Range)
			set.Var(&r, o.name, o.usage)
		default:
			os.Exit(1)
		}
	}
	// parsing now for testing, real usage parse ONLY ONCE ALL commands are processed

	ArgMap[c.name] = &ArgFlag {
		FlagSet: set,
		command: &c,
	}
}

// add parameters for ranges
func CreateSubCmd(name string, usage string, val Value, minv int, maxv int, low_range int, up_range int) Subcommand {
	return Subcommand {
		val_range: Range {
			lower: 0,
			upper: 0,
		},
		name: name,
		usage: usage,
		defval: val,
		minmaxv: Range {
			lower: minv,
			upper: maxv,
		},
	}
}

func CreateCommand(name string, min_opts int, max_opts int, subcommands []Subcommand) Command {

	return Command {
		name: name,
		min_opts: min_opts,
		max_opts: max_opts,
		subcommands: subcommands,
	}
}
