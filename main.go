package main

import "log"
import "nes-emu/hardware"

func main() {
	log.Println("test log")
	hardware.CPURunInstr([2]byte{32, 12})

	log.Println(hardware.Instructions[0x00])
}

