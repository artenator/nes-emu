package hardware

type NES struct {
	CPU  *Cpu
	PPU  *Ppu
	APU  *Apu
	CART *Cartridge
	CARTIO CartridgeIO
}

func NewNES() *NES {
	newNes := NES{}
	newNes.CPU = &Cpu{}
	newNes.PPU = &Ppu{}
	newNes.APU = &Apu{}
	newNes.CPU.nes = &newNes
	newNes.PPU.nes = &newNes
	newNes.APU.nes = &newNes
	newNes.PPU.ppuAddrCounter = 0

	return &newNes
}
