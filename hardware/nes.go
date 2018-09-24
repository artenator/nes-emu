package hardware

type NES struct {
	CPU *Cpu
	PPU *Ppu
}

func NewNES() NES {
	newNes := NES{}
	newNes.CPU = &Cpu{}
	newNes.PPU = &Ppu{}
	newNes.CPU.nes = &newNes
	newNes.PPU.nes = &newNes
	newNes.PPU.ppuAddrCounter = 0

	return newNes
}