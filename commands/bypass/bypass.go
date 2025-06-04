package bypass

import (
	"fmt"
	"errors"
	"strconv"

	"bitbucket.resideo.com/276733/grip2-cli/commands/run"
	"bitbucket.resideo.com/276733/grip2-cli/util"
	mqtt "bitbucket.resideo.com/276733/grip2-cli/mqtt-client"
)

var description = `
    Usage: bypass (on|off) zone# user_code
    Bypass a zone.
    Example:
		./cli bypass --[on|off] --zone=[zone#] --user_code=[user_code#]
		./cli bypass on 1 1234
`


type Bypass struct {
}

var Cmd = util.Command {
	Name: "bypass",
	MinSubCmds: 3,
	MaxSubCmds: 3,
	Description: description,
	Subcommands: []util.Subcommand {

		util.Subcommand {
			Name: "on",
			Usage: "Turn on bypass for the given zone",
			DefVal: util.Value {
				ValType: "bool",
				Val: false,
			},
		},

		util.Subcommand {
			Name: "off",
			Usage: "Turn off bypass for the given zone",
			DefVal: util.Value {
				ValType: "bool",
				Val: false,
			},
		},

		util.Subcommand {
			Name: "zone",
			Usage: "Zone number to turn on/off bypass",
			DefVal: util.Value {
				ValType: "int",
				Val: util.EMPTY_INT,
			},
		},

		// what is user code???
		util.Subcommand {
			Name: "user_code",
			Usage: "The user code ",
			DefVal: util.Value {
				ValType: "int",
				Val: 1234,
			},
		},
	},
}

const (
	SUB_TOPIC = "@/Panel/Partition_/#"
)

// try to find a way to define elsewhere here
var _ run.RunCommand = (*Bypass)(nil)

func (b *Bypass) Run() error {
	topic := "@/SET/Panel/Zone_/*"
var mode  = 0

	on, err:= strconv.ParseBool(util.SubCmdVal["on"])
	off, err:= strconv.ParseBool(util.SubCmdVal["off"])
	if on {
		mode = 1
	} else if off {
		mode = 0
	}
	// can keep as string, integer
	zone := util.SubCmdVal["zone"]
	user_code := util.SubCmdVal["user_code"]

	payload := fmt.Sprintf(`{"=id":%s,"$inp":"Local","$utf":"8>1252","Bypass":%d,"$pin":"%s"}`,zone, mode, user_code)
	response, err := mqtt.SendCommand(SUB_TOPIC, topic, payload)

	if err != nil {
		return errors.New("Error sending MQTT Command")
	}

	fmt.Println(response)

	return nil
}

// should assign this->subcmds
func (b *Bypass) Register() {

	util.ParseCommand(Cmd)
}

