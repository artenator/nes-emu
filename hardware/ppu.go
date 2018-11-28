package hardware

import (
	"encoding/binary"
	"image"
	"image/color"
)

type Ppu struct {
	nes            *NES
	Memory         [0x4000]byte
	ppuAddrCounter uint8
	ppuAddrMSB     uint8
	ppuAddrLSB     uint8
	ppuAddrOffset  uint16
	oamAddr        uint8
	oamSpriteAddr  uint8
	OAM            [0x40]Sprite
	Cycle		   int64
	Scanline	   uint16
	Frame		   *image.RGBA
	FrameReady	   bool
}

func (ppu *Ppu) Write8(value uint8) {
	ppuAddressArr := []uint8{ppu.ppuAddrMSB, ppu.ppuAddrLSB}
	ppuWriteAddress := binary.BigEndian.Uint16(ppuAddressArr)

	//log.Printf("Write to ppu ram at address %x, base: %x, offset: %x, value: %x", ppuWriteAddress + ppu.ppuAddrOffset, ppuWriteAddress, ppu.ppuAddrOffset, value)

	ppu.Memory[ppuWriteAddress+ppu.ppuAddrOffset] = value

	if (ppu.nes.CPU.Memory[0x2000]>>2)&1 == 1 {
		ppu.ppuAddrOffset += 0x20
	} else {
		ppu.ppuAddrOffset++
	}
}

func (ppu *Ppu) Read8(addr uint16) uint8 {
	return ppu.Memory[addr]
}

func (ppu *Ppu) setPpuAddr(addr uint8) {
	if ppu.ppuAddrCounter%2 == 0 {
		ppu.ppuAddrMSB = addr
	} else {
		ppu.ppuAddrLSB = addr
		ppu.ppuAddrOffset = 0
	}

	ppu.ppuAddrCounter++
}

func (ppu *Ppu) SetVBlank() {
	ppu.nes.CPU.Memory[0x2002] |= 1 << 7
}

func (ppu *Ppu) ClearVBlank() {
	ppu.nes.CPU.Memory[0x2002] &= ^(uint8(1) << 7)
}

