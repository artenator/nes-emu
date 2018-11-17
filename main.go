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
	imd.Color = pixel.RGB(float64(float64(c.R)/255), float64(float64(c.G)/255), float64(float64(c.B)/255))
	imd.Push(pixel.V(x+1, y+1))
	imd.Push(pixel.V(x, y))
	imd.Rectangle(0)
}

func runNES(nes hardware.NES, numOfInstructions *uint) {
	for true {
		wait := 0
		opcode := nes.CPU.Read8(nes.CPU.PC)
		nes.CPU.RunInstruction(hardware.Instructions[opcode], false)

		//time.Sleep(1 * time.Nanosecond)
		for wait < 100000 {
			wait++
		}

		*numOfInstructions++

		if (nes.CPU.Memory[0x2002]>>7)&1 == 1 && (nes.CPU.Memory[0x2000]>>7)&1 == 1 && *numOfInstructions%1000 == 0 {
			nes.CPU.HandleNMI()
		}
	}
}

func run() {
	// all of our code will be fired up from here
	cfg := pixelgl.WindowConfig{
		Title:  "Arte's NES Emulator",
		Bounds: pixel.R(0, 0, 256, 240),
		VSync:  false,
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

	// Run
	go runNES(nes, &numOfInstructions)

	// main drawing loop
	for !win.Closed() {
		imd.Clear()

		if numOfInstructions%500 == 0 {
			for y := 0; y < 240; y++ {
				if y == 1 {
					nes.PPU.ClearVBlank()
				}
				for x := 0; x < 256; x++ {
					drawPixel(nes.PPU.GetColorAtPixel(uint8(x), uint8(y)), float64(x), float64(239-y))
				}
				if y == 200 {
					nes.PPU.SetVBlank()

				}
			}

			win.Clear(colornames.Black)
			imd.Draw(win)
			win.Update()
		} else {
			win.Update()
		}

	}
}

func main() {
	log.Println("Arte's NES Emu")
	pixelgl.Run(run)
}
