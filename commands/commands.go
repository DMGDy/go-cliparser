package commands

import (
	"fmt"
	"errors"
	"os"

	"github.com/DMGDy/grip2-cli/util"
	"github.com/DMGDy/grip2-cli/commands/run"
	"github.com/DMGDy/grip2-cli/commands/ble"
)

var Commands = map[string]run.RunCommand {
	"ble": &ble.Ble{},
}

func processCommand(c string, r run.RunCommand) {
		af := util.ArgMap[c]
		util.ValidateValues(af)

		// parse and validate
		af.FlagSet.Parse(os.Args[2:])

}

func RegisterCommands() {
	for _, cmd := range Commands {
		cmd.Register()
	}
}

// either return the response (can be empty or an error)
func RunCommand(command string) (string, error) {
	r, ok := Commands[command]
	if ok {
		processCommand(command, r)
		r.Run()
	} else {
		return " ", errors.New(fmt.Sprintf("Command `%s` does not exist\n"))
	}

	return "", errors.New("TODO")
}

