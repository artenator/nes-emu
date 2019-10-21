package hardware

import (
	"errors"
	"log"
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

func (cpu *Cpu) RunInstruction(instr instruction, doLog bool) {
	var instrByteArr []byte
	for i := uint8(0); i < instr.bytes; i++ {
		instrByteArr = append(instrByteArr, cpu.Read8(cpu.PC + uint16(i)))
	}
	if doLog {
		log.Printf("%x %+v %x PC:%x A: %x SP: %x X: %x Y: %x P: %x PPUADDR: %x PPUDATA: %x",
			cpu.PC,
			instr,
			instrByteArr,
			cpu.PC,
			cpu.A,
			cpu.SP,
			cpu.X,
			cpu.Y,
			cpu.P,
			cpu.Memory[0x2006],
			cpu.Memory[0x2007],)
	}

	var addr uint16
	var arg uint8

	switch instr.mode {
	case A:
		addr = 0
	case imm:
		arg = cpu.Read8(cpu.PC + 1)
		addr = cpu.PC
	case zpg:
		arg = cpu.Read8(cpu.PC + 1)
		addr = uint16(arg)
	case zpgX:
		arg = cpu.Read8(cpu.PC + 1)
		addr = uint16(arg + cpu.X) & 0xFF // wrap around for X
	case zpgY:
		arg = cpu.Read8(cpu.PC + 1)
		addr = uint16(arg + cpu.Y) & 0xFF // wrap around for Y
	case abs:
		arg := cpu.Read16(cpu.PC + 1)
		addr = arg
	case absX:
		arg := cpu.Read16(cpu.PC + 1)
		addr = arg + uint16(cpu.X)
	case absY:
		arg := cpu.Read16(cpu.PC + 1)
		addr = arg + uint16(cpu.Y)
	case ind:
		arg := cpu.Read16(cpu.PC + 1)
		addr = cpu.Read16(arg)
	case indX:
		arg = cpu.Read8(cpu.PC + 1)
		addrLocation := uint16(arg + cpu.X) & 0xFF
		if addrLocation == 0xFF {
			lowByte := cpu.Read8(0xFF)
			highByte := cpu.Read8(0x00)
			addr = uint16(lowByte) | (uint16(highByte) << 8)
		} else {
			addr = cpu.Read16(addrLocation)
		}
	case indY:
		arg = cpu.Read8(cpu.PC + 1)
		addrLocation := uint16(arg)
		if addrLocation == 0xFF {
			lowByte := cpu.Read8(0xFF)
			highByte := cpu.Read8(0x00)
			addr = uint16(lowByte) | (uint16(highByte) << 8) + uint16(cpu.Y)
		} else {
			addr = cpu.Read16(addrLocation) + uint16(cpu.Y)
		}
	case rel:
		arg = cpu.Read8(cpu.PC + 1)
		addr = 0
	case impl:
		addr = 0
	default:
		log.Fatal(errors.New("Fatal: " + string(instr.mode) + " is not a valid addressing mode."))
	}

	// increment the pc based on instruction size
	cpu.PC += uint16(instr.bytes)

	switch instr.assemblyCode {
	case "ADC":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.ADC(instr, addr, value)
	case "AND":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.AND(instr, addr, value)
	case "ASL":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.ASL(instr, addr, value)
	case "BCC":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.BCC(instr, addr, value)
	case "BCS":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.BCS(instr, addr, value)
	case "BEQ":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.BEQ(instr, addr, value)
	case "BIT":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.BIT(instr, addr, value)
	case "BMI":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.BMI(instr, addr, value)
	case "BNE":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.BNE(instr, addr, value)
	case "BPL":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.BPL(instr, addr, value)
	case "BRK":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.BRK(instr, addr, value)
	case "BVC":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.BVC(instr, addr, value)
	case "BVS":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.BVS(instr, addr, value)
	case "CLC":
		cpu.CLC(instr, addr)
	case "CLD":
		cpu.CLD(instr, addr)
	case "CLI":
		cpu.CLI(instr, addr)
	case "CLV":
		cpu.CLV(instr, addr)
	case "CMP":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.CMP(instr, addr, value)
	case "CPX":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.CPX(instr, addr, value)
	case "CPY":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.CPY(instr, addr, value)
	case "DCP":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.DCP(instr, addr, value)
	case "DEC":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.DEC(instr, addr, value)
	case "DEX":
		cpu.DEX(instr, addr)
	case "DEY":
		cpu.DEY(instr, addr)
	case "EOR":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.EOR(instr, addr, value)
	case "INC":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.INC(instr, addr, value)
	case "INX":
		cpu.INX(instr, addr)
	case "INY":
		cpu.INY(instr, addr)
	case "ISC":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.ISC(instr, addr, value)
	case "JMP":
		cpu.JMP(instr, addr)
	case "JSR":
		cpu.JSR(instr, addr)
	case "LAX":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.LAX(instr, addr, value)
	case "LDA":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.LDA(instr, addr, value)
	case "LDX":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.LDX(instr, addr, value)
	case "LDY":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.LDY(instr, addr, value)
	case "LSR":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.LSR(instr, addr, value)
	case "NOP":
		cpu.NOP()
	case "ORA":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.ORA(instr, addr, value)
	case "PHA":
		cpu.PHA()
	case "PHP":
		cpu.PHP()
	case "PLA":
		cpu.PLA()
	case "PLP":
		cpu.PLP()
	case "RLA":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.RLA(instr, addr, value)
	case "ROL":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.ROL(instr, addr, value)
	case "ROR":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.ROR(instr, addr, value)
	case "RRA":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.RRA(instr, addr, value)
	case "RTI":
		cpu.RTI()
	case "RTS":
		cpu.RTS()
	case "SAX":
		cpu.SAX(instr, addr)
	case "SBC":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.SBC(instr, addr, value)
	case "SEC":
		cpu.SEC()
	case "SED":
		cpu.SED()
	case "SEI":
		cpu.SEI()
	case "SLO":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.SLO(instr, addr, value)
	case "SRE":
		value := cpu.getValue(instr.mode, addr, arg)
		cpu.SRE(instr, addr, value)
	case "STA":
		cpu.STA(instr, addr)
	case "STX":
		cpu.STX(instr, addr)
	case "STY":
		cpu.STY(instr, addr)
	case "TAX":
		cpu.TAX()
	case "TAY":
		cpu.TAY()
	case "TSX":
		cpu.TSX()
	case "TXA":
		cpu.TXA()
	case "TXS":
		cpu.TXS()
	case "TYA":
		cpu.TYA()
	default:
		//log.Fatal(errors.New("Fatal: " + string(instr.assemblyCode) + " is not a valid instruction code."))
	}

	cpu.totalCycles += uint64(instr.Cycles)
}

func (cpu *Cpu) getValue(addressingMode uint8, addr uint16, arg uint8) uint8 {
	var value uint8

	switch addressingMode {
	case A:
		value = cpu.A
	case imm:
		value = arg
	case zpg:
		value = cpu.Read8(addr)
	case zpgX:
		value = cpu.Read8(addr)
	case zpgY:
		value = cpu.Read8(addr)
	case abs:
		value = cpu.Read8(addr)
	case absX:
		value = cpu.Read8(addr)
	case absY:
		value = cpu.Read8(addr)
	case ind:
		value = cpu.Read8(addr)
	case indX:
		value = cpu.Read8(addr)
	case indY:
		value = cpu.Read8(addr)
	case rel:
		value = arg
	case impl:
		value = 0
	default:
		log.Fatal(errors.New("Fatal: " + string(addressingMode) + " is not a valid addressing mode."))
	}

	return value
}

// ADC - Add with Carry
// Performs addition with the accumulator and carry bit.
// Sets flags accordingly ZCN
func (cpu *Cpu) ADC(instr instruction, addr uint16, value uint8) {
	// Calculate the result
	result := cpu.A + value + getBit(cpu.P, 0)
	
	// Set the carry flag if unsigned overflow occurs
	if uint16(cpu.A) + uint16(value) + uint16(getBit(cpu.P, 0)) > 0xFF {
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
		cpu.Write8(addr, result)
	}

	newBit7 = getBit(result, 7)
	
	// Set the carry flag old bit 7 is 1
	if oldBit7 > 0 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	// Set the zero flag if the result is zero
	if result == 0 {
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
	if getBit(value, 6) == 1 {
		cpu.setOverflow()
	} else {
		cpu.clearOverflow()
	}

	// Set the sign bit if bit 7 is 1
	if getBit(value, 7) > 0 {
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
	if getBit(cpu.P, 7) == 0 {
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
	cpu.Push8(cpu.P | 0x30)

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
func (cpu *Cpu) CLC(instr instruction, addr uint16) {
	cpu.clearCarry()
}

// CLD - clear decimal mode
func (cpu *Cpu) CLD(instr instruction, addr uint16) {
	cpu.clearDecimal()
}

// CLI - clear interrupt disable flag
func (cpu *Cpu) CLI(instr instruction, addr uint16) {
	cpu.clearInterrupt()
}

// CLV - clear the overflow flag
func (cpu *Cpu) CLV(instr instruction, addr uint16) {
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

// DCP - Decrement from memory with C
func (cpu *Cpu) DCP(instr instruction, addr uint16, value uint8) {
	result := value - 1
	oldBit0 := getBit(value, 0)
	cpu.Write8(addr, result)

	cpu.setCHelper(oldBit0)
}

// DEC - Decrement from memory
func (cpu *Cpu) DEC(instr instruction, addr uint16, value uint8) {
	result := value - 1
	cpu.Write8(addr, result)
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// DEX - Decrement from X register
func (cpu *Cpu) DEX(instr instruction, addr uint16) {
	result := cpu.X - 1
	cpu.X = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// DEY - Decrement from Y register
func (cpu *Cpu) DEY(instr instruction, addr uint16) {
	result := cpu.Y - 1
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
	cpu.Write8(addr, result)
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// INX - Increment from X register
func (cpu *Cpu) INX(instr instruction, addr uint16) {
	result := cpu.X + 1
	cpu.X = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// INY - Increment from Y register
func (cpu *Cpu) INY(instr instruction, addr uint16) {
	result := cpu.Y + 1
	cpu.Y = result
	
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// ISC - Increment from Y register
func (cpu *Cpu) ISC(instr instruction, addr uint16, value uint8) {
	oldBit7 := getBit(cpu.A, 7)

	result := value + 1
	cpu.A -= result

	if oldBit7 != getBit(cpu.A, 7) {
		cpu.setOverflow()
	} else {
		cpu.clearOverflow()
	}

	cpu.setCHelper(oldBit7)
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// JMP - jump to address
func (cpu *Cpu) JMP(instr instruction, addr uint16) {
	cpu.PC = addr
}

// JSR - jump to subroutine
func (cpu *Cpu) JSR(instr instruction, addr uint16) {
	cpu.Push16(cpu.PC - 1)

	cpu.PC = addr
}

// LAX - load acc and Y with mem location
func (cpu *Cpu) LAX(instr instruction, addr uint16, value uint8) {
	cpu.A = value
	cpu.X = value

	cpu.setZHelper(value)
	cpu.setNHelper(value)
}

// LDA - load acc with mem location
func (cpu *Cpu) LDA(instr instruction, addr uint16, value uint8) {
	cpu.A = value

	cpu.setZHelper(value)
	cpu.setNHelper(value)
}

// LDX - load X register with mem location
func (cpu *Cpu) LDX(instr instruction, addr uint16, value uint8) {
	cpu.X = value

	cpu.setZHelper(value)
	cpu.setNHelper(value)
}

// LDY - load Y register with mem location
func (cpu *Cpu) LDY(instr instruction, addr uint16, value uint8) {
	cpu.Y = value

	cpu.setZHelper(value)
	cpu.setNHelper(value)
}

// LSR - Logical Shift Right
func (cpu *Cpu) LSR(instr instruction, addr uint16, value uint8) {
	var result uint8
	oldBit0 := getBit(value, 0)

	// If acc mode, shift the acc
	if instr.mode == A {
		result = cpu.A >> 1
		cpu.A = result
	} else {
		result = value >> 1
		cpu.Write8(addr, result)
	}
	
	// Set the carry flag old bit 7 is 1
	cpu.setCHelper(oldBit0)
	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// NOP - No operation!
func (cpu *Cpu) NOP() {
	
}

// ORA - Logical Inclusive Or
func (cpu *Cpu) ORA(instr instruction, addr uint16, value uint8) {
	result := cpu.A | value
	cpu.A = result

	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// PHA - Push Accumulator
func (cpu *Cpu) PHA() {
	cpu.Push8(cpu.A)
}

// PHP - Push Processor Status
func (cpu *Cpu) PHP() {
	// PHP always sets bit 4
	cpu.Push8(cpu.P | 0x30)
}

// PLA - Pull Accumulator
func (cpu *Cpu) PLA() {
	poppedValue := cpu.Pop8()
	cpu.A = poppedValue

	cpu.setZHelper(poppedValue)
	cpu.setNHelper(poppedValue)
}

// PLP - Pull Processor Status
func (cpu *Cpu) PLP() {
	// bit 4 and 5 are don't cares, but bit 5 is always set
	cpu.P = (cpu.Pop8() & 0xEF) | 0x20
}

func (cpu *Cpu) setCHelper(x uint8) {
	if x > 0 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}
}

// RLA - rotate memory left, then and with the acc
func (cpu *Cpu) RLA(instr instruction, addr uint16, value uint8) {
	oldBit7 := getBit(value, 7)
	var result uint8

	result = (value << 1) | getBit(cpu.P, 0)
	cpu.setCHelper(oldBit7)
	cpu.setNHelper(result)
	cpu.setZHelper(result)

	cpu.Write8(addr, result)

	cpu.A &= result
}

// ROL - rotate left
func (cpu *Cpu) ROL(instr instruction, addr uint16, value uint8) {
	oldBit7 := getBit(value, 7)
	var result uint8
	
	result = (value << 1) | getBit(cpu.P, 0)
	cpu.setCHelper(oldBit7)
	cpu.setNHelper(result)
	cpu.setZHelper(result)

	if instr.mode == A {
		cpu.A = result
	} else {
		cpu.Write8(addr, result)
	}
}

// ROR - rotate right
func (cpu *Cpu) ROR(instr instruction, addr uint16, value uint8) {
	oldBit0 := getBit(value, 0)
	var result uint8
	
	result = (value >> 1) | (getBit(cpu.P, 0) << 7)
	cpu.setCHelper(oldBit0)
	cpu.setNHelper(result)
	cpu.setZHelper(result)

	if instr.mode == A {
		cpu.A = result
	} else {
		cpu.Write8(addr, result)
	}
}

// RRA - rotate memory right, then add to the acc
func (cpu *Cpu) RRA(instr instruction, addr uint16, value uint8) {
	oldBit7 := getBit(value, 7)
	var result uint8

	result = (value >> 1) | getBit(cpu.P, 0)
	cpu.setCHelper(oldBit7)
	cpu.setNHelper(result)
	cpu.setZHelper(result)

	cpu.Write8(addr, result)

	cpu.A += result

	// Set the overflow flag
	if ((cpu.A ^ result) & (value ^ result) & 0x80) > 0 {
		cpu.setOverflow()
	} else {
		cpu.clearOverflow()
	}
}

// RTI - return from interrupt
func (cpu *Cpu) RTI() {
	processorStatus := cpu.Pop8()
	returnAddress := cpu.Pop16()
	
	cpu.P = processorStatus | 0x20
	cpu.PC = returnAddress
}

// RTS - return from subroutine
func (cpu *Cpu) RTS() {
	returnAddress := cpu.Pop16() + 1
	
	cpu.PC = returnAddress
}

// SAX - and x register with acc and store in memory
func (cpu *Cpu) SAX(instr instruction, addr uint16) {
	result := cpu.A & cpu.X
	cpu.Write8(addr, result)

	cpu.setZHelper(result)
	cpu.setNHelper(result)
}

// SBC - Subtract with Carry
func (cpu *Cpu) SBC(instr instruction, addr uint16, value uint8) {
	// TODO: understand why this can be implemented like this
	cpu.ADC(instr, addr, ^value)
}

// SEC - Set Carry Flag
func (cpu *Cpu) SEC() {
	cpu.setCarry()
}

// SED - Set Decimal Flag
func (cpu *Cpu) SED() {
	cpu.setDecimal()
}

// SEI - Set Interrupt Disable
// Sets the interrupt disable flag to one
func (cpu *Cpu) SEI() {
	cpu.setInterrupt()
}

// SLO - shifts memory left, then ors acc with memory
func (cpu *Cpu) SLO(instr instruction, addr uint16, value uint8) {
	oldBit7 := getBit(value, 7)
	result := value << 1
	cpu.Write8(addr, result)
	cpu.A |= result

	cpu.setCHelper(oldBit7)
	cpu.setZHelper(cpu.A)
	cpu.setNHelper(cpu.A)
}

// SRE - shifts memory right, then XORS acc with memory
func (cpu *Cpu) SRE(instr instruction, addr uint16, value uint8) {
	oldBit7 := getBit(value, 7)
	result := value >> 1
	cpu.Write8(addr, result)
	cpu.A ^= result

	cpu.setCHelper(oldBit7)
	cpu.setZHelper(cpu.A)
	cpu.setNHelper(cpu.A)
}

// STA - store acc in memory
func (cpu *Cpu) STA(instr instruction, addr uint16) {
	cpu.Write8(addr, cpu.A)
}

// STX - store X in memory
func (cpu *Cpu) STX(instr instruction, addr uint16) {
	cpu.Write8(addr, cpu.X)
}

// STY - store acc in memory
func (cpu *Cpu) STY(instr instruction, addr uint16) {
	cpu.Write8(addr, cpu.Y)
}

// TAX - transfer acc to X
func (cpu *Cpu) TAX() {
	cpu.X = cpu.A
	cpu.setZHelper(cpu.X)
	cpu.setNHelper(cpu.X)
}

// TAY - transfer acc to Y
func (cpu *Cpu) TAY() {
	cpu.Y = cpu.A
	cpu.setZHelper(cpu.Y)
	cpu.setNHelper(cpu.Y)
}

// TSX - transfer sp to X
func (cpu *Cpu) TSX() {
	cpu.X = cpu.SP

	cpu.setZHelper(cpu.X)
	cpu.setNHelper(cpu.X)
}

// TXA - transfer x to acc
func (cpu *Cpu) TXA() {
	cpu.A = cpu.X

	cpu.setZHelper(cpu.A)
	cpu.setNHelper(cpu.A)
}

// TXS - transfer x to sp
func (cpu *Cpu) TXS() {
	cpu.SP = cpu.X
}

// TYA - transfer y to acc
func (cpu *Cpu) TYA() {
	cpu.A = cpu.Y

	cpu.setZHelper(cpu.A)
	cpu.setNHelper(cpu.A)
}
