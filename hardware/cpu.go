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

func (cpu *Cpu) Reset() {
	// Read first instruction address location
	firstInstruction := cpu.Read16(0xFFFC)
	// Set the PC to be at the address
	cpu.PC = firstInstruction
	
	firstInstructionOpcode := cpu.Read8(firstInstruction)
	
	log.Printf("First Instruction is at address %x", firstInstruction)
	log.Printf("First Instruction has opcode %x", firstInstructionOpcode)
	log.Println("First opcode is ")
	log.Printf("%+v\n", Instructions[firstInstructionOpcode])

	log.Println(cpu.Memory[cpu.PC + 1])

	cpu.RunInstruction(Instructions[0x69])
	cpu.RunInstruction(Instructions[0x0A])
	log.Printf("%+v\n", cpu)
}
