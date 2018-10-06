package main

import (
	"log"
	"nes-emu/hardware"
)

func main() {
	log.Println("Arte's NES Emu")

	// create new nes
	//cpu := hardware.Cpu{}
	nes := hardware.NewNES()

	cart, err := hardware.CreateCartridge("nestest.nes")

	if err != nil {
		log.Println(err)
	} else {
		nes.LoadCartridge(cart)
		//log.Println(cart)
		//log.Println(cpu.Memory)
		nes.CPU.Reset()
	}
	
	// log.Printf("%x", rom)
	// log.Println(err)
}

