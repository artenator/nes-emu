package hardware

import (
	"image"
)

type NES struct {
	CPU  *Cpu
	PPU  *Ppu
	APU  *Apu
	CART *Cartridge
	CARTIO CartridgeIO
	Xwin, Ywin int
}

func NewNES(xWin, yWin int) *NES {
	newNes := NES{}
	newNes.CPU = &Cpu{}
	newNes.PPU = &Ppu{}
	newNes.APU = &Apu{}
	newNes.PPU.Frame = image.NewRGBA(image.Rect(0, 0, xWin, yWin))
	newNes.CPU.nes = &newNes
	newNes.PPU.nes = &newNes
	newNes.APU.nes = &newNes
	newNes.PPU.ppuAddrCounter = 0
	newNes.Xwin = xWin
	newNes.Ywin = yWin

	return &newNes
}
