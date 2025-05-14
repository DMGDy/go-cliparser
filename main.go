package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/DMGDy/grip2-cli/util"
	"github.com/DMGDy/grip2-cli/commands"
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

	commands.RegisterCommands()

	if len(os.Args) < 2 {
		//print usage spiel
		fmt.Println("grip2-cli usage:")
		// generate all help messages
		for _, command := range util.ArgMap {
			command.FlagSet.Usage()
		}

		os.Exit(1)
	}


	_, err := commands.RunCommand(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
