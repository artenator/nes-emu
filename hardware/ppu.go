package hardware

import (
	"encoding/binary"
	"log"
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

	if value == 0x24 {
		log.Println("Hello")
	}

	if ppu.nes.CPU.Memory[0x2002] != 0xA0 {
		log.Println("Hello")
	}

	ppu.Memory[ppuWriteAddress + ppu.ppuAddrOffset] = value
	ppu.ppuAddrOffset++
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