package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/shirou/gopsutil/cpu"
	"golang.org/x/image/colornames"
	"log"
	"nes-emu/hardware"
	"time"
)

var imd = imdraw.New(nil)
var cpuInfo, _ = cpu.Info()

type ColorAtPixel struct {
	color hardware.Color
	x uint8
	y uint8
}

func drawPixel(c hardware.Color, x, y float64) {
	colorRGBA := pixel.RGB(float64(c.R)/255, float64(c.G)/255, float64(c.B)/255)
	imd.Color = colorRGBA
	imd.Push(pixel.V(x+1, y+1))
	imd.Push(pixel.V(x, y))
	imd.Rectangle(0)
}

func runNEStoFrame(nes hardware.NES, numOfInstructions *uint) {
	for !nes.PPU.FrameReady {
		//wait := 0
		opcode := nes.CPU.Read8(nes.CPU.PC)
		instr := hardware.Instructions[opcode]
		nes.CPU.RunInstruction(instr, false)
		nes.PPU.RunPPUCycles(uint16(3 * instr.Cycles))

		inVBlank := (nes.CPU.Memory[0x2002]>>7)&1 == 1
		NMIEnabled := (nes.CPU.Memory[0x2000]>>7)&1 == 1

		*numOfInstructions++

		if inVBlank && NMIEnabled {
			nes.PPU.ClearVBlank()
			nes.CPU.HandleNMI()
		}
	}

	nes.PPU.FrameReady = false
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

	nes := hardware.NewNES()

	cart, err := hardware.CreateCartridge("donkey-kong.nes")

	if err != nil {
		log.Println(err)
	} else {
		nes.LoadCartridge(cart)
		nes.CPU.Reset()
	}

	var numOfInstructions uint = 0
	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	// main drawing loop
	for !win.Closed() {
		runNEStoFrame(nes, &numOfInstructions)
		imd.Clear()

		for y := 0; y < 240; y++ {
			for x := 0; x < 256; x++ {
				c := nes.PPU.GetColorAtPixel(uint8(x), uint8(y))
				drawPixel(c, float64(x), float64(239 - y))
			}
		}

		win.Clear(colornames.Black)
		imd.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	log.Println("Arte's NES Emu")
	pixelgl.Run(run)
}
