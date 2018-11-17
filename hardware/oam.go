package hardware

type Sprite struct {
	yCoord     uint8
	tileNum    uint8
	attributes uint8
	xCoord     uint8
}

func (ppu *Ppu) WriteOAM8(value uint8) {
	switch ppu.oamAddr % 4 {
	case 0:
		ppu.OAM[ppu.oamSpriteAddr%64].yCoord = value
	case 1:
		ppu.OAM[ppu.oamSpriteAddr%64].tileNum = value
	case 2:
		ppu.OAM[ppu.oamSpriteAddr%64].attributes = value
	case 3:
		ppu.OAM[ppu.oamSpriteAddr%64].xCoord = value
		ppu.oamSpriteAddr++
	}

	ppu.oamAddr++
}

func (ppu *Ppu) SetOamAddr(addr uint8) {
	ppu.oamAddr = addr
}
