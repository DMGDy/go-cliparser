package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"bitbucket.resideo.com/276733/grip2-cli/util"
	"bitbucket.resideo.com/276733/grip2-cli/commands"
	mqtt "bitbucket.resideo.com/276733/grip2-cli/mqtt-client"
)

// looking for literal argument of "verbose"

func consumeConfFlags(args *[]string) {
	a := slices.Index(*args, "-v")
	var b = 0 
	for i, s := range *args {
		if strings.Contains(s, "--ip=") {
			b = i
		}
	}

	if a > 0 {
		util.Verbose = true
		*args = slices.Delete(*args, a, a+1)
		fmt.Printf("**Verbose output enbabled**\n\n")
	}

	if b > 0 {
		parts := strings.Split((*args)[b], "=")
		if len(parts) < 2 {
			fmt.Println("Argument not provided correctly.")
			fmt.Println("--ip=[IP]")
			fmt.Println("ie.) --ip=10.1.195.139")
			os.Exit(1)
		}
		mqtt.IP = parts[1]

		*args = slices.Delete(*args, b, b+1)
		fmt.Printf("**Using IP: %s**\n\n" ,mqtt.IP)
	}
}

// consume "--verbose" "-v:" if provided
func main() {

	if len(os.Args) < 2 {
		util.PrintHelpFull()
		os.Exit(1)
	}

	consumeConfFlags(&os.Args)


	if len(os.Args) < 2 {
		util.PrintHelpFull()
		os.Exit(1)
	}

	// register the defined commands/subcommands to flag module
	commands.RegisterCommands()

	// init mqtt client
	err := mqtt.InitClient()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}



	err = commands.RunCommand(os.Args[1])


	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	mqtt.CloseClient()
}
