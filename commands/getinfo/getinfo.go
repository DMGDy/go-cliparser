package getinfo

import (
	"fmt"
	"errors"
	"strconv"
	"encoding/json"

	"bitbucket.resideo.com/276733/grip2-cli/commands/run"
	"bitbucket.resideo.com/276733/grip2-cli/util"
	mqtt "bitbucket.resideo.com/276733/grip2-cli/mqtt-client"
)

var description = `
    Usage: getinfo delay
    Continuously arm (FullSet) and disarm partition one forever.  The
    delay field is the number of seconds to sleep between each
    arm/disarm operation.  The number can be less than 1 (for example
    0.5). 

	Providing a value to 'delay' will ignore the 'arm' and 'disarm' subcommands.
    Example:
      ./cli getinfo --delay=2.5
	  ./cli getinfo -arm
	  ./cli getinfo -disarm
`


type Getinfo struct {
}

var Cmd = util.Command {
	Name: "getinfo",
	MinSubCmds: 1,
	MaxSubCmds: 2,
	Description: description,
	Subcommands: []util.Subcommand {

		util.Subcommand {
			Name: "panel",
			Usage: "Get the panel status",
			MinMaxv: util.Range {
				Lower: 0,
				Upper: 9,
			},
			DefVal: util.Value {
				ValType: "int",
				Val: 999,
			},
		},

		util.Subcommand {
			Name: "partition",
			Usage: "Get the status of a partition",
			MinMaxv: util.Range {
				Lower: 0,
				Upper: 4,
			},
			DefVal: util.Value {
				ValType: "int",
				Val: 999,
			},

		},
	},
}

type PartitionInfo struct {
	Account string `json:"Account"`
	Name string `json:"Name"`
	Id int `json:"id"`
	Status map[string]interface{} `json:"_Status"`
	State string `json:"State"`
	SubState map[string]interface{} `json:"_SubState"`
}

type PanelInfo struct {
	Door []interface{} `json:"Door"`
	Name string `json:"Name"`
	Privilege []interface{} `json:"Privilege"`
	Type int `json:"_Type:"`
	Id int `json:"id"`
	Pin string `json:"pin"`
}

const (
)

var _ run.RunCommand = (*Getinfo)(nil)

func (a *Getinfo) Run() error {
	panel, err := strconv.ParseInt(util.SubCmdVal["panel"], 10, 32)
	partition, err := strconv.ParseInt(util.SubCmdVal["partition"], 10, 32)

	if err !=nil {
		return errors.New("Could not convert subcommand argument type")
	}
	var response string

	if panel < 9{
		sub := "@/Panel/User_/#"
		pub := fmt.Sprintf("@/GET/Panel/User_/%d", panel)
		response, err = mqtt.SendCommand(sub, pub, "")
		if err != nil {
			return err
		}
		// format JSON 

		var panel_info PanelInfo
		if err = json.Unmarshal([]byte(response), &panel_info); err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Panel Info")
		fmt.Println("Door:",panel_info.Door)
		fmt.Println("Name:",panel_info.Name)
		fmt.Println("Privilege:",panel_info.Privilege)
		fmt.Println("Type:",panel_info.Type)
		fmt.Println("Id:",panel_info.Id)
		fmt.Println("Pin:",panel_info.Pin)

	}

	if partition < 9{
		partition = partition

		//                  @/Panel/Partition_/1/#
		sub := fmt.Sprintf("@/Panel/Partition_/%d/#", partition)
		pub := fmt.Sprintf("@/GET/Panel/Partition_/*")

		response, err = mqtt.SendCommand(sub, pub, "")


		if err != nil {
			return err
		}

		var part_info PartitionInfo
		if err = json.Unmarshal([]byte(response), &part_info); err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Partition Info")
		fmt.Println("Account:",part_info.Account)
		fmt.Println("Name:",part_info.Name)
		fmt.Println("ID:",part_info.Id)
		fmt.Println("Status:")
		for k,v := range part_info.Status {
			fmt.Printf("    %s: %.0f\n",k,v)
		}
		fmt.Println("State:",part_info.State)
		fmt.Println("SubState:")
		for k,v := range part_info.SubState{
			fmt.Printf("    %s: %.0f\n",k,v)
		}
	}

	return nil
}

// should assign this->subcmds
func (a *Getinfo) Register() {

	util.ParseCommand(Cmd)
}