func (ppu *Ppu) get8x8Tile(base uint16, pos uint16) [8][8]uint8 {
	b0 := ppu.Memory[base+(pos*0x10) : base+(pos*0x10)+8]
	b1 := ppu.Memory[base+(pos*0x10)+8 : base+(pos*0x10)+16]

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

	b := ppu.Memory[0x2000+0x3C0+uint16(pos)]

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

func (ppu *Ppu) getSpriteColorPalette(pos uint8) [4]Color {
	switch pos {
	case 0:
		return [4]Color{
			palette[ppu.Read8(0x3F00)],
			palette[ppu.Read8(0x3F11)],
			palette[ppu.Read8(0x3F12)],
			palette[ppu.Read8(0x3F13)],
		}
	case 1:
		return [4]Color{
			palette[ppu.Read8(0x3F00)],
			palette[ppu.Read8(0x3F15)],
			palette[ppu.Read8(0x3F16)],
			palette[ppu.Read8(0x3F17)],
		}
	case 2:
		return [4]Color{
			palette[ppu.Read8(0x3F00)],
			palette[ppu.Read8(0x3F19)],
			palette[ppu.Read8(0x3F1A)],
			palette[ppu.Read8(0x3F1B)],
		}
	case 3:
		return [4]Color{
			palette[ppu.Read8(0x3F00)],
			palette[ppu.Read8(0x3F1D)],
			palette[ppu.Read8(0x3F1E)],
			palette[ppu.Read8(0x3F1F)],
		}
	}

	return [4]Color{}
}

func (ppu *Ppu) getBackgroundColorAtPixel(x, y uint8) Color {
	backgroundTileBase := uint16((ppu.nes.CPU.Memory[0x2000]>>4)&1) * 0x1000
	backgroundTileOffset := (uint16(y/8) * 32) + (uint16(x/8) % 32)
	nameTableSelect := ppu.nes.CPU.Memory[0x2000] & 0x03
	nameTableBase := 0x2000 + uint16(uint16(nameTableSelect) * 0x400)
	backgroundTilePos := ppu.Memory[nameTableBase+backgroundTileOffset]

	backgroundTile := ppu.get8x8Tile(backgroundTileBase, uint16(backgroundTilePos))
	xBG := x % 8
	yBG := y % 8

	attributePalettePos := ((y / 32) * 8) + ((x / 32) % 32)
	attributeTile := ppu.get2x2Attribute(backgroundTileBase, attributePalettePos)
	xAttr := ((x / 16) % 16) % 2
	yAttr := ((y / 16) % 16) % 2

	bgColorPalette := ppu.getBackgroundColorPalette(attributeTile[yAttr][xAttr])
	bgColor := bgColorPalette[backgroundTile[yBG][xBG]]

	return bgColor
}

func (ppu *Ppu) getBackgroundColorAtPixelOptimized(x, y uint8, backgroundTile [8][8]uint8, attributeTile [2][2]uint8) Color {
	xBG := x % 8
	yBG := y % 8

	xAttr := ((x / 16) % 16) % 2
	yAttr := ((y / 16) % 16) % 2

	bgColorPalette := ppu.getBackgroundColorPalette(attributeTile[yAttr][xAttr])
	bgColor := bgColorPalette[backgroundTile[yBG][xBG]]

	return bgColor
}

func (ppu *Ppu) getSpriteColorAtPixel(x, y uint8, s Sprite) Color {
	flipHorizontal := false
	if (s.attributes >> 6) & 1 == 1 {
		flipHorizontal = true
	}

	backgroundTileBase := uint16((ppu.nes.CPU.Memory[0x2000]>>3)&1) * 0x1000
	backgroundTilePos := s.tileNum

	backgroundTile := ppu.get8x8Tile(backgroundTileBase, uint16(backgroundTilePos))
	xBG := x % 8
	if flipHorizontal {
		xBG = 7 - x % 8
	}

	yBG := y % 8

	spriteColorPalette := ppu.getSpriteColorPalette(s.attributes & 0x03)
	spriteColor := spriteColorPalette[backgroundTile[yBG][xBG]]

	// Check if sprites hide background
	if backgroundTile[yBG][xBG] == 0 {
		spriteColor.A = 0
	}

	return spriteColor
}

func (ppu *Ppu) GetColorAtPixel(x, y uint8) Color {
	var color Color

	for _, sprite := range ppu.OAM {
		if sprite.yCoord > 0x00 && sprite.yCoord < 0xEF {
			// TODO: 8x16 sprites
			inRangeX := x >= sprite.xCoord && x < sprite.xCoord+8
			inRangeY := y >= sprite.yCoord && y < sprite.yCoord+8
			if inRangeX && inRangeY {
				spriteColor := ppu.getSpriteColorAtPixel(x - sprite.xCoord, y - sprite.yCoord, sprite)
				if spriteColor.A > 0 {
					return spriteColor
				} else {
					color = ppu.getBackgroundColorAtPixel(x, y)
					return color
				}
			}
		}
	}

	color = ppu.getBackgroundColorAtPixel(x, y)

	return color
}

func (ppu *Ppu) GetColorAtPixelOptimized(x, y uint8, backgroundTile [8][8]uint8, attributeTile [2][2]uint8) Color {
	var color Color

	for _, sprite := range ppu.OAM {
		if sprite.yCoord > 0x00 && sprite.yCoord < 0xEF {
			// TODO: 8x16 sprites
			inRangeX := x >= sprite.xCoord && x < sprite.xCoord+8
			inRangeY := y >= sprite.yCoord && y < sprite.yCoord+8
			if inRangeX && inRangeY {
				spriteColor := ppu.getSpriteColorAtPixel(x - sprite.xCoord, y - sprite.yCoord, sprite)
				if spriteColor.A > 0 {
					return spriteColor
				} else {
					color = ppu.getBackgroundColorAtPixelOptimized(x, y, backgroundTile, attributeTile)
					return color
				}
			}
		}
	}

	color = ppu.getBackgroundColorAtPixelOptimized(x, y, backgroundTile, attributeTile)

	return color
}

func (ppu *Ppu) PPURun() {

	if ppu.Cycle == 340 && ppu.Scanline < 240 {
		backgroundTileBase := uint16((ppu.nes.CPU.Memory[0x2000]>>4)&1) * 0x1000
		backgroundTileOffset := (uint16(0/8) * 32) + (uint16(0/8) % 32)
		nameTableSelect := ppu.nes.CPU.Memory[0x2000] & 0x03
		nameTableBase := 0x2000 + uint16(uint16(nameTableSelect) * 0x400)
		backgroundTilePos := ppu.Memory[nameTableBase+backgroundTileOffset]
		backgroundTile := ppu.get8x8Tile(backgroundTileBase, uint16(backgroundTilePos))
		attributePalettePos := uint8((0 / 32) * 8) + ((0 / 32) % 32)
		attributeTile := ppu.get2x2Attribute(backgroundTileBase, attributePalettePos)
		sl := ppu.Scanline
		for x := 0; x < 256; x++ {
			if x % 8 == 0 {
				backgroundTileOffset = (uint16(sl/8) * 32) + (uint16(x/8) % 32)
				nameTableBase = 0x2000 + uint16(uint16(nameTableSelect) * 0x400)
				backgroundTilePos = ppu.Memory[nameTableBase+backgroundTileOffset]
				backgroundTile = ppu.get8x8Tile(backgroundTileBase, uint16(backgroundTilePos))
				attributePalettePos = uint8((sl / 32) * 8) + ((uint8(x) / 32) % 32)
				attributeTile = ppu.get2x2Attribute(backgroundTileBase, attributePalettePos)
			}
			c := ppu.GetColorAtPixelOptimized(uint8(x), uint8(sl), backgroundTile, attributeTile)
			ppu.Frame.SetRGBA(x, int(sl), color.RGBA{c.R, c.G, c.B, uint8(c.A)})
		}
	}

	if ppu.Scanline == 259 && ppu.Cycle == 340 {
		ppu.ClearVBlank()
	}
	if ppu.Scanline == 241  && ppu.Cycle == 0 {
		ppu.SetVBlank()
	}

	ppu.Cycle = (ppu.Cycle + 1) % 341
	if ppu.Cycle == 340 {
		ppu.Scanline = (ppu.Scanline + 1) % 260
	}

	// if a frame is ready, set bool
	if ppu.Cycle == 0 && ppu.Scanline == 0 {
		ppu.FrameReady = true
	}
}

func (ppu *Ppu) RunPPUCycles(numOfCycles uint16) {
	for i := uint16(0); i < numOfCycles; i++ {
		ppu.PPURun()
	}
}

func (ppu *Ppu) GenerateFrame() *image.RGBA {
	var img *image.RGBA = image.NewRGBA(image.Rect(0, 0, 256, 240))

	for y := 0; y < 240; y++ {
		for x := 0; x < 256; x++ {
			c := ppu.GetColorAtPixel(uint8(x), uint8(y))
			img.SetRGBA(x, y, color.RGBA{c.R, c.G, c.B, uint8(c.A)})
		}
	}

	return img
}