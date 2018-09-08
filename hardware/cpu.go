package hardware

import (
	"log"
	"encoding/binary"
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

	// Memory
	Memory [0xFFFF]byte
}

// CPURunInstr Runs cpu instruction
func CPURunInstr(instr [2]byte) {
	opcode := instr[0]

	log.Printf("Received instruction %d\n", int(opcode))

	log.Println(Instructions[opcode])

}

func (cpu *Cpu) Reset() {
	// Read first instruction address location
	firstInstruction := binary.LittleEndian.Uint16(cpu.Memory[0xFFFC:0xFFFE])
	// Set the PC to be at the address
	cpu.PC = firstInstruction
	
	firstInstructionOpcode := cpu.Read8(firstInstruction)
	
	log.Printf("First Instruction is at address %x", firstInstruction)
	log.Printf("First Instruction has opcode %x", firstInstructionOpcode)
	log.Println("First opcode is ")
	log.Printf("%+v\n", Instructions[firstInstructionOpcode])

	cpu.CPY(Instructions[0xCD])
	log.Printf("%+v\n", cpu)
}
