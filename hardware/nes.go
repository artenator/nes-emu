package hardware

import "image"

type NES struct {
	CPU *Cpu
	PPU *Ppu
}

func NewNES() NES {
	newNes := NES{}
	newNes.CPU = &Cpu{}
	newNes.PPU = &Ppu{}
	newNes.PPU.Frame = image.NewRGBA(image.Rect(0, 0, 256, 240))
	newNes.CPU.nes = &newNes
	newNes.PPU.nes = &newNes
	newNes.PPU.ppuAddrCounter = 0

	return newNes
}