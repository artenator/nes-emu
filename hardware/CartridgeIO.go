package hardware

import (
	"encoding/binary"
	"log"
)

type CartridgeIO interface {
	read8(addr uint16) uint8
	read16(addr uint16) uint16
	write8(addr uint16, value uint8)
	initCartIO(cartridge *Cartridge)
}

type Mapper0CIO struct {
	cartridge *Cartridge
}

func (m *Mapper0CIO) setMirrorStyle() {
	mirrorFlag := m.cartridge.flags6 & 1

	switch mirrorFlag {
	case 0:
		m.cartridge.mirrorStyle = horizontal
	case 1:
		m.cartridge.mirrorStyle = vertical
	}
}

func (m *Mapper0CIO) initCartIO(cartridge *Cartridge) {
	m.cartridge = cartridge
	m.setMirrorStyle()
}

func (m *Mapper0CIO) read8(addr uint16) uint8 {
	if addr < 0x2000 {
		return m.cartridge.chrRom[addr]
	} else if addr >= 0x8000 {
		if m.cartridge.prgRomBlocks == 1 {
			return m.cartridge.prgRom[(addr - 0x8000) % 0x4000]
		} else if m.cartridge.prgRomBlocks == 2 {
			return m.cartridge.prgRom[(addr - 0x8000)]
		}
	}

	return 0
}

func (m *Mapper0CIO) read16(addr uint16) uint16 {
	if addr >= 0x8000 {
		if m.cartridge.prgRomBlocks == 1 {
			return binary.LittleEndian.Uint16(m.cartridge.prgRom[(addr - 0x8000) % 0x4000 : (addr - 0x8000) % 0x4000 + 2])
		} else if m.cartridge.prgRomBlocks == 2 {
			return binary.LittleEndian.Uint16(m.cartridge.prgRom[(addr - 0x8000) : (addr - 0x8000) + 2])
		}

	}

	return binary.LittleEndian.Uint16(m.cartridge.prgRom[(addr - 0x8000) : (addr - 0x8000) + 2])
}

func (m *Mapper0CIO) write8(addr uint16, value uint8) {

}

type Mapper1CIO struct {
	cartridge *Cartridge
	shiftReg byte

	//Control (internal, $8000-$9FFF)
	//4bit0
	//-----
	//CPPMM
	//|||||
	//|||++- Mirroring (0: one-screen, lower bank; 1: one-screen, upper bank;
	//|||               2: vertical; 3: horizontal)
	//|++--- PRG ROM bank mode (0, 1: switch 32 KB at $8000, ignoring low bit of bank number;
	//|                         2: fix first bank at $8000 and switch 16 KB bank at $C000;
	//|                         3: fix last bank at $C000 and switch 16 KB bank at $8000)
	//+----- CHR ROM bank mode (0: switch 8 KB at a time; 1: switch two separate 4 KB banks)
	controlBank byte

	//CHR bank 0 (internal, $A000-$BFFF)
	//4bit0
	//-----
	//CCCCC
	//|||||
	//+++++- Select 4 KB or 8 KB CHR bank at PPU $0000 (low bit ignored in 8 KB mode)
	chrBank1 byte

	//CHR bank 1 (internal, $C000-$DFFF)
	//4bit0
	//-----
	//CCCCC
	//|||||
	//+++++- Select 4 KB CHR bank at PPU $1000 (ignored in 8 KB mode)
	chrBank2 byte

	//PRG bank (internal, $E000-$FFFF)
	//	//4bit0
	//	//-----
	//	//RPPPP
	//	//|||||
	//	//|++++- Select 16 KB PRG ROM bank (low bit ignored in 32 KB mode)
	//	//+----- PRG RAM chip enable (0: enabled; 1: disabled; ignored on MMC1A)
	prgBank byte

	chrRomBankMode byte
	prgRomBankMode byte
}

const (
	//Chr modes
	//(0: switch 8 KB at a time; 1: switch two separate 4 KB banks)
	chrBankMode0 = iota
	chrBankMode1
)

const (
	//Prg Modes
	//(0, 1: switch 32 KB at $8000, ignoring low bit of bank number;
	//2: fix first bank at $8000 and switch 16 KB bank at $C000;
	//3: fix last bank at $C000 and switch 16 KB bank at $8000)
	prgBankMode0 = iota
	prgBankMode1
	prgBankMode2
	prgBankMode3
)

func (m *Mapper1CIO) initCartIO(cartridge *Cartridge) {
	m.cartridge = cartridge
	m.shiftReg = 0x10
}

