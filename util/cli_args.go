package util

// usage format to add
// ./cli [COMMAND] --[SUBCOMMAND] | --[SUBCOMMAND] [VAL] | -[SUBCOMMAND] [VAL] | --[SUBCOMMAND]=[VAL] | (--verbose | -verbose)
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
	Description string
	MinSubCmds int // minimum options to provide
	MaxSubCmds int //
}

// rename
type Subcommand struct {
	// if wanted to do something in a range of 
	ValRange Range // range to operate on
	Name string // name of subcommand
	Usage string // description on how to use subcommand
	DefVal Value // default value
	// This program will keep [minv, maxv) as integer values wrapped in Range type
	MinMaxv Range
}

// Actual CLI arguments are given by Flag
type ArgFlag struct {
	FlagSet *flag.FlagSet
	command *Command
}

var usage = `
./cli [COMMAND] --[SUBCOMMAND] | --[SUBCOMMAND] [VAL] | -[SUBCOMMAND] [VAL] | --[SUBCOMMAND]=[VAL] | (--verbose | -verbose)
`

var examples = `python3 grip_cli.py ping
    cli ping
    cli ping ip=158.100.69.89
    cli ping ip=158.100.69.89 port=12350
    cli ping ip=158.100.69.89 verbose
    cli zoneset trip 3
`

// to be used and exported outside this package
// string being the command (ie. 'ble', 'camera', 'ls')
var ArgMap = make(map[string]*ArgFlag)

var SubCmdVal = make(map[string]string)

// COULD BE BETTER
const (
	EMPTY_STR = ""
	EMPTY_INT = 0
	EMPTY_F64 = 0.0

	MAX_I32 = 2<<31
)

func floatInRange(f float64, r Range) bool {
	return f >= float64(r.Lower) && f < float64(r.Upper)
}

func EmptyRange() Range {
	return Range {
		Lower: 0,
		Upper: 0,
	}
}

// very unlikely to be hit so good numbers to check
func RequiredRange() Range {
	return Range {
		Lower: 420,
		Upper: 69,
	}
}


// convert string in form of 'n-m' to Range object or panic
func stringToRange(s string) Range {
	split := strings.Split(s, "-")

	if len(split) < 2 {
		fmt.Printf("Not in form of a propper Range: %s\n", s)
		fmt.Printf("Hint: Range should be '1-10', '5-100', '4-5'\n")
		os.Exit(1)
	}

	lower,err := strconv.Atoi(split[0])
	upper, err:= strconv.Atoi(split[1])

	if err != nil {
		fmt.Printf("Not in form of a propper Range: %s\n", s)
		fmt.Printf("Hint: Range should be '1-10', '5-100', '4-5'\n")
		os.Exit(1)
	}

	return Range {
		Lower: lower,
		Upper: upper,
	}

}

// Satisfy flag.Value interface 'Set(string)-> error' and 'String()->string'
func (r *Range) Set(s string) error {
	val := stringToRange(s)
	r.Lower = val.Lower
	r.Upper = val.Upper
	// no error returned, program will exit with error
	return nil
}

func (r *Range) String() string {
	return fmt.Sprintf("%d-%d", r.Lower, r.Upper)
}


func isInRange(r1 Range, r2 Range) bool {
	return r1.Lower >= r2.Lower && r1.Upper < r2.Upper
}

func isEmptyRange(r Range) bool {
	return r.Lower == EmptyRange().Lower && r.Upper == EmptyRange().Upper
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
	af.FlagSet.Parse(os.Args[2:])
	command := af.command
	subcommands := af.command.Subcommands


	// 2 includes filename and command
	if len(os.Args) - 2 < command.MinSubCmds {
		fmt.Printf("Not enough subcommands provided. Minimum of %d, supplied %d\n", command.MinSubCmds, len(os.Args) - 2)
		fmt.Println("============================================================")
		PrintHelpCmd(os.Args[1])
		os.Exit(1)
	}
	for _, subcommand := range subcommands{
		name := subcommand.Name
		flag := af.FlagSet.Lookup(name)
		// if 'nil' subcommand is invalid
		//	- if time permits, can implement levenshtein distance to see closest subcommand
		//	- same check can be used on main Command
		if flag == nil {
			fmt.Printf("Unrecognized subcommand: %s\n", name)
			os.Exit(1)
		}

		val := flag.Value.String()

		val_type := subcommand.DefVal.ValType

		

		// need to cast as returned value is just string
		switch val_type {
		case "int":
			int_val, err := strconv.Atoi(val)
			if err != nil {
				fmt.Printf("Could not parse as an integer: %s\n", val)
				os.Exit(1)
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
			float, err := strconv.ParseFloat(val, 64)
			if err != nil  {
				fmt.Printf("Could not convert string to float: %s\n", val)
				os.Exit(1)
			}

			// implement range value testing
			r := subcommand.MinMaxv
			if !floatInRange(float, subcommand.MinMaxv) {
				fmt.Printf("Provided float is out of bounds: %f, should be between [%d-%d)\n", float, r.Lower, r.Upper)
				os.Exit(1)
			}
			SubCmdVal[subcommand.Name] = val
		// TODO: implement Range type parsing
		case "range":
			// string -> Range

			provided := stringToRange(val)
			fmt.Println(val)

			if isEmptyRange(provided) {
				continue
			}
			boundary := subcommand.MinMaxv

			if ! isInRange(provided, boundary) {
				fmt.Printf("Provided range exceeds boundaries of %d-%d (supplied %d-%d)\n", boundary.Lower, boundary.Upper, provided.Lower, provided.Upper)
				os.Exit(2)
			}
		default:
			fmt.Printf("Unrecognized type: %s\n", val_type)
			os.Exit(1)
		}
		
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
			fmt.Printf("not a recognized argument type\n")
			os.Exit(1)
		}
	}
	// parsing now for testing, real usage parse ONLY ONCE ALL commands are processed

	ArgMap[c.Name] = &ArgFlag {
		FlagSet: set,
		command: &c,
	}
}

func PrintHelpCmd(cmd string) {
	af := ArgMap[cmd]
	fmt.Println(af.command.Description)
	af.FlagSet.Usage()
}

func PrintHelpFull() {
	fmt.Println("cli usage: "+usage)
	for _, command := range ArgMap {
		command.FlagSet.Usage()
	}

	fmt.Println(examples)
}
