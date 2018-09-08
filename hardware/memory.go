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
		return binary.LittleEndian.Uint16(cpu.Memory[addr ^ 0x7FF : (addr ^ 0x7FF) + 2])
	} else {
		return binary.LittleEndian.Uint16(cpu.Memory[addr : addr + 2])
	}
}
