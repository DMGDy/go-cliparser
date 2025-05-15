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
	ValType string
	// empty interface for accepting any time
	Val interface{}
}

// maybe better way ?
type Range struct {
	Lower int
	Upper int
}

type Command struct {
	Subcommands []Subcommand
	Name string // ie 'ble', 'ping', 'ls', 'camera'
	// rename to args
	MinSubCmds int // minimum options to provide
	MaxSubCmds int //
}

// rename
type Subcommand struct {
	// if wanted to do something in a range of 
	ValRange Range
	Name string
	Usage string
	DefVal Value
	// This program will keep [minv, maxv) as integer values wrapped in Range type
	MinMaxv Range
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
	val := RangeVal(s).Val.(Range)
	r.Lower = val.Lower
	r.Upper = val.Upper
	// no error returned, program will exit with error
	return nil
}

func rangeInRange(r1 Range, r2 Range) bool {
	in_range := false
	if r1.Lower >= r2.Lower && r1.Upper < r2.Upper {
		in_range := true
	}
	return in_range
}

func (r *Range) String() string {
	return fmt.Sprintf("%d-%d", r.Lower, r.Upper)
}

func intInRange(n int, r Range) bool {
	in_range := false

	if r.Upper == r.Lower {
		in_range = true
	}

	if n >= r.Lower && n < r.Upper {
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
	for _, subcommand := range af.command.Subcommands{
		name := subcommand.Name
		flag := af.FlagSet.Lookup(name)
		// if 'nil' subcommand is invalid
		//	- if time permits, can implement levenshtein distance to see closest subcommand
		//	- same check can be used on main Command
		if flag == nil {
			fmt.Errorf("Unrecognized subcommand: %s\n", name)
			os.Exit(1)
		}

		val := flag.Value.String()

		val_type := subcommand.DefVal.ValType
		

		// need to cast as returned value is just string
		switch val_type {
		case "int":
			int_val, err := strconv.Atoi(val)
			if err != nil {
				fmt.Errorf("Could not parse as an integer: %s\n", val)
			}
			// add to map if it is in range
			if intInRange(int_val , subcommand.ValRange) {
				// can provide just string, since it will become string in sending the command
				SubCmdVal[subcommand.Name] = val
			}

		case "string":
				SubCmdVal[subcommand.Name] = val
		case "bool":
				SubCmdVal[subcommand.Name] = val
		case "float64":
			// implement range value testing
			SubCmdVal[subcommand.Name] = val
		// TODO: implement Range type parsing
		case "range":
			// string -> Range
			split := strings.Split(val, "-")
			if len(split) < 2 {
				// should never get here
				fmt.Errorf("Range in incorrect form: %s\n", val)
				os.Exit(1)
			}

			if rangeInRange(
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
		ValType: "range",
		Val: Range {
			Lower: lower,
			Upper: upper,
		},
	}
}

func BoolVal(b bool) Value {
	return Value {
		ValType: "bool",
		Val: b,
	}
}

func IntVal(i int) Value {
	return Value {
		ValType: "int",
		Val: i,
	}
}

func StringVal(s string) Value {
	return Value {
		ValType: "string",
		Val: s,
	}
}



func Float64Val(f float64) Value {
	return Value {
		ValType: "float64",
		Val: f,
	}
}

func ParseCommand(c Command) {
	set := flag.NewFlagSet(c.Name, flag.ExitOnError)

	for _, o := range c.Subcommands{
		switch o.DefVal.ValType {
		case "int":
			n := o.DefVal.Val.(int)
			set.Int(o.Name, n , o.Usage)
		case "string":
			s := o.DefVal.Val.(string)
			set.String(o.Name, s, o.Usage)
		case "float64":
			f := o.DefVal.Val.(float64)
			set.Float64(o.Name, f, o.Usage)
		case "bool":
			b := o.DefVal.Val.(bool)
			set.Bool(o.Name, b, o.Usage)
		case "range":
			r := o.DefVal.Val.(Range)
			set.Var(&r, o.Name, o.Usage)
		default:
			os.Exit(1)
		}
	}
	// parsing now for testing, real usage parse ONLY ONCE ALL commands are processed

	ArgMap[c.Name] = &ArgFlag {
		FlagSet: set,
		command: &c,
	}
}

// add parameters for ranges
func CreateSubCmd(name string, usage string, val Value, minv int, maxv int, low_range int, up_range int) Subcommand {
	return Subcommand {
		ValRange: Range {
			Lower: 0,
			Upper: 0,
		},
		Name: name,
		Usage: usage,
		DefVal: val,
		MinMaxv: Range {
			Lower: minv,
			Upper: maxv,
		},
	}
}

func CreateCommand(name string, min_opts int, max_opts int, subcommands []Subcommand) Command {

	return Command {
		Name: name,
		MinSubCmds : min_opts,
		MaxSubCmds: max_opts,
		subcommands: subcommands,
	}
}
