package hardware

import (
	"log"
	"errors"
)

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

// setZero - Sets the overflow bit on the cpu
func (cpu *Cpu) setOverflow() {
	cpu.P |= 0x40
}

// clearOverflow - Clear the overflow bit on the cpu
func (cpu *Cpu) clearOverflow() {
	cpu.P &= 0xBF
}

// setInterrupt - Sets the Interrupt disable bit on the cpu
func (cpu *Cpu) setInterrupt() {
	cpu.P |= 0x04
}

// clearInterrupt - Clear the Interrupt disable bit on the cpu
func (cpu *Cpu) clearInterrupt() {
	cpu.P &= 0xFB
}

// setBreak
func (cpu *Cpu) setBreak() {
	cpu.P |= 0x10
}

// clearBreak
func (cpu *Cpu) clearBreak() {
	cpu.P &= 0xEF
}

// setDecimal
func (cpu *Cpu) setDecimal() {
	cpu.P |= 0x08
}

// clearDecimal
func (cpu *Cpu) clearDecimal() {
	cpu.P &= 0xF7
}

func (cpu *Cpu) doRelativeBranch(value uint8) {
	if value >= 0x80 {
		cpu.PC = cpu.PC + uint16(value) - 0x100
	} else {
		cpu.PC = cpu.PC + uint16(value)
	}
}

func (cpu *Cpu) RunInstruction(instr instruction) {
	log.Printf("%+v\n", instr)
	
	// Increment the PC
	cpu.PC++

	var addr uint16
	var value uint8

	switch instr.mode {
	case A:
		addr = 0
		value = cpu.A
	case imm:
		arg := cpu.Read8(cpu.PC)
		addr = cpu.PC
                value = arg
		cpu.PC++
	case zpg:
		arg := cpu.Read8(cpu.PC)
		addr = uint16(arg)
		value = cpu.Read8(addr)
		cpu.PC++
	case zpgX:
		arg := cpu.Read8(cpu.PC)
		addr = uint16(arg + cpu.X) & 0xFF // wrap around for X
		value = cpu.Read8(addr)
		cpu.PC++
	case zpgY:
		arg := cpu.Read8(cpu.PC)
		addr = uint16(arg + cpu.Y) & 0xFF // wrap around for Y
		value = cpu.Read8(addr)
		cpu.PC++
	case abs:
		arg := cpu.Read16(cpu.PC)
		addr = arg
		value = cpu.Read8(addr)
		cpu.PC += 2
	case absX:
                arg := cpu.Read16(cpu.PC)
		addr = arg + uint16(cpu.X)
		value = cpu.Read8(addr)
		cpu.PC += 2
	case absY:
		arg := cpu.Read16(cpu.PC)
		addr = arg + uint16(cpu.Y)
		value = cpu.Read8(addr)
		cpu.PC += 2
	case ind:
		arg := cpu.Read16(cpu.PC)
		addr = cpu.Read16(arg)
		value = cpu.Read8(addr)
		cpu.PC += 2
	case indX:
		arg := cpu.Read8(cpu.PC)
		addr = cpu.Read16(uint16(arg + cpu.X) & 0xFF)
		value = cpu.Read8(addr)
		cpu.PC++
	case indY:
		arg := cpu.Read8(cpu.PC)
		addr = cpu.Read16(uint16(arg)) + uint16(cpu.Y)
		value = cpu.Read8(addr)
		cpu.PC++
	case rel:
		arg := cpu.Read8(cpu.PC)
		addr = 0
		value = arg
		cpu.PC++
	case impl:
		addr = 0
                value = 0
	default:
		log.Fatal(errors.New("Fatal: " + string(instr.mode) + " is not a valid addressing mode."))
	}

	switch instr.assemblyCode {
	case "ADC":
		cpu.ADC(instr, addr, value)
	case "CPY":
		cpu.CPY(instr, addr, value)
	case "AND":
		cpu.AND(instr, addr, value)
	case "ASL":
		cpu.ASL(instr, addr, value)
	case "BCC":
		cpu.BCC(instr, addr, value)
	case "BCS":
		cpu.BCS(instr, addr, value)
	case "BEQ":
		cpu.BEQ(instr, addr, value)
	case "BIT":
		cpu.BIT(instr, addr, value)
	case "SEI":
		cpu.SEI(instr, addr, value)
	default:
                log.Fatal(errors.New("Fatal: " + string(instr.assemblyCode) + " is not a valid instruction code."))
	}
}

