package armdisarm

import (
	"fmt"
	"errors"

	"github.com/DMGDy/grip2-cli/commands/run"
	"github.com/DMGDy/grip2-cli/util"
)

var description = `
    Usage: armdisarm delay
    Continuously arm (FullSet) and disarm partition one forever.  The
    delay field is the number of seconds to sleep between each
    arm/disarm operation.  The number can be less than 1 (for example
    0.5).
    Example:
      ./cli armdisarm 2.5
`


type Armdisarm struct {
}


		util.Subcommand {
var Cmd = util.Command {
	Name: "armdisarm",
	MinSubCmds: 1,
	MaxSubCmds: 1,
	Description: description,
	Subcommands: []util.Subcommand {
		util.Subcommand {
			Name: delay,
			Usage: "Number of seconds to sleep between each arm/disarm operation",
			DefValue: Value {
				Type: "float64",
				Val: util.EMPTY_F64
			},
			MinMaxv: Range {
				// delay time cannot be negative
				Lower: 0;
				Upper: util.EMPTY_INT
			},
		},
	},
}

// try to find a way to define elsewhere here
var _ run.RunCommand = (*ArmDisarm)(nil)

func (a *Armdisarm) Run() (string, error) {
	for k,v := range util.SubCmdVal {
		fmt.Printf("%s: %s\n",k,v)
	}
	return "", errors.New("TODO")
}

// should assign this->subcmds
func (a *Armdisarm) Register() {

	util.ParseCommand(Cmd)
}

