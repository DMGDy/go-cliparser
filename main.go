package main

import (
	"fmt"
	util "grip2-clu/util"
)

func main() {
	fmt.Printf("Hello, World!\n")

	fmt.Println("Hello, World!")
	var opts []Option

	bvalue := util.BoolVal(false)
	opts = append(opts, util.CreateOpts("scanon", "Turns BLE scan on.", bvalue, 0, 0))

	bvalue = BoolVal(true)
	opts = append(opts, util.CreateOpts("scanoff", "Turns BLE scan on.", bvalue, 0, 0))

	ivalue := IntVal(1000800)
	opts = append(opts, util.CreateOpts("uuid", "This indicates which BLE process to send this message to.", ivalue, 0, 99999999))

	command := util.CreateCommand("ble", 1, 1, opts)

	util.ParseCommand(command)

	af := util.ArgMap["ble"]

	flag := af.flag_set.Lookup("scanon")
	fmt.Println(flag.Value)
}
