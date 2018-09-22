package hardware

import (
	"log"
	"time"
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
	Memory [0x10000]byte
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

	// Set initial flags
	cpu.P = 0x04

	// PPU register initial state
	cpu.Memory[0x2002] = 0xA0

	// initialize the stack pointer
	cpu.SP = 0xFF

	for true {
		opcode := cpu.Read8(cpu.PC)
		cpu.RunInstruction(Instructions[opcode])

		time.Sleep(500 * time.Nanosecond)
	}

	// print the whole CPU and memory!!
	log.Printf("%+v\n", cpu)
}
