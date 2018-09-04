package main

import ("log"
	"nes-emu/hardware")

func main() {
	log.Println("Arte's NES Emu")
	// create cpu
	cpu := hardware.Cpu{}
	hardware.CPURunInstr([2]byte{32, 12})

	cart, err := hardware.CreateCartridge("donkey-kong.nes")

	if err != nil {
		log.Println(err)
	} else {
                cpu.LoadCartridge(cart)
		//log.Println(cart)
		log.Println(cpu.Memory)
		cpu.Reset()
	}
	
	// log.Printf("%x", rom)
	// log.Println(err)
}

