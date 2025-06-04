package commands

import (
	"fmt"
	"errors"
	"os"

	"bitbucket.resideo.com/276733/grip2-cli/util"
	"bitbucket.resideo.com/276733/grip2-cli/commands/run"
	"bitbucket.resideo.com/276733/grip2-cli/commands/ble"
	"bitbucket.resideo.com/276733/grip2-cli/commands/armdisarm"
	"bitbucket.resideo.com/276733/grip2-cli/commands/bypass"
	"bitbucket.resideo.com/276733/grip2-cli/commands/getinfo"
)

var Commands = map[string]run.RunCommand {
	"ble": &ble.Ble{},
	"armdisarm": &armdisarm.Armdisarm{},
	"bypass": &bypass.Bypass{},
	"getinfo": &getinfo.Getinfo{},
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
func RunCommand(command string) error {
	r, ok := Commands[command]
	if ok {
		processCommand(command, r)
		err := r.Run()
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("Command `%s` does not exist\n", command))
	}

	return nil
}

