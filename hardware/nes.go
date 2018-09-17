package hardware

type NES struct {
	CPU Cpu
	PPU Ppu
}

func NewNES() NES {
	newNes := NES{}
	newNes.CPU = Cpu{}
	newNes.PPU = Ppu{}
	newNes.PPU.CPU = &newNes.CPU

	return newNes
}