func (m *Mapper1CIO) read8(addr uint16) uint8 {
	if addr >= 0x8000 {
		switch m.prgRomBankMode {
		case prgBankMode0, prgBankMode1:
			truncOffsetAddr := addr & 0x7FFF
			baseAddrIdx := m.prgBank & 0xE // ignore last bit in 32kb mode
			return m.cartridge.prgRom[uint16(baseAddrIdx % m.cartridge.prgRomBlocks) * 0x4000 + truncOffsetAddr]
		case prgBankMode2:
			truncOffsetAddr := addr & 0x7FFF
			if truncOffsetAddr < 0x4000 {
				return m.cartridge.prgRom[truncOffsetAddr]
			} else {
				baseAddrIdx := m.prgBank & 0xF
				return m.cartridge.prgRom[uint16(baseAddrIdx % m.cartridge.prgRomBlocks) * 0x4000 + truncOffsetAddr]
			}
		case prgBankMode3:
			truncOffsetAddr := addr & 0x7FFF
			if truncOffsetAddr < 0x4000 {
				baseAddrIdx := m.prgBank & 0xF
				return m.cartridge.prgRom[uint16(baseAddrIdx % m.cartridge.prgRomBlocks) * 0x4000 + truncOffsetAddr]
			} else {
				return m.cartridge.prgRom[uint16(m.cartridge.prgRomBlocks - 1) * 0x4000 + truncOffsetAddr]
			}
		default:
			log.Fatalf("Invalid prg bank mode %d", m.prgRomBankMode)
		}
	}
	return 0
}

func (m *Mapper1CIO) read16(addr uint16) uint16 {
	if addr >= 0x8000 {
		switch m.prgRomBankMode {
		case prgBankMode0:
		case prgBankMode1:
			truncOffsetAddr := addr & 0x7FFF
			baseAddrIdx := m.prgBank & 0xE // ignore last bit in 32kb mode
			baseAddr := uint16(baseAddrIdx % m.cartridge.prgRomBlocks) * 0x4000
			return binary.LittleEndian.Uint16(m.cartridge.prgRom[baseAddr + truncOffsetAddr : baseAddr + truncOffsetAddr + 2])
		case prgBankMode2:
			truncOffsetAddr := addr & 0x7FFF
			if truncOffsetAddr < 0x4000 {
				return binary.LittleEndian.Uint16(m.cartridge.prgRom[truncOffsetAddr : truncOffsetAddr + 2])
			} else {
				baseAddrIdx := m.prgBank & 0xF
				baseAddr := uint16(baseAddrIdx % m.cartridge.prgRomBlocks) * 0x4000
				return binary.LittleEndian.Uint16(m.cartridge.prgRom[baseAddr + truncOffsetAddr : baseAddr + truncOffsetAddr + 2])
			}
		case prgBankMode3:
			truncOffsetAddr := addr & 0x7FFF
			if truncOffsetAddr < 0x4000 {
				baseAddrIdx := m.prgBank & 0xF
				baseAddr := uint16(baseAddrIdx % m.cartridge.prgRomBlocks) * 0x4000
				return binary.LittleEndian.Uint16(m.cartridge.prgRom[baseAddr + truncOffsetAddr : baseAddr + truncOffsetAddr + 2])
			} else {
				baseAddr := uint16(m.cartridge.prgRomBlocks - 1) * 0x4000
				return binary.LittleEndian.Uint16(m.cartridge.prgRom[baseAddr + truncOffsetAddr : baseAddr + truncOffsetAddr + 2])
			}
		default:
			log.Fatalf("Invalid prg bank mode %d", m.prgRomBankMode)
		}
	}

	return binary.LittleEndian.Uint16(m.cartridge.prgRom[(addr - 0x8000) : (addr - 0x8000) + 2])
}

func (m *Mapper1CIO) setMirrorStyle() {
	mirrorFlag := m.controlBank & 0x3

	switch mirrorFlag {
	case 0:
		m.cartridge.mirrorStyle = singleScreen
	case 1:
		m.cartridge.mirrorStyle = singleScreen
	case 2:
		m.cartridge.mirrorStyle = vertical
	case 3:
		m.cartridge.mirrorStyle = horizontal
	}
}

func (m *Mapper1CIO) setBankModes() {
	m.chrRomBankMode = (m.controlBank >> 4) & 1
	m.prgRomBankMode = (m.controlBank >> 2) & 0x3
}

func (m *Mapper1CIO) setRegister(addr uint16, regValue byte) {
	if addr >= 0x8000 && addr < 0xA000 {
		m.controlBank = regValue
		m.setMirrorStyle()
		m.setBankModes()
		log.Printf("Setting Control Register for MMC1 %x", regValue)
		log.Printf("Bank Modes: chr %x prg %x", m.chrRomBankMode, m.prgRomBankMode)
	} else if addr >= 0xA000 && addr < 0xC000 {
		m.chrBank1 = regValue
	} else if addr >= 0xC000 && addr < 0xE000 {
		m.chrBank2 = regValue
	} else if addr >= 0xE000 && addr <= 0xFFFF {
		m.prgBank = regValue
	}
}

func (m *Mapper1CIO) write8(addr uint16, value uint8) {
	if addr >= 0x8000 {
		isReset := getBit(value, 7) == 1
		isFifthWrite := m.shiftReg & 1 == 1

		m.shiftReg = (m.shiftReg >> 1) | ((value & 1) << 4)

		if isFifthWrite && !isReset {
			m.setRegister(addr, m.shiftReg)
		}

		if isReset || isFifthWrite {
			m.shiftReg = 0x10
		}
	}
}