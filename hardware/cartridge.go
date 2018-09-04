package hardware

import ("io/ioutil"
	"log"
	"errors")

type Cartridge struct {
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

			// load chr and prg rom data
			c.prgRom = romNoHeader[0:0x4000 * uint(c.prgRomBlocks)]
			c.chrRom = romNoHeader[0:0x2000 * uint(c.chrRomBlocks)]
                        
		} else {
			log.Println("This is not a valid NES rom.")
			return c, errors.New("This is not a valid NES rom.")
		}
	}

	return c, nil
}

func (cpu *Cpu) LoadCartridge(cartridge Cartridge) {
        // if there's only one 16KB prg slot, mirror it in the cpu memory
	if uint(cartridge.prgRomBlocks) == 1 {
		copy(cpu.Memory[0x8000:0xC000], cartridge.prgRom[0:0x4000])
		copy(cpu.Memory[0xC000:0xFFFF], cartridge.prgRom[0:0x4000])
	}
}
