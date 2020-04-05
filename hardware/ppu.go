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
	ppuScrollMSB     uint8
	ppuScrollLSB     uint8
	ppuAddrOffset  uint16
	oamAddr        uint8
	oamSpriteAddr  uint8
	OAM            [0x40]Sprite
	Cycle		   int64
	Scanline	   uint16
	Frame		   *image.RGBA
	FrameReady	   bool

	ppumask     PpuMask
	ppuctrl PpuCtrl
	NmiOccurred bool
	PpuReady bool

	currentTiles [0x21][8][8]uint8
	currentAttributes [0x21][2][2]uint8

	currentSprites [0x40]Sprite
	spriteCount int

	scalingFactor int
}

type PpuCtrl struct {
	nmiGenerate uint8
	ppuMasterSlaveSelect uint8
	spriteSize uint8
	backgroundPatternTableAddr uint8
	spritePatternTableAddr uint8
	vramAddressIncrement uint8
	baseNametableAddr uint8
}

type PpuMask struct {
	emphBlue bool
	emphGreen bool
	emphRed bool
	spriteEnable bool
	backgroundEnable bool
	spriteLeftColumnEnable bool
	backgroundLeftColumnEnable bool
	greyScale bool
}

func (mask *PpuCtrl) setValues(value uint8) {
	mask.nmiGenerate = (value >> 7) & 1
	mask.ppuMasterSlaveSelect = (value >> 6) & 1
	mask.spriteSize = (value >> 5) & 1
	mask.backgroundPatternTableAddr = (value >> 4) & 1
	mask.spritePatternTableAddr = (value >> 3) & 1
	mask.vramAddressIncrement = (value >> 2) & 1
	mask.baseNametableAddr = (value) & 0x3
}

func (mask *PpuMask) setValues(value uint8) {
	mask.emphBlue = (value >> 7) & 1 == 1
	mask.emphGreen = (value >> 6) & 1 == 1
	mask.emphRed = (value >> 5) & 1 == 1
	mask.spriteEnable = (value >> 4) & 1 == 1
	mask.backgroundEnable = (value >> 3) & 1 == 1
	mask.spriteLeftColumnEnable = (value >> 2) & 1 == 1
	mask.backgroundLeftColumnEnable = (value >> 1) & 1 == 1
	mask.greyScale = value & 1 == 1
}

func (ppu *Ppu) DataRead() uint8 {
	ppuAddressArr := []uint8{ppu.ppuAddrMSB, ppu.ppuAddrLSB}
	ppuWriteAddress := binary.BigEndian.Uint16(ppuAddressArr)

	absReadAddress := ((ppuWriteAddress+ppu.ppuAddrOffset) - 1) % 0x3FFF

	//log.Printf("reading ppu 0x%x, value: 0x%x, OFFSET: %d", absReadAddress, ppu.Read8(absReadAddress), ppu.ppuAddrOffset)

	ppu.incrementAddress()

	return ppu.Read8(absReadAddress)
}

func (ppu *Ppu) Write8(value uint8) {
	ppuAddressArr := []uint8{ppu.ppuAddrMSB, ppu.ppuAddrLSB}
	ppuWriteAddress := binary.BigEndian.Uint16(ppuAddressArr)

	absWriteAddress := (ppuWriteAddress+ppu.ppuAddrOffset) & 0x3FFF

	ppu.Memory[absWriteAddress] = value

	switch ppu.nes.CART.mirrorStyle {
	case horizontal:
		if absWriteAddress >= 0x2000 && absWriteAddress < 0x2400 {
			ppu.Memory[absWriteAddress + 0x400] = value
		} else if absWriteAddress >= 0x2400 && absWriteAddress < 0x2800 {
			ppu.Memory[absWriteAddress - 0x400] = value
		} else if absWriteAddress >= 0x2800 && absWriteAddress < 0x2C00 {
			ppu.Memory[absWriteAddress + 0x400] = value
		} else if absWriteAddress >= 0x2C00 && absWriteAddress < 0x3000 {
			ppu.Memory[absWriteAddress - 0x400] = value
		}
	case vertical:
		if absWriteAddress >= 0x2000 && absWriteAddress < 0x2400 {
			ppu.Memory[absWriteAddress + 0x800] = value
		} else if absWriteAddress >= 0x2400 && absWriteAddress < 0x2800 {
			ppu.Memory[absWriteAddress + 0x800] = value
		} else if absWriteAddress >= 0x2800 && absWriteAddress < 0x2C00 {
			ppu.Memory[absWriteAddress - 0x800] = value
		} else if absWriteAddress >= 0x2C00 && absWriteAddress < 0x3000 {
			ppu.Memory[absWriteAddress - 0x800] = value
		}
	}


	//log.Printf("writing ppu 0x%x, value: 0x%x, OFFSET: %d", absWriteAddress, ppu.Read8(absWriteAddress), ppu.ppuAddrOffset)

	ppu.incrementAddress()
}

