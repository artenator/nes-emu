package hardware

import (
	"encoding/binary"
)

func (cpu *Cpu) Read8(addr uint16) uint8 {
	if addr < 0x2000 {
		return cpu.Memory[addr&0x7FF]
	} else if addr >= 0x2000 && addr < 0x4000 {
		truncAddr := addr&0x2007
		readValue := cpu.Memory[truncAddr]
		if truncAddr == 0x2002 {
			cpu.nes.PPU.ClearVBlank()
			cpu.nes.PPU.ppuAddrCounter = 0
		}
		if truncAddr == 0x2007 {
			readValue = cpu.nes.PPU.DataRead()
		}
		return readValue
	} else if addr == 0x4016 {
		val := (cpu.Controller >> (7 - (cpu.ControllerIdx % 8))) & 1
		cpu.ControllerIdx++
		return val
	} else if addr >= 0x6000 {
		readFromMapper := cpu.nes.CARTIO.read8(addr)
		return readFromMapper
	} else {
		return cpu.Memory[addr]
	}
}

func (cpu *Cpu) Read16(addr uint16) uint16 {

	// TODO: @artenator this read totally messes up at the end of the address space. Fix it

	if addr < 0x2000 {
		return binary.LittleEndian.Uint16(cpu.Memory[addr : uint32(addr)+2])
	} else if addr >= 0x2000 && addr < 0x4000 {
		return binary.LittleEndian.Uint16(cpu.Memory[addr&0x2007 : (uint32(addr)&0x2007)+2])
	} else if addr >= 0x6000 {
		return cpu.nes.CARTIO.read16(addr)
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

		// PPUCTRL
		if truncAddr == 0x2000 {
			cpu.nes.PPU.ppuctrl.setValues(value)
		}
		// PPUMASK
		if truncAddr == 0x2001 {
			cpu.nes.PPU.ppumask.setValues(value)
		}
		// PPUADDR
		if truncAddr == 0x2006 {
			cpu.nes.PPU.setPpuAddr(value)
		}
		// PPUDATA
		if truncAddr == 0x2007 {
			cpu.nes.PPU.Write8(value)
		}
		// OAMADDR
		if truncAddr == 0x2003 {
			cpu.nes.PPU.SetOamAddr(value)
		}
		// OAMDATA
		if truncAddr == 0x2004 {
			cpu.nes.PPU.WriteOAM8(value)
		}
		// PPUSCROLL
		if truncAddr == 0x2005 {
			cpu.nes.PPU.setPpuScrollAddr(value)
			//log.Printf("writing to ppuscroll 0x%x", value)
			//log.Println(strconv.FormatInt(int64(cpu.Memory[0x2000]), 2))
			//log.Printf("+%v", cpu.nes.PPU.Memory[0x2000:0x3000])
		}
	} else {
		cpu.Memory[addr] = value

		if addr == 0x4001 {
			cpu.nes.APU.sweep1.setSweepValues(value)
		} else if addr == 0x4003 {
			cpu.nes.APU.pulse1.setTargetTimer()
		} else if addr == 0x4005 {
			cpu.nes.APU.sweep2.setSweepValues(value)
		} else if addr == 0x4007 {
			cpu.nes.APU.pulse2.setTargetTimer()
		} else if addr == 0x4008 {
			cpu.nes.APU.triangle.setLinearCounterValues(value)
		} else if addr == 0x400B {
			cpu.nes.APU.triangle.linearReload = true
			if cpu.nes.APU.enableTriangle {
				cpu.nes.APU.triangle.setLengthCounter(value)
			}
			// OAMDMA at 0x4014 write
		} else if addr == 0x4014 {
			// write all the sprites to oam
			startPos := uint16(value) << 8
			//log.Printf("start pos is %x", startPos)
			for idx, _ := range cpu.Memory[startPos:startPos + 0x100] {
				byteRead := cpu.Read8(startPos+uint16(idx))

				cpu.nes.PPU.WriteOAM8(byteRead)
			}
			//log.Printf("OAM %+v", cpu.nes.PPU.OAM)
			//log.Printf("OAM %s", cpu.nes.PPU.OAM)
			cpu.nes.PPU.SetOamAddr(0)
			cpu.nes.PPU.oamSpriteAddr = 0

		} else if addr == 0x4015 {
			cpu.nes.APU.enablePulseChannel1 = (value >> 0) & 1 == 1
			if !cpu.nes.APU.enablePulseChannel1 {
				cpu.nes.APU.pulse1.lengthTimer = 0
			}
			cpu.nes.APU.enablePulseChannel2 = (value >> 1) & 1 == 1
			if !cpu.nes.APU.enablePulseChannel2 {
				cpu.nes.APU.pulse2.lengthTimer = 0
			}
			cpu.nes.APU.enableTriangle = (value >> 2) & 1 == 1
			if !cpu.nes.APU.enableTriangle {
				cpu.nes.APU.triangle.lengthCounter = 0
			}
		} else if addr == 0x4016 {
			if value == 0 {
				cpu.ControllerIdx = 0
			}
		} else if addr == 0x4017 {
			cpu.nes.APU.setFrameCounterValues(value)
		} else if addr >= 0x6000 {
			cpu.nes.CARTIO.write8(addr, value)
		}
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
