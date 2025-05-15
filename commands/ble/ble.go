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
	MinSubCmds :1,
	MaxSubCmds: 5,
	Subcommands: []util.Subcommand {
		util.Subcommand {
			Name: "scanon",
			Usage: "turns BLE scan on",
			ValRange: util.Range {
				Lower: 0,
				Upper: 0,
			},
			MinMaxv: util.Range {
				Lower: 0,
				Upper: 0,
			},
			DefVal: util.Value {
				ValType: "bool",
				Val: false,
			},
		},
		util.Subcommand {
			Name: "scanoff",
			Usage: "turns BLE scan off",
			ValRange: util.Range {
				Lower: 0,
				Upper: 0,
			},
			MinMaxv: util.Range {
				Lower: 0,
				Upper: 0,
			},
			DefVal: util.Value {
				ValType: "bool",
				Val: false,
			},
		},
		util.Subcommand {
			Name: "uid",
			Usage: "BLE process to send message to (default to 1000800)",
			ValRange: util.Range {
				Lower: 0,
				Upper: 0,
			},
			MinMaxv: util.Range {
				Lower: 0,
				Upper: 0,
			},
			DefVal: util.Value {
				ValType: "string",
				Val: "1000800",
			},
		},
		
		util.Subcommand {
			Name: "list",
			Usage: "Shows a list of all learned devices",
			ValRange: util.Range {
				Lower: 0,
				Upper: 0,
			},
			MinMaxv: util.Range {
				Lower: 0,
				Upper: 0,
			},
			DefVal: util.Value {
				ValType: "string",
				Val: "N/A",
			},
		},

		util.Subcommand {
			Name: "delete",
			Usage: "Deletes the device of given uid shown from list",
			ValRange: util.Range {
				Lower: 0,
				Upper: 0,
			},
			MinMaxv: util.Range {
				Lower: 0,
				Upper: 0,
			},
			DefVal: util.Value {
				ValType: "string",
				Val: " ",
			},
		},

		util.Subcommand {
			Name: "pinack",
			Usage: "sends a pin acknowledge message",
			ValRange: util.Range {
				Lower: 0,
				Upper: 0,
			},
			MinMaxv: util.Range {
				Lower: 0,
				Upper: 0,
			},

			DefVal: util.Value {
				ValType: "string",
				Val: " ",
			},
		},

		util.Subcommand {
			Name: "range-example",
			Usage: "test if range implementation works",
			ValRange: util.Range {
				Lower: 0,
				Upper: 0,
			},
			MinMaxv: util.Range {
				Lower: 1,
				Upper: 5,
			},

			DefVal: util.Value {
				ValType: "range",
				Val: "0-1",
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

	util.ParseCommand(Cmd)
}

