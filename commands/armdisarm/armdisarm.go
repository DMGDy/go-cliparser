package armdisarm

import (
	"fmt"
	"errors"
	"strconv"
	"time"

	"bitbucket.resideo.com/276733/grip2-cli/commands/run"
	"bitbucket.resideo.com/276733/grip2-cli/util"
	mqtt "bitbucket.resideo.com/276733/grip2-cli/mqtt-client"
)

var description = `
    Usage: armdisarm delay
    Continuously arm (FullSet) and disarm partition one forever.  The
    delay field is the number of seconds to sleep between each
    arm/disarm operation.  The number can be less than 1 (for example
    0.5). 

	Providing a value to 'delay' will ignore the 'arm' and 'disarm' subcommands.
    Example:
      ./cli armdisarm --delay=2.5
	  ./cli armdisarm -arm
	  ./cli armdisarm -disarm
`


type Armdisarm struct {
}


var Cmd = util.Command {
	Name: "armdisarm",
	MinSubCmds: 1,
	MaxSubCmds: 1,
	Description: description,
	Subcommands: []util.Subcommand {
		
		util.Subcommand {
			Name: "delay",
			Usage: "Number of seconds to sleep between each arm/disarm operation",
			DefVal: util.Value {
				ValType: "float64",
				Val: util.EMPTY_F64,
			},
			MinMaxv: util.Range {
				// delay time cannot be negative
				Lower: 0,
				Upper: util.MAX_I32,
			},
		},

		util.Subcommand {
			Name: "arm",
			Usage: "Arm Partition 1 once",
			DefVal: util.Value {
				ValType: "bool",
				Val: false,
			},
		},

		util.Subcommand {
			Name: "disarm",
			Usage: "Disarm Partition 1 once",
			DefVal: util.Value {
				ValType: "bool",
				Val: false,
			},
		},

	},
}

const (
	SUB_TOPIC = "@/Panel/Partition_/#"
)

var _ run.RunCommand = (*Armdisarm)(nil)

// what is difference between 'ArmAway' and 'ArmStay'?

func (a *Armdisarm) Run() error {
	topic := "@/Panel/Partition_/1"

	delay, err := strconv.ParseFloat(util.SubCmdVal["delay"], 32)
	if err != nil {
		return errors.New("Could not convert subcommand argument type")
	}
	arm, err := strconv.ParseBool(util.SubCmdVal["arm"])
	disarm ,err := strconv.ParseBool(util.SubCmdVal["disarm"])

	if err != nil {
		return errors.New("Could not convert subcommand argument type")
	}

	// ignore arm disarm if greater than 0
	if delay > 0.0 {
		for {
			payload := `{"SoundUI":[0,128,0,0,"0,B3",0,0],"State":"ArmStay","_SubState":{"ExitTimer":30}}`
			response, err:= mqtt.SendCommand(SUB_TOPIC, topic, payload)

			if err != nil {
				return errors.New("Error sending MQTT Command")
			}
			fmt.Println(response)

			time.Sleep(time.Millisecond * time.Duration(delay * 1000))

			payload = `{"SoundUI":[0,0,0,0,"0,B1",1,0],"State":"Disarmed","_SubState":{"ExitTimer":0}}`
			response, err = mqtt.SendCommand(SUB_TOPIC, topic, payload)

			if err != nil {
				return errors.New("Error sending MQTT Command")
			}
			fmt.Println(response)
		}
	} else if arm {
		payload := `{"SoundUI":[0,128,0,0,"0,B3",0,0],"State":"ArmStay","_SubState":{"ExitTimer":30}}`
		response, err:= mqtt.SendCommand(SUB_TOPIC, topic, payload)

		if err != nil {
			return errors.New("Error sending MQTT Command")
		}
		fmt.Println(response)

	} else if disarm {
		payload := `{"SoundUI":[0,0,0,0,"0,B1",1,0],"State":"Disarmed","_SubState":{"ExitTimer":0}}`
		response, err := mqtt.SendCommand(SUB_TOPIC, topic, payload)

		if err != nil {
			return errors.New("Error sending MQTT Command")
		}
		fmt.Println(response)

	}
	
	/*
	for k,v := range util.SubCmdVal {
		fmt.Printf("%s: %s\n",k,v)
	}
	*/
	return nil
}

// should assign this->subcmds
func (a *Armdisarm) Register() {

	util.ParseCommand(Cmd)
}

