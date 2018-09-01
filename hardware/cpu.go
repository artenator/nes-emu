package hardware

import (
	"log"
)

// cpu speed
const cpuSpeed = 1789773

type Cpu struct {
	// Program Counter
	PC uint16

	// Stack Pointer
	SP uint8

	// Accumulator
	A uint8

	// Index Register X
	X uint8

	// Index Register Y
	Y uint8

	// Processor Status
	// 7 6 5 4 3 2 1 0
	// N V   B D I Z C
	P uint8
}

// CPURunInstr Runs cpu instruction
func CPURunInstr(instr [2]byte) {
	opcode := instr[0]

	log.Printf("Received instruction %d\n", int(opcode))

	log.Println(Instructions[opcode])

}
