package main

import "C"
import (
	"errors"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/shirou/gopsutil/cpu"
	"golang.org/x/image/colornames"
	"log"
	"nes-emu/hardware"
	"os"
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

func runNESInstruction(nes hardware.NES, numOfInstructions *uint) {
	opcode := nes.CPU.Read8(nes.CPU.PC)
	instr := hardware.Instructions[opcode]
	nes.CPU.RunInstruction(instr, false)
	nes.PPU.RunPPUCycles(uint16(3 * instr.Cycles))
	nes.APU.RunAPUCycles(uint16(instr.Cycles))

	inVBlank := (nes.CPU.Memory[0x2002]>>7)&1 == 1
	NMIEnabled := (nes.CPU.Memory[0x2000]>>7)&1 == 1

	*numOfInstructions++

	if inVBlank && NMIEnabled {
		nes.PPU.ClearVBlank()
		nes.CPU.HandleNMI()
	}
}

func runNEStoFrame(nes hardware.NES, numOfInstructions *uint) {
	for !nes.PPU.FrameReady {
		runNESInstruction(nes, numOfInstructions)
	}

	//log.Printf("%+v", nes.APU.GetPulseFrequency(0x4000))
	//log.Printf("%+v", nes.APU.GetPulseFrequency(0x4004))

	nes.PPU.FrameReady = false
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Arte's NES Emulator",
		Bounds: pixel.R(0, 0, 256, 240),
		VSync:  false,
	}

	var gameName string
	if len(os.Args) > 1 {
		gameName = os.Args[1]
	} else {
		panic(errors.New("Please pass in a game name."))
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	nes := hardware.NewNES()

	cart, err := hardware.CreateCartridge(gameName)

	if err != nil {
		log.Println(err)
	} else {
		nes.LoadCartridge(cart)
		nes.CPU.Reset()
	}

	var (
		numOfInstructions uint = 0
		frames = 0
		us = time.Tick(16666 * time.Microsecond)
		second = time.Tick(time.Second)
	)

	// main drawing loop
	for !win.Closed() {
		imd.Clear()

		nes.CPU.CheckControllerPresses(win)

		runNEStoFrame(nes, &numOfInstructions)

		pic := pixel.PictureDataFromImage(nes.PPU.Frame)

		sprite := pixel.NewSprite(pic, pic.Bounds())

		win.Clear(colornames.Black)

		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		win.Update()

		<-us
		frames++

		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("FPS: %d %s", frames, cfg.Title))
			nes.APU.Cyclelimit = uint8(29829 * frames / 44100)
			frames = 0
		default:
		}
	}
}

func main() {
	log.Println("Arte's NES Emu")
	pixelgl.Run(run)
}
