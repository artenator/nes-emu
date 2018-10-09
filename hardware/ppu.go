package hardware

import (
	"encoding/binary"
)

type Ppu struct {
	nes *NES
	Memory [0x4000]byte
	ppuAddrCounter uint8
	ppuAddrMSB uint8
	ppuAddrLSB uint8
	ppuAddrOffset uint16
}

func (ppu *Ppu) Write8(value uint8){
	ppuAddressArr := []uint8{ppu.ppuAddrMSB, ppu.ppuAddrLSB}
	ppuWriteAddress := binary.BigEndian.Uint16(ppuAddressArr)

	//log.Printf("Write to ppu ram at address %x, base: %x, offset: %x, value: %x", ppuWriteAddress + ppu.ppuAddrOffset, ppuWriteAddress, ppu.ppuAddrOffset, value)

	ppu.Memory[ppuWriteAddress + ppu.ppuAddrOffset] = value

	if (ppu.nes.CPU.Memory[0x2000] >> 2) & 1 == 1 {
		ppu.ppuAddrOffset += 0x20
	} else {
		ppu.ppuAddrOffset++
	}
}

func (ppu *Ppu) Read8(addr uint16) uint8 {
	return ppu.Memory[addr]
}

func (ppu *Ppu) setPpuAddr(addr uint8) {
	if ppu.ppuAddrCounter % 2 == 0 {
		ppu.ppuAddrMSB = addr
	} else {
		ppu.ppuAddrLSB = addr
		ppu.ppuAddrOffset = 0
	}

	ppu.ppuAddrCounter++
}

func (ppu *Ppu) setVBlank() {
	ppu.nes.CPU.Memory[0x2002] |= 1 << 7

	if (ppu.nes.CPU.Memory[0x2000] >> 7) & 1 == 1 {
		//log.Println("NMI Interrupt")
		ppu.nes.CPU.handleNMI()
	}
}

func (ppu *Ppu) clearVBlank() {
	ppu.nes.CPU.Memory[0x2002] &= ^(uint8(1) << 7)
}

func (ppu *Ppu) get8x8Tile(base uint16, pos uint8) [8][8]uint8 {
	b0 := ppu.Memory[base : base + 8]
	b1 := ppu.Memory[base + 8 : base + 16]

	var result [8][8]uint8

	for i := 0; i < 8; i++ {
		barr0 := b0[i]
		barr1 := b1[i]
		for j := uint8(0); j < 8; j++ {
			var biResult uint8
			biResult |= (barr0 >> j) & 1
			biResult |= ((barr1 >> j) & 1) << 1
			result[i][j] = biResult
		}
	}

	return result
}