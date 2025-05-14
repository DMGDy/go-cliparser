package ble

import (
	"fmt"
	"errors"

	"github.com/DMGDy/grip2-cli/commands/run"
	"github.com/DMGDy/grip2-cli/util"
)

var _ run.RunCommand = (*Ble)(nil)

type Ble struct {
	subcmds map[string]util.Value
}

func (b *Ble) Run() (string, error) {
	for k,v := range util.SubCmdVal {
		fmt.Printf("%s: %s\n",k,v)
	}
	return "", errors.New("TODO")
}

// should assign this->subcmds
func (b *Ble) Register() {
	var subcmd []util.Subcommand

	// ble command begin
	bvalue := util.BoolVal(false)
	subcmd = append(subcmd, util.CreateSubCmd("scanon", "Turns BLE scan on.", bvalue, 0, 0, 0, 0))

	bvalue = util.BoolVal(false)
	subcmd = append(subcmd, util.CreateSubCmd("scanoff", "Turns BLE scan on.", bvalue, 0, 0, 0, 0))

	ivalue := util.IntVal(1000800)
	subcmd = append(subcmd, util.CreateSubCmd("uid", `This indicates which BLE
process to send this message to.`, ivalue, 0, 99999999, 0, 0))

	c := util.CreateCommand("ble", 1, 1, subcmd)

	util.ParseCommand(c)

}