func (ppu *Ppu) incrementAddress() {
	incrementDown := (ppu.nes.CPU.Memory[0x2000]>>2)&1 == 1
	if incrementDown {
		ppu.ppuAddrOffset += 0x20
	} else {
		ppu.ppuAddrOffset++
	}
}

func (ppu *Ppu) Read8(addr uint16) uint8 {
	if addr < 0x2000 {
		return ppu.nes.CARTIO.read8(addr)
	} else {
		return ppu.Memory[addr]
	}
}

func (ppu *Ppu) setPpuAddr(addr uint8) {
	if ppu.ppuAddrCounter == 0 {
		ppu.ppuAddrMSB = addr
		ppu.ppuAddrCounter = 1
	} else {
		ppu.ppuAddrLSB = addr
		ppu.ppuAddrCounter = 0
		ppu.ppuAddrOffset = 0
	}

	//log.Printf("setting ppu address 0x%x", addr)
}

func (ppu *Ppu) setPpuScrollAddr(addr uint8) {
	if ppu.ppuAddrCounter == 0 {
		ppu.ppuScrollMSB = addr
		ppu.ppuAddrCounter = 1
	} else {
		ppu.ppuScrollLSB = addr
		ppu.ppuAddrCounter = 0
		ppu.ppuAddrOffset = 0
	}
}

func (ppu *Ppu) SetVBlank() {
	ppu.nes.CPU.Memory[0x2002] |= 1 << 7
}

func (ppu *Ppu) ClearVBlank() {
	ppu.nes.CPU.Memory[0x2002] &= ^(uint8(1) << 7)
}

func (ppu *Ppu) get8x8Tile(base uint16, pos uint16) [8][8]uint8 {
	b0 := []byte{}
	for i := uint16(0); i < 8; i++ {
		b0 = append(b0, ppu.Read8(base+(pos*0x10)+i))
	}
	b1 := []byte{}
	for i := uint16(8); i < 16; i++ {
		b1 = append(b1, ppu.Read8(base+(pos*0x10)+i))
	}

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

func (ppu *Ppu) get8x16Tile(base uint16, pos uint16) [16][8]uint8 {
	tile8x8_1 := ppu.get8x8Tile(base, pos)
	tile8x8_2 := ppu.get8x8Tile(base, pos + 1)


	var result [16][8]uint8

	copy(result[0:8], tile8x8_1[:])
	copy(result[8:16], tile8x8_2[:])

	return result
}

func (ppu *Ppu) get2x2Attribute(base uint16, pos uint8) [2][2]uint8 {
	var result [2][2]uint8

	b := ppu.Read8(base+0x3C0+uint16(pos))

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
			palette[ppu.Read8(0x3F00) % 64],
			palette[ppu.Read8(0x3F01) % 64],
			palette[ppu.Read8(0x3F02) % 64],
			palette[ppu.Read8(0x3F03) % 64],
		}
	case 1:
		return [4]Color{
			palette[ppu.Read8(0x3F00) % 64],
			palette[ppu.Read8(0x3F05) % 64],
			palette[ppu.Read8(0x3F06) % 64],
			palette[ppu.Read8(0x3F07) % 64],
		}
	case 2:
		return [4]Color{
			palette[ppu.Read8(0x3F00) % 64],
			palette[ppu.Read8(0x3F09) % 64],
			palette[ppu.Read8(0x3F0A) % 64],
			palette[ppu.Read8(0x3F0B) % 64],
		}
	case 3:
		return [4]Color{
			palette[ppu.Read8(0x3F00) % 64],
			palette[ppu.Read8(0x3F0D) % 64],
			palette[ppu.Read8(0x3F0E) % 64],
			palette[ppu.Read8(0x3F0F) % 64],
		}
	}

	return [4]Color{}
}

