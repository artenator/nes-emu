package hardware

import "encoding/binary"

func (cpu *Cpu) Read8(addr uint16) uint8 {
	if addr < 0x2000 {
		return cpu.Memory[addr & 0x7FF]
	} else {
		return cpu.Memory[addr]
	}
}

func (cpu *Cpu) Read16(addr uint16) uint16 {
	
	if addr < 0x2000 {
		return binary.LittleEndian.Uint16(cpu.Memory[addr & 0x7FF : (uint32(addr) & 0x7FF) + 2])
	} else {
		return binary.LittleEndian.Uint16(cpu.Memory[addr : uint32(addr) + 2])
	}
}

func (cpu *Cpu) Write8(addr uint16, value uint8) {
	if addr < 0x2000 {
                cpu.Memory[addr & 0x7FF] = value
	} else {
                cpu.Memory[addr] = value
	}
}

func (cpu *Cpu) Write16(addr, value uint16) {
	if addr < 0x2000 {
		binary.LittleEndian.PutUint16(cpu.Memory[addr & 0x7FF : (addr & 0x7FF) + 2], value)
	} else {
		binary.LittleEndian.PutUint16(cpu.Memory[addr : addr + 2], value)
	}
}

func (cpu *Cpu) Push16(value uint16) {
	high, low := uint8(value >> 8), uint8(value & 0xFF)
	
	cpu.Memory[0x100 | uint16(cpu.SP)] = high
	cpu.Memory[(0x100 | uint16(cpu.SP)) - 1] = low
        cpu.SP -= 2
}

func (cpu *Cpu) Push8(value uint8) {
        cpu.Memory[cpu.SP] = value
}

func (cpu *Cpu) Pop16() uint16 {
	cpu.SP++
	
	barr := []uint8{cpu.Memory[cpu.SP + 1], cpu.Memory[cpu.SP]}

	cpu.SP++

	return binary.BigEndian.Uint16(barr)
}

func (cpu *Cpu) Pop8() uint8 {
	cpu.SP++

	b := cpu.Memory[cpu.SP]

	return uint8(b)
}