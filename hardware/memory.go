package hardware

import (
	"encoding/binary"
	"log"
)

func (cpu *Cpu) Read8(addr uint16) uint8 {
	if addr < 0x2000 {
		return cpu.Memory[addr&0x7FF]
	} else if addr >= 0x2000 && addr < 0x4000 {
		return cpu.Memory[addr&0x2007]
	} else {
		return cpu.Memory[addr]
	}
}

func (cpu *Cpu) Read16(addr uint16) uint16 {

	// TODO: @artenator this read totally messes up at the end of the address space. Fix it

	if addr < 0x2000 {
		return binary.LittleEndian.Uint16(cpu.Memory[addr&0x7FF : (uint32(addr)&0x7FF)+2])
	} else if addr >= 0x2000 && addr < 0x4000 {
		if addr == 0x2002 {
			log.Println("status read")
			cpu.nes.PPU.ClearVBlank()
		}
		return binary.LittleEndian.Uint16(cpu.Memory[addr&0x2007 : (uint32(addr)&0x2007)+2])
	} else {
		return binary.LittleEndian.Uint16(cpu.Memory[addr : uint32(addr)+2])
	}
}

func (cpu *Cpu) Write8(addr uint16, value uint8) {
	if addr < 0x2000 {
		cpu.Memory[addr&0x7FF] = value
	} else if addr >= 0x2000 && addr < 0x4000 {
		truncAddr := addr & 0x2007

		cpu.Memory[truncAddr] = value

		if addr >= 0x2000 && addr < 0x4000 {
			// PPUADDR
			if truncAddr == 0x2006 {
				cpu.nes.PPU.setPpuAddr(cpu.A)
			}
			// PPUDATA
			if truncAddr == 0x2007 {
				cpu.nes.PPU.Write8(cpu.A)
			}
			// OAMADDR
			if truncAddr == 0x2003 {
				cpu.nes.PPU.SetOamAddr(cpu.A)
			}
			// OAMDATA
			if truncAddr == 0x2004 {
				cpu.nes.PPU.WriteOAM8(cpu.A)
			}

		}
	} else {
		cpu.Memory[addr] = value
	}
}

func (cpu *Cpu) Write16(addr, value uint16) {
	if addr < 0x2000 {
		binary.LittleEndian.PutUint16(cpu.Memory[addr&0x7FF:(addr&0x7FF)+2], value)
	} else if addr >= 0x2000 && addr < 0x4000 {
		binary.LittleEndian.PutUint16(cpu.Memory[addr&0x2007:(addr&0x2007)+2], value)
	} else {
		binary.LittleEndian.PutUint16(cpu.Memory[addr:addr+2], value)
	}
}

func (cpu *Cpu) Push16(value uint16) {
	high, low := uint8(value>>8), uint8(value&0xFF)

	cpu.Memory[0x100|uint16(cpu.SP)] = high
	cpu.Memory[(0x100|uint16(cpu.SP))-1] = low
	cpu.SP -= 2
}

func (cpu *Cpu) Push8(value uint8) {
	cpu.Memory[0x100|uint16(cpu.SP)] = value
	cpu.SP--
}

func (cpu *Cpu) Pop16() uint16 {
	cpu.SP++

	barr := []uint8{cpu.Memory[(0x100|uint16(cpu.SP))+1], cpu.Memory[0x100|uint16(cpu.SP)]}

	cpu.SP++

	return binary.BigEndian.Uint16(barr)
}

func (cpu *Cpu) Pop8() uint8 {
	cpu.SP++

	b := cpu.Memory[0x100|uint16(cpu.SP)]

	return uint8(b)
}
