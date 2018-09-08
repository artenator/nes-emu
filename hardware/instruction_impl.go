package hardware

import "log"

func (cpu *Cpu) CPY(instr instruction) {
	log.Println("CPY")
	log.Println(instr)
}