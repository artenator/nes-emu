package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"log"
	"nes-emu/hardware"
)

var imd = imdraw.New(nil)

func drawPixel(x, y float64) {
	imd.Color = pixel.RGB(1, 0, 0)
	imd.Push(pixel.V(x + 1, y + 1))
	imd.Push(pixel.V(x, y))
	imd.Rectangle(0)
}

func run() {
	// all of our code will be fired up from here
	cfg := pixelgl.WindowConfig{
		Title:  "Arte's NES Emulator",
		Bounds: pixel.R(0, 0, 256, 240),
		VSync: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	drawPixel(0, 0)
	drawPixel(100, 100)

	for !win.Closed() {
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
	log.Println("Arte's NES Emu")

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
	// log.Printf("%x", rom)
	// log.Println(err)
}

