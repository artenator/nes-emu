package hardware

import (
	"errors"
	"io/ioutil"
	"log"
)

type Cartridge struct {
	nes *NES

	// Constant for ines headers
	nesLabel [4]byte

	// Size of the PRG rom in 16KB chunks
	prgRomBlocks byte

	// Size of the CHR rom in 8KB chunks
	chrRomBlocks byte

	// Flags 6
	flags6 byte

	// Flags 7
	flags7 byte

	// Size of PRG ram in 8KB chunks
	prgRamBlocks byte

	// Flags 9
	flags9 byte

	// Flags 10
	flags10 byte

	// Zero filled rest of header
	zeroBuffer [5]byte

	// actual prgRom data
	prgRom []byte

	// actual chrRom data
	chrRom []byte

	//Mapper type
	mapperType byte

	//Mirroring style
	mirrorStyle byte
}

const (
	horizontal = iota
	vertical = iota
	singleScreen = iota
	fourScreen = iota
)

const (
	mapper0 = iota
	mmc1 = iota
)

func (c *Cartridge) setMapperType() {
	mapperLow := (c.flags6 & 0xF0) >> 4
	mapperHigh := c.flags7 & 0xF0

	c.mapperType = mapperHigh | mapperLow
}

func CreateCartridge(filename string) (Cartridge, error) {
	// Read nes rom into memory
	rom, err := ioutil.ReadFile(filename)

	romNoHeader := rom[16:len(rom)]

	var c Cartridge

	// If file read was unsuccessful, log it
	if err != nil {
		log.Fatal("Something went wrong during file read. Error: " + err.Error())
	} else {
		copy(c.nesLabel[:], rom[0:4])

		// Make sure this is an NES rom
		if string(c.nesLabel[:]) == "NES" + string(0x1a) {
			c.prgRomBlocks = rom[4]
			c.chrRomBlocks = rom[5]
			c.flags6 = rom[6]
			c.flags7 = rom[7]
			c.prgRamBlocks = rom[8]
			c.flags9 = rom[9]
			c.flags10 = rom[10]
			copy(c.zeroBuffer[:], rom[11:15])

			c.setMapperType()
			log.Printf("Mapper type %d", c.mapperType)

			// load chr and prg rom data
			log.Println(c.prgRomBlocks, c.chrRomBlocks)
			if c.prgRomBlocks > 0 {
				c.prgRom = romNoHeader[0:0x4000 * uint(c.prgRomBlocks)]
			}

			if c.chrRomBlocks > 0 {
				c.chrRom = romNoHeader[0x4000:0x4000 + (0x2000 * uint(c.chrRomBlocks))]
			}
                        
		} else {
			log.Println("This is not a valid NES rom.")
			return c, errors.New("This is not a valid NES rom.")
		}
	}

	return c, nil
}

func (nes *NES) LoadCartridge(cartridge Cartridge) {
	nes.CART = &cartridge
	nes.CART.nes = nes

	switch cartridge.mapperType {
	case mapper0:
		mapper := &Mapper0CIO{}
		mapper.initCartIO(&cartridge)
		nes.CARTIO = mapper
	case mmc1:
		mapper := &Mapper1CIO{}
		mapper.initCartIO(&cartridge)
		nes.CARTIO = mapper
		//copy(nes.PPU.Memory[0:0x2000], cartridge.chrRom[0:0x2000])

	default:
		log.Fatalf("Unsupported mapper %d", cartridge.mapperType)
	}
}
