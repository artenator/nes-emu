package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"log"
	"nes-emu/hardware"
)

var imd = imdraw.New(nil)

func drawPixel(c hardware.Color, x, y float64) {
	imd.Color = pixel.RGB(float64(c.R), float64(c.G), float64(c.B))
	imd.Push(pixel.V(x + 1, y + 1))
	imd.Push(pixel.V(x, y))
	imd.Rectangle(0)
}

func runNES(nes hardware.NES, numOfInstructions * uint) {
	for true {
		opcode := nes.CPU.Read8(nes.CPU.PC)
		nes.CPU.RunInstruction(hardware.Instructions[opcode], false)

		//time.Sleep(500 * time.Nanosecond)

		*numOfInstructions++

		if *numOfInstructions % 1000 == 0 {
			if (nes.CPU.Memory[0x2002] >> 7) & 1 == 0 {
				nes.PPU.SetVBlank()
			} else {
				nes.PPU.ClearVBlank()
			}
		}
	}
}

func run() {
	// all of our code will be fired up from here
	cfg := pixelgl.WindowConfig{
		Title:  "Arte's NES Emulator",
		Bounds: pixel.R(0, 0, 256, 240),
		VSync: false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// create new nes
	//cpu := hardware.Cpu{}
	nes := hardware.NewNES()

	cart, err := hardware.CreateCartridge("donkey-kong.nes")

	if err != nil {
		log.Println(err)
	} else {
		nes.LoadCartridge(cart)
		//log.Println(cart)
		//log.Println(cpu.Memory)

		nes.CPU.Reset()
	}

	var numOfInstructions uint = 0

	go runNES(nes, &numOfInstructions)

	// main drawing loop
	for !win.Closed() {
		imd.Clear()



		if numOfInstructions % 1000 == 0 {

			for y := 0; y < 240; y++ {
				for x := 0; x < 256; x++ {
					drawPixel(nes.PPU.GetColorAtPixel(uint8(x), uint8(y)), float64(x), float64(239 - y))
				}
			}

			win.Clear(colornames.Black)
			imd.Draw(win)
			win.Update()
			//log.Println("drawing to screen...")
			//log.Printf("%+v", nes.PPU.Memory[0x2000:0x2050])
		} else {
			win.Update()
		}


	}
}

func main() {
	pixelgl.Run(run)
	log.Println("Arte's NES Emu")


	// log.Printf("%x", rom)
	// log.Println(err)
}

