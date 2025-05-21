package ble

import (
	"fmt"
	"errors"

	"bitbucket.resideo.com/276733/grip2-cli/commands/run"
	"bitbucket.resideo.com/276733/grip2-cli/util"
)

var description = `
    Usage: scanon|scanoff|list|delete|pinack [uid]
    Perform various operation on the Bluetooth system.

    scanon turns BLE scan on.
    scanoff turns BLE scan off.
    pinack sends a pin acknowledge message.
    Those take an optional uid.  This indicates which BLE process to send this message to.
    1000800 is for the one on the panel, and is also the default if no uid is specified.
    1000801 is for tablet 1, 1000802 for tablet two etc (check that this is true??).

    list shows a list of all learned devices.

    delete deletes the device.  This requires the uid of the device, which is the ID shown
    from the list command.

    Examples:
      ./cli ble list
      ./cli ble scanon
      ./cli ble scanon 1000802
      ./cli ble delete 8001
`


type Ble struct {
}


var Cmd = util.Command {
	Name: "ble",
	MinSubCmds: 1,
	MaxSubCmds: 5,
	Description: description,
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
				Val: "",
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
				Val: "",
			},
		},

		util.Subcommand {
			Name: "pinack",
			Usage: "sends a pin acknowledge message",

			MinMaxv: util.Range {
				Lower: 0,
				Upper: 0,
			},

			DefVal: util.Value {
				ValType: "string",
				Val: "",
			},
		},

		util.Subcommand {
			Name: "range-example",
			Usage: "test if range implementation works",
			MinMaxv: util.Range {
				Lower: 1,
				Upper: 5,
			},

			DefVal: util.Value {
				ValType: "range",
				Val: util.RequiredRange(),
			},
		},
	},
}

// try to find a way to define elsewhere here
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