// ADC - Add with Carry
// Performs addition with the accumulator and carry bit.
// Sets flags accordingly ZCN
func (cpu *Cpu) ADC(instr instruction, addr uint16, value uint8) {
	log.Println("ADC")

	// Calculate the result
	result := cpu.A + value + getBit(cpu.P, 7)
	
	// Set the carry flag if unsigned overflow occurs
	if uint16(cpu.A) + uint16(value) + uint16(getBit(cpu.P, 7)) > 0xFF {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}
	
	// Set the overflow flag
	if ((cpu.A ^ result) & (value ^ result) & 0x80) > 0 {
		cpu.setOverflow()
	} else {
		cpu.clearOverflow()
	}

	// Set the zero flag if the result is zero
	if result == 0 {
		cpu.setZero()
	} else {
		cpu.clearZero()
	}

	// Set the sign bit if bit 7 is 1
	if getBit(result, 7) > 0 {
		cpu.setNegative()
	} else {
		cpu.clearNegative()
	}

	// Set accumulator value
	cpu.A = result
}

// AND - Logical And
// Performs a logical and on the acc
// sets flags ZN
func (cpu *Cpu) AND(instr instruction, addr uint16, value uint8) {
        // Calculate the result
	result := cpu.A & value

	// Set the zero flag if the result is zero
	if result == 0 {
		cpu.setZero()
	} else {
		cpu.clearZero()
	}

	// Set the sign bit if bit 7 is 1
	if (result & 0x80) > 0 {
		cpu.setNegative()
	} else {
		cpu.clearNegative()
	}

	cpu.A = result
}

// ASL - Arithmetic Shift Left
// Shifts bits in the acc or the memory one bit left
// sets flags ZCN
func (cpu *Cpu) ASL(instr instruction, addr uint16, value uint8) {
	var result uint8
	oldBit7 := getBit(value, 7)
	var newBit7 uint8

	// If acc mode, shift the acc
	if instr.mode == A {
                result = cpu.A << 1
		cpu.A = result
	} else {
		result = value << 1
		cpu.Memory[addr] = result
	}

	newBit7 = getBit(result, 7)
	
	// Set the carry flag old bit 7 is 1
	if oldBit7 > 0 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	// Set the zero flag if the result is zero
	if cpu.A == 0 {
		cpu.setZero()
	} else {
		cpu.clearZero()
	}

	// Set the sign bit if new bit 7 is 1
	if newBit7 > 0 {
		cpu.setNegative()
	} else {
		cpu.clearNegative()
	}
	
}

// BCC - Branch if Carry Clear
// If carry bit is clear, cause a relative branch to occur
func (cpu *Cpu) BCC(instr instruction, addr uint16, value uint8) {
	if getBit(cpu.P, 0) == 0 {
		cpu.doRelativeBranch(value)
	}
}

// BCS - Branch if Carry Set
// If carry bit is set, cause a relative branch to occur
func (cpu *Cpu) BCS(instr instruction, addr uint16, value uint8) {
	if getBit(cpu.P, 0) > 0 {
		cpu.doRelativeBranch(value)
	}
}

// BEQ - Branch if Equal
// If zero flag is set, cause relative branch
func (cpu *Cpu) BEQ(instr instruction, addr uint16, value uint8) {
	if getBit(cpu.P, 1) > 0 {
		cpu.doRelativeBranch(value)
	}
}

// BIT - Bit Test
// Ands Acc with mem location but does not store
//  Sets flags
func (cpu *Cpu) BIT(instr instruction, addr uint16, value uint8) {
	result := cpu.A & value

	// Set zero flag if result is 0
        if result == 0 {
		cpu.setZero()
	} else {
		cpu.clearZero()
	}

	// Set Overflow to bit 6
	if getBit(result, 6) == 1 {
		cpu.setOverflow()
	} else {
		cpu.clearOverflow()
	}

	// Set the sign bit if bit 7 is 1
	if getBit(result, 7) > 0 {
		cpu.setNegative()
	} else {
		cpu.clearNegative()
	}
}

// BMI - Branch if Minus
// If the negative flag is set, then do a relative branch
func (cpu *Cpu) BMI(instr instruction, addr uint16, value uint8) {
	// check if negative flag is set
	if getBit(cpu.P, 7) > 0 {
		cpu.doRelativeBranch(value)
	}
}

// BNE - Branch Not Equal
// If the zero flag is clear, then do a relative branch
func (cpu *Cpu) BNE(instr instruction, addr uint16, value uint8) {
	// check if negative flag is set
	if getBit(cpu.P, 1) == 0 {
		cpu.doRelativeBranch(value)
	}
}

