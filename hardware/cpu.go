package hardware

import (
	"log"
	"math"
	"time"
)

// cpu speed
const cpuSpeed = 1789773
var NsPerCycle = (1 / float64(1789773)) * math.Pow10(9)

type Cpu struct {
	// pointer to base struct
	nes *NES

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

func (cpu *Cpu) setCpuInitialState() {
	// Set initial flags
	cpu.P = 0x24

	// PPU register initial state
	cpu.Memory[0x2000] = 0x80
	cpu.Memory[0x2002] = 0xA0

	// initialize the stack pointer
	cpu.SP = 0xFD
}

func (cpu *Cpu) runMainCpuLoop() {
	// number of instructions ran
	var numOfInstructions uint = 0

	for true {
		opcode := cpu.Read8(cpu.PC)
		cpu.RunInstruction(Instructions[opcode], false)

		time.Sleep(500 * time.Nanosecond)

		numOfInstructions++

		if numOfInstructions % 100 == 0 {
			if (cpu.Memory[0x2002] >> 7) & 1 == 0 {
				cpu.nes.PPU.SetVBlank()
			} else {
				cpu.nes.PPU.ClearVBlank()
			}
		}

		if numOfInstructions % 50000 == 0 {
			//log.Printf("%+v", cpu.nes.PPU.Memory[0x2000:0x2400])
		}
	}
}

func (cpu *Cpu) Reset() {
	// Read first instruction address location
	firstInstruction := cpu.Read16(0xFFFC)
	//firstInstruction := uint16(0xC000)

	// Set the PC to be at the address
	cpu.PC = firstInstruction

	cpu.setCpuInitialState()

	//cpu.runMainCpuLoop()

	// print the whole CPU and memory!!
	log.Printf("%+v\n", cpu)
}

func (cpu *Cpu) HandleNMI() {
	// Push current pc to the stack
	cpu.Push16(cpu.PC)

	// push current processor status to the stack
	cpu.Push8(cpu.P)

	// Set the PC to the NMI vector at 0xFFFA
	cpu.PC = cpu.Read16(0xFFFA)
}