func (ppu *Ppu) getSpriteColorPalette(pos uint8) [4]Color {
	switch pos {
	case 0:
		return [4]Color{
			palette[ppu.Read8(0x3F00) % 64],
			palette[ppu.Read8(0x3F11) % 64],
			palette[ppu.Read8(0x3F12) % 64],
			palette[ppu.Read8(0x3F13) % 64],
		}
	case 1:
		return [4]Color{
			palette[ppu.Read8(0x3F00) % 64],
			palette[ppu.Read8(0x3F15) % 64],
			palette[ppu.Read8(0x3F16) % 64],
			palette[ppu.Read8(0x3F17) % 64],
		}
	case 2:
		return [4]Color{
			palette[ppu.Read8(0x3F00) % 64],
			palette[ppu.Read8(0x3F19) % 64],
			palette[ppu.Read8(0x3F1A) % 64],
			palette[ppu.Read8(0x3F1B) % 64],
		}
	case 3:
		return [4]Color{
			palette[ppu.Read8(0x3F00) % 64],
			palette[ppu.Read8(0x3F1D) % 64],
			palette[ppu.Read8(0x3F1E) % 64],
			palette[ppu.Read8(0x3F1F) % 64],
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

func (ppu *Ppu) getColorFromSprite(x, y uint8, s Sprite) Color {
	flipHorizontal := (s.attributes >> 6) & 1 == 1
	flipVertical := (s.attributes >> 7) & 1 == 1

	var backgroundTileBase uint16
	var tile8x8 [8][8]uint8
	var tile8x16 [16][8]uint8
	var backgroundTilePos uint8

	// Check if we are in 8x16 sprite mode
	if ppu.ppuctrl.spriteSize == 0 {
		backgroundTileBase = uint16(ppu.ppuctrl.spritePatternTableAddr) * 0x1000
		backgroundTilePos = s.tileNum
		tile8x8 = ppu.get8x8Tile(backgroundTileBase, uint16(backgroundTilePos))

		xBG := x % 8
		if flipHorizontal {
			xBG = 7 - x % 8
		}

		yBG := y % 8
		if flipVertical {
			yBG = 7 - y % 8
		}

		spriteColorPalette := ppu.getSpriteColorPalette(s.attributes & 0x03)
		spriteColor := spriteColorPalette[tile8x8[yBG][xBG]]

		// Check if sprites hide background
		if tile8x8[yBG][xBG] == 0 {
			spriteColor.A = 0
		}

		return spriteColor
	} else {
		backgroundTileBase = uint16(s.tileNum & 1) * 0x1000
		backgroundTilePos = s.tileNum & ^uint8(0x01)
		tile8x16 = ppu.get8x16Tile(backgroundTileBase, uint16(backgroundTilePos))

		xBG := x % 8
		if flipHorizontal {
			xBG = 7 - x % 8
		}

		yBG := y % 16
		if flipVertical {
			yBG = 15 - y % 16
		}

		spriteColorPalette := ppu.getSpriteColorPalette(s.attributes & 0x03)
		spriteColor := spriteColorPalette[tile8x16[yBG][xBG]]

		// Check if sprites hide background
		if tile8x16[yBG][xBG] == 0 {
			spriteColor.A = 0
		}

		return spriteColor
	}
}

func (ppu *Ppu) getSpriteColorAtPixel(x, y uint8) Color {
	for id := 0; id < ppu.spriteCount; id++ {
		sprite := ppu.currentSprites[id]

		var ySpriteOffset uint8
		if ppu.ppuctrl.spriteSize == 1 {
			ySpriteOffset = 16
		} else {
			ySpriteOffset = 8
		}

		inRangeX := x >= sprite.xCoord && x < sprite.xCoord+8
		inRangeY := y >= sprite.yCoord && y < sprite.yCoord+ySpriteOffset

		if inRangeX && inRangeY {
			// trigger sprite 0 hit
			if id == 0 && ppu.ppumask.backgroundEnable && ppu.ppumask.spriteEnable {
				ppu.setSpriteHit()
			}

			spriteColor := ppu.getColorFromSprite(x-sprite.xCoord, y-sprite.yCoord, sprite)
			if spriteColor.A > 0 {
				return spriteColor
			}
		}
	}

	return Color{}
}

func (ppu *Ppu) fetchSprites() {
	ppu.spriteCount = 0

	var heightOffset uint16

	if ppu.ppuctrl.spriteSize == 0 {
		heightOffset = 8
	} else {
		heightOffset = 16
	}

	for _, sprite := range ppu.OAM {
		isSpriteVisible := sprite.yCoord > 0x00 && sprite.yCoord < 0xEF
		isSpriteOnCurrentScanline := ppu.Scanline +  1 - uint16(sprite.yCoord) < heightOffset
		if isSpriteVisible && isSpriteOnCurrentScanline {
			if ppu.spriteCount > 8 {
				return
			}

			ppu.currentSprites[ppu.spriteCount] = sprite
			ppu.spriteCount++
		}
	}
}

func (ppu *Ppu) is8x16Mode() bool {
	return (ppu.nes.CPU.Read8(0x2000) >> 5) & 1 == 1
}

func (ppu *Ppu) setSpriteHit() {
	ppu.nes.CPU.Memory[0x2002] |= 0x40
}

func (ppu *Ppu) clearSpriteHit() {
	ppu.nes.CPU.Memory[0x2002] &= 0xBF
}

func (ppu *Ppu) getAllTilesForScanline(nameTableY, nameTableX uint16) {
	ppu.currentTiles = [0x21][8][8]uint8{}
	ppu.currentAttributes = [0x21][2][2]uint8{}
	nameTableSelect := ppu.nes.CPU.Memory[0x2000] & 0x03

	for x := uint16(0); x < 0x21; x++ {
		backgroundTileOffset := ((nameTableY % 240/8) * 32) + (((nameTableX + 8 * x) % 256 / 8) % 32)
		nameTableBase := (0x2000 + (uint16(nameTableSelect)*0x400) + (nameTableY / 240) * 0x800 + ((nameTableX + 8 * x) / 256) * 0x400) & 0x2FFF
		backgroundTileBase := uint16((ppu.nes.CPU.Memory[0x2000]>>4)&1) * 0x1000
		backgroundTilePos := ppu.Memory[nameTableBase+backgroundTileOffset]
		backgroundTile := ppu.get8x8Tile(backgroundTileBase, uint16(backgroundTilePos))
		attributePalettePos := uint8((nameTableY % 240/32)*8) + ((uint8((nameTableX + 8 * x) % 256) / 32) % 32)
		attributeTile := ppu.get2x2Attribute(nameTableBase, attributePalettePos)
		copy(ppu.currentTiles[x][0:8], backgroundTile[0:8])
		copy(ppu.currentAttributes[x][0:2], attributeTile[0:2])
	}
}

func (ppu *Ppu) PPURun() {

	if ppu.Scanline == 0 {
		ppu.NmiOccurred = false
		ppu.clearSpriteHit()
	}

	if ppu.Cycle == 340 && ppu.Scanline < 240 {
		sl := ppu.Scanline
		nameTableY := sl + uint16(ppu.ppuScrollLSB)

		if nameTableY % 8 == 0 {
			ppu.getAllTilesForScanline(nameTableY, uint16(ppu.ppuScrollMSB))
		}

		for x := 0; x < 256; x++ {
			nameTableX := uint16(x) + uint16(ppu.ppuScrollMSB)

			currentTileIdx := (int(ppu.ppuScrollMSB % 8) + x) / 8 % 0x21


			c := ppu.getSpriteColorAtPixel(uint8(x), uint8(sl))
			if c.A == 0 {
				c = ppu.getBackgroundColorAtPixelOptimized(uint8(nameTableX % 256), uint8(nameTableY % 240), ppu.currentTiles[currentTileIdx], ppu.currentAttributes[currentTileIdx])
			}

			for i := 0; i < ppu.scalingFactor; i++ {
				for j := 0; j < ppu.scalingFactor; j++ {
					ppu.Frame.SetRGBA(x * ppu.scalingFactor + i, int(sl) * ppu.scalingFactor + j, color.RGBA{c.R, c.G, c.B, uint8(c.A)})
				}
			}
		}
	}

	if ppu.Scanline >= 257 && ppu.Scanline <= 320 {
		ppu.SetOamAddr(0)
	}

	if ppu.Scanline == 259 && ppu.Cycle == 340 {
		ppu.ClearVBlank()
	}
	if ppu.Scanline == 241 && ppu.Cycle == 0 {
		ppu.SetVBlank()
	}

	if ppu.Scanline <= 240 && ppu.Cycle == 257 {
		ppu.fetchSprites()
	}

	ppu.Cycle = (ppu.Cycle + 1) % 341
	if ppu.Cycle == 340 {
		ppu.Scanline = (ppu.Scanline + 1) % 260
	}

	// if a frame is ready, set bool
	if ppu.Cycle == 0 && ppu.Scanline == 240 {
		ppu.FrameReady = true
	}
}

func (ppu *Ppu) RunPPUCycles(numOfCycles uint16) {
	for i := uint16(0); i < numOfCycles; i++ {
		ppu.PPURun()
	}
}

func (ppu *Ppu) InitFrame(scalingFactor int) {
	ppu.scalingFactor = scalingFactor
	ppu.Frame = image.NewRGBA(image.Rect(0, 0, 256 * scalingFactor, 240 * scalingFactor))
}