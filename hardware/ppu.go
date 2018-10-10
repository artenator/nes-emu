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

func (ppu *Ppu) SetVBlank() {
	ppu.nes.CPU.Memory[0x2002] |= 1 << 7

	if (ppu.nes.CPU.Memory[0x2000] >> 7) & 1 == 1 {
		//log.Println("NMI Interrupt")
		ppu.nes.CPU.handleNMI()
	}
}

func (ppu *Ppu) ClearVBlank() {
	ppu.nes.CPU.Memory[0x2002] &= ^(uint8(1) << 7)
}

func (ppu *Ppu) get8x8Tile(base uint16, pos uint16) [8][8]uint8 {
	b0 := ppu.Memory[base + (pos * 0x10) : base + (pos * 0x10) + 8]
	b1 := ppu.Memory[base + (pos * 0x10) + 8 : base + (pos * 0x10) + 16]

	var result [8][8]uint8

	for i := 0; i < 8; i++ {
		barr0 := b0[i]
		barr1 := b1[i]
		for j := uint8(0); j < 8; j++ {
			var biResult uint8
			biResult |= (barr0 >> (7 - j)) & 1
			biResult |= ((barr1 >> (7 - j)) & 1) << 1
			result[i][j] = biResult
		}
	}

	return result
}

func (ppu *Ppu) get2x2Attribute(base uint16, pos uint8) [2][2]uint8 {
	var result [2][2]uint8

	b := ppu.Memory[uint16(pos) + 0x3C0]

	result[0][0] = b & 0x03
	result[0][1] = (b & 0x0C) >> 2
	result[1][0] = (b & 0x30) >> 4
	result[1][1] = (b & 0xC0) >> 6

	return result
}

func (ppu *Ppu) getBackgroundColorPalette(pos uint8) [4]Color {
	switch pos {
	case 0:
		return [4]Color{
			palette[ppu.Read8(0x3F00)],
			palette[ppu.Read8(0x3F01)],
			palette[ppu.Read8(0x3F02)],
			palette[ppu.Read8(0x3F03)],
		}
	case 1:
		return [4]Color{
			palette[ppu.Read8(0x3F00)],
			palette[ppu.Read8(0x3F05)],
			palette[ppu.Read8(0x3F06)],
			palette[ppu.Read8(0x3F07)],
		}
	case 2:
		return [4]Color{
			palette[ppu.Read8(0x3F00)],
			palette[ppu.Read8(0x3F09)],
			palette[ppu.Read8(0x3F0A)],
			palette[ppu.Read8(0x3F0B)],
		}
	case 3:
		return [4]Color{
			palette[ppu.Read8(0x3F00)],
			palette[ppu.Read8(0x3F0D)],
			palette[ppu.Read8(0x3F0E)],
			palette[ppu.Read8(0x3F0F)],
		}
	}

	return [4]Color{}
}

func (ppu *Ppu) GetColorAtPixel(x, y uint8) Color {
	backgroundTileBase := uint16((ppu.nes.CPU.Memory[0x2000] >> 4) & 1) * 0x1000
	backgroundTileOffset := (uint16(y / 8) * 32) + (uint16(x / 8) % 32)
	backgroundTilePos := ppu.Memory[0x2000 + backgroundTileOffset]


	backgroundTile := ppu.get8x8Tile(backgroundTileBase, uint16(backgroundTilePos))
	xBG := x % 8
	yBG := y % 8

	attributePalettePos := ((y / 32) * 8) + ((x / 32) % 32)
	attributeTile := ppu.get2x2Attribute(backgroundTileBase, attributePalettePos)
	xAttr := ((x / 16) % 16) % 2
	yAttr := ((y / 16) % 16) % 2

	bgColorPalette := ppu.getBackgroundColorPalette(attributeTile[xAttr][yAttr])
	bgColor := bgColorPalette[backgroundTile[yBG][xBG]]

	if ppu.Memory[0x3f00] == 0x0f && x > 250 && y > 230 {
		//log.Println("Hello ")
	}

	return bgColor
}