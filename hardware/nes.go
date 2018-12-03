package hardware

import (
	"github.com/hajimehoshi/oto"
	"image"
	"log"
)

type NES struct {
	CPU *Cpu
	PPU *Ppu
	APU *Apu
}

func NewNES() NES {
	newNes := NES{}
	newNes.CPU = &Cpu{}
	newNes.PPU = &Ppu{}
	newNes.APU = &Apu{}
	newNes.PPU.Frame = image.NewRGBA(image.Rect(0, 0, 256, 240))
	newNes.CPU.nes = &newNes
	newNes.PPU.nes = &newNes
	newNes.APU.nes = &newNes
	newNes.PPU.ppuAddrCounter = 0

	// init audio player
	newNes.APU.cycleLimit = 40
	var err error
	if newNes.APU.audioDevice, err = oto.NewPlayer(44100, 1, 1, 4096); err != nil {
		log.Fatal("Audio could not be initialized")
	}

	return newNes
}