// BPL - Branch if Positive
// If the negative flag is clear, then do a relative branch
func (cpu *Cpu) BPL(instr instruction, addr uint16, value uint8) {
	// check if negative flag is set
	if getBit(cpu.P, 7) > 0 {
		cpu.doRelativeBranch(value)
	}
}

// BRK - Force Interrupt
// forces generation of an interrupt request
// sets break command flag
func (cpu *Cpu) BRK(instr instruction, addr uint16, value uint8) {
	// push program counter and processor status to stack
	cpu.Push16(cpu.PC)

	cpu.setBreak()
	cpu.Push8(cpu.P)

	cpu.setInterrupt()
	
	cpu.PC = cpu.Read16(0xFFFE)
}

// BVC - Branch if overflow clear
func (cpu *Cpu) BVC(instr instruction, addr uint16, value uint8) {
	if getBit(cpu.P, 6) == 0 {
		cpu.doRelativeBranch(value)
	}
}

// BVS - Branch if overflow clear
func (cpu *Cpu) BVS(instr instruction, addr uint16, value uint8) {
	if getBit(cpu.P, 6) > 0 {
		cpu.doRelativeBranch(value)
	}
}

// CLC - set the carry flag to zero
func (cpu *Cpu) CLC(instr instruction, addr uint16, value uint8) {
	cpu.clearCarry()
}

// CLD - clear decimal mode
func (cpu *Cpu) CLD(instr instruction, addr uint16, value uint8) {
	cpu.clearDecimal()
}

// CLI - clear interrupt disable flag
func (cpu *Cpu) CLI(instr instruction, addr uint16, value uint8) {
	cpu.clearInterrupt()
}

// CLV - clear the overflow flag
func (cpu *Cpu) CLV(instr instruction, addr uint16, value uint8) {
	cpu.clearOverflow()
}

func (cpu *Cpu) cmpVals(x, y uint8) {
	// Compute the result
        compareResult := x - y

	// Set the carry flag if Y >= value
	if x >= y {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	// Set the zero flag if the result is zero
	if x == y {
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

// CMP - compares acc with value and sets flags
func (cpu *Cpu) CMP(instr instruction, addr uint16, value uint8) {
	cpu.cmpVals(cpu.A, value)
}

// CPX - compares x register with value and sets flags
func (cpu *Cpu) CPX(instr instruction, addr uint16, value uint8) {
        cpu.cmpVals(cpu.X, value)
}

// CPY - Compare Y Register
// Performs a subtraction on Y register and src.
// Sets flags accordingly. ZCN
func (cpu *Cpu) CPY(instr instruction, addr uint16, value uint8) {
	cpu.cmpVals(cpu.Y, value)
}

func (cpu *Cpu) setZHelper(value uint8) {
	if value == 0 {
		cpu.setZero()
	} else {
		cpu.clearZero()
	}
}

func (cpu *Cpu) setNHelper(value uint8) {
	if (value & 0x80) > 0 {
		cpu.setNegative()
	} else {
		cpu.clearNegative()
	}
}

// DEC - Decrement from memory
func (cpu *Cpu) DEC(instr instruction, addr uint16, value uint8) {
	result := value - 1
        cpu.Memory[addr] = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// DEX - Decrement from X register
func (cpu *Cpu) DEX(instr instruction, addr uint16, value uint8) {
	result := value - 1
        cpu.X = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// DEY - Decrement from Y register
func (cpu *Cpu) DEY(instr instruction, addr uint16, value uint8) {
	result := value - 1
        cpu.Y = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// EOR - Exclusive OR with acc
func (cpu *Cpu) EOR(instr instruction, addr uint16, value uint8) {
	result := cpu.A ^ value
        cpu.A = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// INC - Increment from memory
func (cpu *Cpu) INC(instr instruction, addr uint16, value uint8) {
	result := value + 1
        cpu.Memory[addr] = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// INX - Increment from X register
func (cpu *Cpu) INX(instr instruction, addr uint16, value uint8) {
	result := value + 1
        cpu.X = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// INY - Increment from X register
func (cpu *Cpu) INY(instr instruction, addr uint16, value uint8) {
	result := value + 1
        cpu.Y = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// JMP - Increment from X register
func (cpu *Cpu) JMP(instr instruction, addr uint16, value uint8) {
        cpu.PC = addr
}

// JSR - Increment from X register
func (cpu *Cpu) JSR(instr instruction, addr uint16, value uint8) {
	result := value + 1
        cpu.Y = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// SEI - Set Interrupt Disable
// Sets the interrupt disable flag to one
func (cpu *Cpu) SEI(instr instruction, addr uint16, value uint8) {
	log.Println("SEI")
	
	cpu.setInterrupt()
}
