package hardware

import "log"

func getBit(i uint8, pos uint) uint8 {
	return (i >> pos) & 1
}

// setCarry - Sets the carry bit on the cpu
func (cpu *Cpu) setCarry() {
	cpu.P |= 0x01
}

// clearCarry - Clears the carry bit on the cpu
func (cpu *Cpu) clearCarry() {
	cpu.P &= 0xFE
}

// setNegative - Sets the negative bit on the cpu
func (cpu *Cpu) setNegative() {
	cpu.P |= 0x80
}

// clearNegative - Clears the negative bit on the cpu
func (cpu *Cpu) clearNegative() {
	cpu.P &= 0x7F
}

// setZero - Sets the zero bit on the cpu
func (cpu *Cpu) setZero() {
	cpu.P |= 0x02
}

// clearZero - Clear the zero bit on the cpu
func (cpu *Cpu) clearZero() {
	cpu.P &= 0xFD
}

// ADC - Add with Carry
// Performs addition with the accumulator and carry bit.
// Sets flags accordingly
func (cpu *Cpu) ADC(instr instruction) {
	log.Println("ADC")
	log.Printf("%+v\n", instr)

	// Increase the PC
	cpu.PC++

	var value uint8

	switch instr.mode {
	case imm:
		arg := cpu.Read8(cpu.PC)
                value = arg
		cpu.PC++
	case zpg:
		arg := cpu.Read8(cpu.PC)
		value = cpu.Memory[arg]
		cpu.PC++
	case zpgX:
		arg := cpu.Read8(cpu.PC)
		value = cpu.Memory[uint16(arg + cpu.X) & 0xFF] // wrap around for X
		cpu.PC++
	case abs:
		arg := cpu.Read16(cpu.PC)
		value = cpu.Memory[arg]
		cpu.PC += 2
	case absX:
                arg := cpu.Read16(cpu.PC)
		value = cpu.Memory[arg + uint16(cpu.X)]
		cpu.PC += 2
	case absY:
		arg := cpu.Read16(cpu.PC)
		value = cpu.Memory[arg + uint16(cpu.Y)]
		cpu.PC += 2
	case indX:
		arg := cpu.Read8(cpu.PC)
		indirectAddress := cpu.Read16(uint16(cpu.Memory[uint16(arg + cpu.X) & 0xFF]))
		value = cpu.Memory[indirectAddress]
		cpu.PC++
	case indY:
		arg := cpu.Read8(cpu.PC)
		indirectAddress := cpu.Read16(uint16(cpu.Memory[uint16(arg)])) + uint16(cpu.Y)
		value = cpu.Memory[indirectAddress]
		cpu.PC++
	}

	// Set accumulator value
	cpu.A = cpu.A + value + getBit(cpu.P, 7)

	
	
}

// CPY - Compare Y Register
// Performs a subtraction on Y register and src.
// Sets flags accordingly.
func (cpu *Cpu) CPY(instr instruction) {
	log.Println("CPY")
        log.Printf("%+v\n", instr)

	// increase PC
        cpu.PC++

	var value uint8

	// Figure out addressing mode
	switch instr.mode {
	case imm:
		arg := cpu.Read8(cpu.PC)
		value = arg
		cpu.PC++
	case zpg:
		arg := cpu.Read8(cpu.PC)
		value = cpu.Memory[arg]
		cpu.PC++
	case abs:
		arg := cpu.Read16(cpu.PC)
		log.Printf("Reading %x", arg)
		value = cpu.Memory[arg]
		cpu.PC += 2
	}

	// Compute the result
        compareResult := cpu.Y - value

	// Set the carry flag if Y >= value
	if cpu.Y >= value {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	// Set the zero flag if the result is zero
	if cpu.Y == value {
		cpu.setZero()
	} else {
		cpu.clearZero()
	}

	// Set the sign bit if bit 7 is 1
	if (compareResult & 0x80) > 0 {
		cpu.setNegative()
	} else {
		cpu.clearNegative()
	}
}
