package ble

import (
	"fmt"
	"errors"

	"github.com/DMGDy/grip2-cli/commands/run"
	"github.com/DMGDy/grip2-cli/util"
)



type Ble struct {
}

var Cmd = util.Command {
	Name: "ble",
	min_opts: 1,
	max_opts: 5,
	subcommands: []util.Subcommands {
		Subcommand {
			Name: "scanon",
			usage: "turns BLE scan on",
			val_range: Range {
				lower: 0,
				upper: 0,
			},
			minmaxv: Range {
				lower: 0,
				upper: 0,
			},
			defval: Value {
				val_type: "bool",
				val: false,
			},
		},
		Subcommand {
			name: "scanoff",
			usage: "turns BLE scan off",
			val_range: Range {
				lower: 0,
				upper: 0,
			},
			minmaxv: Range {
				lower: 0,
				upper: 0,
			},
			defval: Value {
				val_type: "bool",
				val: false,
			},
		},
		Subcommand {
			name: "uid",
			usage: "BLE process to send message to (default to 1000800)",
			val_range: Range {
				lower: 0,
				upper: 0,
			},
			minmaxv: Range {
				lower: 0,
				upper: 0,
			},
			defval: Value {
				val_type: "string",
				val: 1000800,
			},
		},
		
		Subcommand {
			name: "list",
			usage: "Shows a list of all learned devices",
			val_range: Range {
				lower: 0,
				upper: 0,
			},
			minmaxv: Range {
				lower: 0,
				upper: 0,
			},
			defval: Value {
				val_type: "string",
				val: "N/A",
			},
		},

		Subcommand {
			name: "delete",
			usage: "Deletes the device of given uid shown from list",
			val_range: Range {
				lower: 0,
				upper: 0,
			},
			minmaxv: Range {
				lower: 0,
				upper: 0,
			},
			defval: Value {
				val_type: "string",
				val: " ",
			},
		},

		Subcommand {
			name: "pinack",
			usage: "sends a pin acknowledge message",
			val_range: Range {
				lower: 0,
				upper: 0,
			},
			minmaxv: Range {
				lower: 0,
				upper: 0,
			},

			defval: Value {
				val_type: "string",
				val: " ",
			},
		},

		Subcommand {
			name: "range-example",
			usage: "test if range implementation works",
			val_range: Range {
				lower: 0,
				upper: 0,
			},
			minmaxv: Range {
				lower: 0,
				upper: 0,
			},

			defval: Value {
				val_type: "string",
				val: " ",
			},
		},
	},
}

var _ run.RunCommand = (*Ble)(nil)

func (b *Ble) Run() (string, error) {
	for k,v := range util.SubCmdVal {
		fmt.Printf("%s: %s\n",k,v)
	}
	return "", errors.New("TODO")
}

// should assign this->subcmds
func (b *Ble) Register() {
	var subcmd []util.Subcommand

	util.ParseCommand(cmd)
}

