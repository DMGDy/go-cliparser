package main

import (
	"fmt"
	"os"
	"slices"

	"bitbucket.resideo.com/276733/grip2-cli/util"
	"bitbucket.resideo.com/276733/grip2-cli/commands"
	mqtt "bitbucket.resideo.com/276733/grip2-cli/mqtt-client"
)

// looking for literal argument of "verbose"
var verbose = false

func consumeVerbose(args *[]string) {
	n := -1
	a := slices.Index(*args, "--verbose")
	b := slices.Index(*args, "-v")
	if a > b {
		n = a
	} else {
		n = b
	}
	if n != -1 {
		// Essentially "popping" element from slice
		*args = slices.Delete(*args, n, n+1)
		fmt.Printf("**Verbose output enbabled**\n\n")
	}
}

// consume "--verbose" "-v:" if provided
func main() {

	// init mqtt client
	err := mqtt.InitClient()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}


	commands.RegisterCommands()

	if len(os.Args) < 2 {
		util.PrintHelpFull()
		os.Exit(1)
	}


	_, err = commands.RunCommand(os.Args[1])

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	mqtt.CloseClient()
}
